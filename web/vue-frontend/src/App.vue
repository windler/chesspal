<template>
  <v-app>
    <v-main>
      <v-app-bar dark>
        <v-img class="logo" src="/chesspal.svg"></v-img>

        <v-tabs align-with-title v-model="tab">
          <v-tab href="#tab-1">
            <v-icon>fas fa-chess-pawn</v-icon>&nbsp;Game
          </v-tab>
          <v-tab href="#tab-2">
            <v-icon>fas fa-clock-rotate-left</v-icon>&nbsp;History
          </v-tab>
          <v-tab href="#tab-3">
            <v-icon>fas fa-robot</v-icon>&nbsp;Bot games
          </v-tab>
        </v-tabs>
      </v-app-bar>

      <v-footer fixed dark padless>
        <v-row justify="center" no-gutters>
          <v-icon class="mx-2" :color="connected ? 'green' : 'red'"
            >fa fa-signal</v-icon
          >
          <v-btn
            class="my-auto"
            icon
            href="https://github.com/windler/chesspal"
            target="_blank"
          >
            <v-icon>fab fa-github</v-icon>
          </v-btn>

          <v-btn class="my-auto" icon @click.stop="toggleDarkTheme()">
            <v-icon>fa fa-moon</v-icon>
          </v-btn>
        </v-row>
      </v-footer>

      <v-tabs-items v-model="tab">
        <v-tab-item key="1" value="tab-1">
          <v-container>
            <v-row class="justify-center">
              <v-col cols="12" lg="8">
                <ChessBoard
                  :svg="
                    nextBestPosition != '' && showHint
                      ? nextBestPosition
                      : currentPosition
                  "
                  :fen="fen"
                  :outcome="outcome"
                  :pgn="pgn"
                  class="my-4"
                />
              </v-col>

              <v-col cols="12" lg="3">
                <EvalInfo
                  :pawn="pawn"
                  :class="evalMode == 1 ? 'my-4' : 'd-none'"
                />
                <MoveList
                  :movesBlack="movesBlack"
                  :movesWhite="movesWhite"
                  :showEvaluation="evalMode == 1"
                  class="my-4"
                  height="350px"
                />
                <GameActions
                  v-on:undoMoves="undoMoves($event)"
                  v-on:draw="draw()"
                  v-on:resign="resign()"
                  class="my-4"
                  v-on:showHint="showHint = true"
                  v-on:changeMode="evalMode = $event"
                />
                <v-dialog
                  overlay-opacity="0.95"
                  max-width="350px"
                  v-model="startdialog"
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      :disabled="!connected"
                      color="primary"
                      width="100%"
                      v-bind="attrs"
                      v-on="on"
                      v-show="!started"
                    >
                      New game
                    </v-btn>
                    
                  </template>
                  <v-container>
                    <v-row class="justify-center">
                      <v-col cols="12">
                        <ChessPlayer
                          v-on:nameChange="white.name = $event"
                          v-on:modeChange="white.mode = $event"
                          :locked="started"
                          color="white"
                          :name="white.name"
                          :bots="bots"
                          class="my-4"
                        />
                        <ChessPlayer
                          v-on:nameChange="black.name = $event"
                          v-on:modeChange="black.mode = $event"
                          :locked="started"
                          color="black"
                          :name="black.name"
                          class="my-4"
                          :bots="bots"
                        />
                        <SettingsCard
                          :locked="started"
                          v-on:upsideDownChange="upsideDown = $event"
                          v-on:speakChange="
                            white.speak = Boolean($event);
                            black.speak = Boolean($event);
                          "
                          :speak="white.speak || black.speak ? 'true' : 'false'"
                          class="my-4"
                        />
                        <v-btn
                          color="primary"
                          width="100%"
                          @click.stop="startGame()"
                        >
                          Start game
                        </v-btn>
                        <v-progress-linear
                          indeterminate
                          width="100%"
                          color="primary"
                          v-if="startSend"
                          height="20"
                        ></v-progress-linear>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-dialog>
                <div class="my-12"></div>
              </v-col>
            </v-row>
          </v-container>
        </v-tab-item>

        <v-tab-item key="2" value="tab-2">
          <GameHistory
            showArchived="false"
            showBotGames="false"
            showHumanGames="true"
          />
        </v-tab-item>
        <v-tab-item key="3" value="tab-3">
          <GameHistory
            showArchived="false"
            showBotGames="true"
            showHumanGames="false"
          />
        </v-tab-item>
      </v-tabs-items>
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
import GameHistory from "./components/GameHistory.vue";

export default {
  name: "App",

  components: {
    ChessPlayer,
    EvalInfo,
    MoveList,
    ChessBoard,
    SettingsCard,
    GameActions,
    GameHistory,
  },

  data: () => ({
    startdialog: false,
    startSend: false,
    tab: null,
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
      name: "black",
      mode: 0,
      speak: false,
    },
    white: {
      name: "white",
      mode: 0,
      speak: false,
    },
    evalMode: 0,
    pgn: "",
    fen: "r5nr/ppk2pp1/7p/2Bp1b2/8/7P/PPP1PPP1/RN2KB1R",
    outcome: "*",
    bots: [],
  }),

  methods: {
    toggleDarkTheme() {
      this.$vuetify.theme.dark = !this.$vuetify.theme.dark;
    },
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

        this.speak(text);
        this.lastSpoken = move + player.name;
      }
    },
    speak: function (text) {
      this.speech.text = text;

      this.speech.rate = 0.4;
      window.speechSynthesis.speak(this.speech);
    },
    startGame: function () {
      if (!this.started) {
        this.startSend = true;
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

  mounted() {
    if (localStorage.white) {
      this.white.name = localStorage.white;
    }
    if (localStorage.black) {
      this.black.name = localStorage.black;
    }
  },

  watch: {
    "white.name": function (val) {
      localStorage.white = val;
    },
    "black.name": function (val) {
      localStorage.black = val;
    },
  },

  created: async function () {
    const connect = () => {
      this.speech = new SpeechSynthesisUtterance();
      this.voices = window.speechSynthesis.getVoices();
      this.speech.lang = "en";

      console.log("Starting connection to WebSocket Server");
      var host = location.host;
      if (process.env.VUE_APP_CHESSPAL_HOST !== undefined) {
        host = process.env.VUE_APP_CHESSPAL_HOST;
      }
      this.connection = new WebSocket("ws://" + host + "/ws");
      var that = this;

      this.connection.onmessage = function (event) {
        var data = JSON.parse(event.data);

        if (data.bots != null) {
          that.bots = data.bots;
          return;
        }

        if (data.started) {
          that.speak("Game started!");
          that.started = true;
          that.startdialog = false;
          that.startSend = false;
          return;
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

        if (data.moves != null) {
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
        }

        that.movesWhite = movesWhite;
        that.movesBlack = movesBlack;

        if (data.turn != that.turn) {
          that.showHint = false;
        }

        that.turn = data.turn;

        if (movesBlack.length > 0) {
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
        }

        that.pgn = data.pgn;
        that.fen = data.fen;
        that.outcome = data.outcome;

        if (that.outcome != "*") {
          that.started = false;
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
    };
    const sleep = (milliseconds) => {
      return new Promise((resolve) => setTimeout(resolve, milliseconds));
    };

    while (!this.connected) {
      connect();
      await sleep(3000);
    }
  },
};
</script>
<style >
body {
  overflow: hidden;
  width: 100vw;
}
.logo {
  filter: invert(100%);
  width: 80px;
}
</style>