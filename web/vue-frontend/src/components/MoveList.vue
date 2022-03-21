<template>
  <v-card
    variant="outlined"
    min-height="400px"
    max-height="400px"
    class="overflow-y-auto"
  >
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fa fa-list</v-icon>
    </v-card-title>

    <v-container>
      <v-row class="justify-center">
        <v-col cols="12" sm="6">
          <v-list dense>
            <v-list-item
              v-for="(move, index) in movesWhite"
              :key="index"
              >{{ index + 1 }}:
              <span :class="getTextColor(move.accuracy)">
                {{ move.notation }}</span
              ></v-list-item
            >
          </v-list>
        </v-col>
        <v-col cols="12" sm="6">
          <v-list dense>
            <v-list-item v-for="(move, index) in movesBlack" :key="index"
              ><span :class="getTextColor(move.accuracy)">
                {{ move.notation }}</span
              ></v-list-item
            >
          </v-list>
        </v-col>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>
export default {
  name: "MoveList",
  props: ["movesBlack", "movesWhite"],
  methods: {
    getTextColor(acc) {
      if (acc == "Blunder") {
        return "red--text";
      }
      if (acc == "Inaccuracy") {
        return "blue--text";
      }
      if (acc == "Mistake") {
        return "orange--text";
      }
      return "";
    },
  },
  data() {
    return {};
  },
};
</script>
