<template>
  <!-- <v-card variant="outlined" min-width="720px" min-height="700px"> -->
  <v-card darkrounded min-width="280px" min-height="280px">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fas fa-chess-board</v-icon>
    </v-card-title>

    <div>
      <v-overlay :absolute="absolute" :value="outcome != '*'" class="dimmed white--text">
        {{ outcome != "*" ? outcome : "" }}
      </v-overlay>
      
      <div v-html="svg" :class="boardClass() + ' ma-auto'"></div>
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
    boardClass: function() {
      switch (this.$vuetify.breakpoint.name) {
        case "xs": return 'board-xsmall'
        case "sm": return 'board-small'
        case "md": return 'board-small'
        case "lg": return 'board-medium'
        case "xl": return 'board'
      }
       return 'board'
      
    },
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
      absolute: true,
    };
  },
};
</script>

<style>

.board {
    transform: scale(2);
    transform-origin: 0 0;
    width: 720px; 
    height: 720px;
}

.board-xsmall {
  transform: scale(0.8);
    transform-origin: 0 0;
    width: 280px; 
    height: 280px;
}

.board-small {
    transform: scale(1);
    transform-origin: 0 0;
    width: 350px; 
    height: 350px;
}

.board-medium {
    transform: scale(1.5);
    transform-origin: 0 0;
    width: 525px; 
    height: 525px;
}

.dimmed {
  text-align: center;
  font-size: 70pt;
  line-height: 100%;
}
</style>
