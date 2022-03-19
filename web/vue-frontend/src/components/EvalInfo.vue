<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fa fa-flask</v-icon>
    </v-card-title>

    <v-progress-linear
      height="50"
      :value="pawn"
      background-color="grey darken-4"
      color="grey lighten"
    >
      <span class="white--text">
        {{ Number(pawn - 50).toLocaleString() }}
      </span>
    </v-progress-linear>
  </v-card>
</template>

<script>
export default {
  name: "EvalInfo",

  props: ["pawn"],
  data() {
    return {};
  },
};
</script>
