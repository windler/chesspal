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
              <EvaluationMode
                v-on:changeMode="evalMode = $event"
                :locked="started"
                class="my-6"
              />
            </v-sheet>
          </v-col>

          <v-col cols="12" sm="6">
            <v-sheet min-height="70vh" rounded="lg">
              <v-row justify="center">
                <ChessBoard :svg="currentPosition" class="my-6" />
              </v-row>
            </v-sheet>
          </v-col>

          <v-col cols="12" sm="3">
            <v-sheet rounded="lg" min-height="268">
              <EvalInfo :pawn="pawn" :show="evalMode == 1" class="my-6" />
              <MoveList
                :movesBlack="movesBlack"
                :movesWhite="movesWhite"
                class="my-6"
              />
              <PGNCard :pgn="pgn" class="my-6" />
            </v-sheet>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import ChessPlayer from "./components/ChessPlayer.vue";
import EvaluationMode from "./components/EvaluationMode.vue";
import EvalInfo from "./components/EvalInfo.vue";
import MoveList from "./components/MoveList.vue";
import ChessBoard from "./components/ChessBoard.vue";
import PGNCard from "./components/PGNCard.vue";

export default {
  name: "App",

  components: {
    ChessPlayer,
    EvaluationMode,
    EvalInfo,
    MoveList,
    ChessBoard,
    PGNCard,
  },

  data: () => ({
    connection: null,
    speech: null,
    lastMove: "",
    movesBlack: [],
    movesWhite: [],
    connected: false,
    started: false,
    currentPosition: "",
    pawn: 50.0,
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
  }),
  methods: {
    speakMove: function (player, move) {
      if (player.speak && !window.speechSynthesis.pending) {
        var text = move;
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
        window.speechSynthesis.speak(this.speech);
      }
    },
    startGame: function () {
      if (!this.started) {
        var msg = JSON.stringify({
          action: "start",
          options: {
            white: {
              name: this.white.name,
              type: Number(this.white.mode),
            },
            black: {
              name: this.black.name,
              type: Number(this.black.mode),
            },
            evalMode: Number(this.evalMode),
          },
        });

        this.connection.send(msg);
        console.log(msg);
        this.started = true;
      }
    },
  },

  created: function () {
    this.speech = new SpeechSynthesisUtterance();
    this.speech.lang = "en";

    console.log("Starting connection to WebSocket Server");
    this.connection = new WebSocket("ws://localhost:8080/ws");
    var that = this;

    this.connection.onmessage = function (event) {
      var data = JSON.parse(event.data);
      if (data.svgPosition != "") {
        that.currentPosition = data.svgPosition;
      }

      if (data.pawn !== 0.0) {
        that.pawn = data.pawn;
      }

      if (data.accuracy != "") {
        if (data.turn == "b") {
          that.movesWhite = that.movesWhite.map(function (move) {
            if (move.notation == data.lastMove) {
              move.accuracy = data.accuracy;
            }
            return move;
          });
        } else {
          that.movesBlack = that.movesBlack.map(function (move) {
            if (move.notation == data.lastMove) {
              move.accuracy = data.accuracy;
            }
            return move;
          });
        }
      }

      if (data.lastMove != that.lastMove) {
        that.lastMove = data.lastMove;
        if (data.turn == "b") {
          that.movesWhite.push({ notation: that.lastMove });
          that.speakMove(that.white, that.lastMove);
        } else {
          that.movesBlack.push({ notation: that.lastMove });
          that.speakMove(that.black, that.lastMove);
        }

        that.pgn = data.pgn;
      }
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
