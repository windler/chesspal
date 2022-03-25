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
  </v-card>
</template>

<script>
export default {
  name: "ChessBoard",

  props: ["svg", "fen", "outcome"],
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
