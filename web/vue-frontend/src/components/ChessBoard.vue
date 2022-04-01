<template>
  <v-card variant="outlined" min-width="720px" min-height="700px">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fas fa-chess-board</v-icon>
    </v-card-title>

    <div>
      <div :class="outcome != '*' ? 'dimmed white--text' : ''">
        {{ outcome != "*" ? outcome : "" }}
      </div>
      <div v-html="svg" class="board">{{ svg }}</div>
    </div>

    <v-divider class="ma-4"></v-divider>

    <v-card-actions>
      <v-btn color="primary" text> SHOW PGN </v-btn>

      <v-spacer></v-spacer>

      <v-btn icon @click="showPGN = !showPGN">
        <v-icon>{{ showPGN ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
      </v-btn>
    </v-card-actions>
    <v-expand-transition>
      <div v-show="showPGN" style="z-index: 100">
        <v-divider></v-divider>
        <v-card-text ref="pgn">
          <pre>{{ pgn }}</pre>
        </v-card-text>
        <v-btn class="ma-2" icon @click="importLichess()"
          ><v-icon>fas fa-magnifying-glass-chart</v-icon>
        </v-btn>
        <v-btn class="my-2" icon @click="copy()"
          ><v-icon>fas fa-copy</v-icon>
        </v-btn>
        <v-snackbar
          v-model="copied"
          timeout="2000"
          color="deep-purple accent-4"
        >
          <span class="text-center">Copied to clipboard!</span>
        </v-snackbar>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script>
export default {
  name: "ChessBoard",

  props: ["svg", "fen", "outcome", "pgn"],
  methods: {
    importLichess: async function () {
      const requestOptions = {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: "pgn=" + this.pgn,
      };
      const response = await fetch(
        "https://lichess.org/api/import",
        requestOptions
      );
      const data = await response.json();
      window.open(data.url, "_blank");
    },
    copy: function () {
      navigator.clipboard.writeText(this.pgn).then(
        () => {
          this.copied = true;
        },
        (err) => {
          console.error("Could not copy text: ", err);
        }
      );
    },
  },
  data() {
    return {
      showPGN: false,
      copied: false,
    };
  },
};
</script>

<style>
.board {
  transform: scale(2);
  transform-origin: 0 0;
  width: 700px;
  height: 700px;
}

.dimmed {
  position: absolute;
  width: 720px;
  height: 720px;
  text-align: center;
  font-size: 70pt;
  line-height: 700px;
  z-index: 10;
  background: rgba(0, 0, 0, 0.7);
}
</style>
