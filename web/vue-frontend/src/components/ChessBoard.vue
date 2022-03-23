<template>
  <v-card variant="outlined" min-width="720px" min-height="700px">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fas fa-chess</v-icon>
    </v-card-title>

    <div v-html="svg" class="board">
      {{ svg }}
    </div>
  </v-card>
</template>

<script>
export default {
  name: "ChessBoard",

  props: ["svg", "fen"],
  data() {
    return {};
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
</style>
