<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fa fa-flask</v-icon>
    </v-card-title>

    <v-progress-linear
      height="50"
      :active="show"
      :value="getPawnValue()"
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

  methods: {
    getPawnValue() {
      let base = 50;
      if (this.pawn == base) {
        return base;
      }
      let advantage = Math.abs(base - this.pawn);
      let translatedAdvantage = base - base / advantage;
      if (advantage < 1) {
        translatedAdvantage = advantage;
      }
      if (this.pawn < base) {
        return base - translatedAdvantage;
      }
      return base + translatedAdvantage;
    },
  },

  props: ["pawn", "show"],
  data() {
    return {};
  },
};
</script>
