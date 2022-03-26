<template>
  <v-app>
    <v-main>
      <v-app-bar color="deep-purple accent-4" dense dark>
        <v-toolbar-title>Chesspal</v-toolbar-title>
        <v-spacer></v-spacer>

        <v-btn v-if="!started" icon @click.stop="startGame()">
          <v-icon>fas fa-play</v-icon>
        </v-btn>
        <v-icon :color="connected ? 'green' : 'red'">fa fa-signal</v-icon>
      </v-app-bar>

      <v-container>
        <v-row class="justify-center">
          <v-col cols="12" sm="3">
            <v-sheet rounded="lg" min-height="268">
              <ChessPlayer
                v-on:nameChange="white.name = $event"
                v-on:modeChange="white.mode = $event"
                v-on:speakChange="white.speak = Boolean($event)"
                :locked="started"
                color="white"
                class="my-6"
              />
              <ChessPlayer
                v-on:nameChange="black.name = $event"
                v-on:modeChange="black.mode = $event"
                v-on:speakChange="black.speak = Boolean($event)"
                :locked="started"
                color="black"
                class="my-6"
              />
              <SettingsCard
                v-on:upsideDownChange="upsideDown = $event"
                :locked="started"
                class="my-6"
              />
            </v-sheet>
          </v-col>

          <v-col cols="12" sm="6">
            <v-sheet min-height="70vh" rounded="lg">
              <v-row justify="center">
                <ChessBoard
                  :svg="
                    nextBestPosition != '' && showHint
                      ? nextBestPosition
                      : currentPosition
                  "
                  :fen="fen"
                  :outcome="outcome"
                  :pgn="pgn" 
                  class="my-6"
                />
              </v-row>
            </v-sheet>
          </v-col>

          <v-col cols="12" sm="3">
            <v-sheet rounded="lg" min-height="268">
              <EvalInfo :pawn="pawn" :show="evalMode == 1" class="my-6" />
              <MoveList
                :movesBlack="movesBlack"
                :movesWhite="movesWhite"
                :showEvaluation="evalMode == 1"
                class="my-6"
              />
              <GameActions
                v-on:undoMoves="undoMoves($event)"
                v-on:draw="draw()"
                v-on:resign="resign()"
                v-on:showHint="showHint = true"
                v-on:changeMode="evalMode = $event"
              />
            </v-sheet>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import ChessPlayer from "./components/ChessPlayer.vue";
import EvalInfo from "./components/EvalInfo.vue";
import MoveList from "./components/MoveList.vue";
import ChessBoard from "./components/ChessBoard.vue";
import SettingsCard from "./components/SettingsCard.vue";
import GameActions from "./components/GameActions.vue";

export default {
  name: "App",

  components: {
    ChessPlayer,
    EvalInfo,
    MoveList,
    ChessBoard,
    SettingsCard,
    GameActions,
  },

  data: () => ({
    showHint: true,
    connection: null,
    upsideDown: false,
    speech: null,
    lastMove: "",
    movesBlack: [],
    movesWhite: [],
    connected: false,
    started: false,
    currentPosition: "",
    nextBestPosition: "",
    pawn: 50.0,
    turn: "w",
    black: {
      name: "Black",
      mode: 0,
      speak: false,
    },
    white: {
      name: "White",
      mode: 0,
      speak: false,
    },
    evalMode: 0,
    pgn: "",
    fen: "r5nr/ppk2pp1/7p/2Bp1b2/8/7P/PPP1PPP1/RN2KB1R",
    outcome: "*",
  }),
  methods: {
    speakMove: function (player, move) {
      if (
        player.speak &&
        !window.speechSynthesis.pending &&
        this.lastSpoken != move + player.name
      ) {
        var text = move;

        if (text.indexOf("-") == -1) {
          text = text.replace(/.{1}/g, "$&-");
        }
        text = text.replace(/K/g, "King ");
        text = text.replace(/N/g, "Knight ");
        text = text.replace(/B/g, "Bishop ");
        text = text.replace(/R/g, "Rook ");
        text = text.replace(/Q/g, "Queen ");
        text = text.replace(/x/g, " takes ");
        text = text.replace(/O-O-O/g, "Long Castles ");
        text = text.replace(/O-O/g, "Castles ");
        text = text.replace(/\+/g, " Check ");
        text = text.replace(/#/g, " Check mate ");

        this.speech.text = text;

        this.speech.rate = 0.4;
        window.speechSynthesis.speak(this.speech);
        this.lastSpoken = move + player.name;
      }
    },
    startGame: function () {
      if (!this.started) {
        var msg = JSON.stringify({
          action: "start",
          startOptions: {
            white: {
              name: this.white.name,
              type: Number(this.white.mode),
            },
            black: {
              name: this.black.name,
              type: Number(this.black.mode),
            },
            evalMode: 1, //always use eval but only show based on ui // Number(this.evalMode),
            upsideDown: Boolean(this.upsideDown),
          },
        });

        this.connection.send(msg);
        console.log(msg);
      }
    },
    undoMoves: function (n) {
      var msg = JSON.stringify({
        action: "undo",
        undoMoves: n,
      });

      this.connection.send(msg);
      console.log(msg);
    },
    draw: function () {
      var msg = JSON.stringify({
        action: "result",
        result: "draw",
      });

      this.connection.send(msg);
      console.log(msg);
    },
    resign: function () {
      var msg = JSON.stringify({
        action: "result",
        result: "resign",
      });

      this.connection.send(msg);
      console.log(msg);
    },
  },

  created: function () {
    this.speech = new SpeechSynthesisUtterance();
    this.voices = window.speechSynthesis.getVoices();
    this.speech.lang = "en";

    console.log("Starting connection to WebSocket Server");
    this.connection = new WebSocket("ws://localhost:8080/ws");
    var that = this;

    this.connection.onmessage = function (event) {
      var data = JSON.parse(event.data);

      if (data.started) {
        that.started = true;
        return
      }

      if (data.svgPosition != "") {
        that.currentPosition = data.svgPosition;
      }
      if (data.svgNextBestMove != "") {
        that.nextBestPosition = data.svgNextBestMove;
      }

      if (data.pawn !== 0.0) {
        that.pawn = data.pawn;
      }

      var movesWhite = [];
      var movesBlack = [];

      data.moves.forEach((m) => {
        var data = {
          notation: m.move,
          accuracy: m.accuracy,
        };
        if (m.color == "b") {
          movesBlack.push(data);
        } else {
          movesWhite.push(data);
        }
      });

      that.movesWhite = movesWhite;
      that.movesBlack = movesBlack;

      if (data.turn != that.turn) {
        that.showHint = false;
      }

      that.turn = data.turn;

      if (data.turn == "b") {
        that.speakMove(
          that.white,
          that.movesWhite[movesWhite.length - 1].notation
        );
      } else {
        that.speakMove(
          that.black,
          that.movesBlack[movesBlack.length - 1].notation
        );
      }

      that.pgn = data.pgn;
      that.fen = data.fen;
      that.outcome = data.outcome;
    };

    this.connection.onopen = function () {
      console.log("Successfully connected to the echo websocket server...");

      that.connected = true;
    };
    this.connection.onclose = function () {
      console.log("WS connection closed");
      that.connected = false;
    };
  },
};
</script>
