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
            <v-list-item v-for="(move, index) in movesWhite" :key="index"
              >{{ index + 1 }}:
              <v-icon
                :class="move.accuracy ? 'outlined' : ''"
                :color="getAccColor(move.accuracy)"
              >
                {{ getAccIcon(move.accuracy) }}</v-icon
              >
              &nbsp; {{ move.notation }}
            </v-list-item>
          </v-list>
        </v-col>
        <v-col cols="12" sm="6">
          <v-list dense>
            <v-list-item v-for="(move, index) in movesBlack" :key="index"
              ><v-icon
                :class="move.accuracy ? 'outlined' : ''"
                :color="getAccColor(move.accuracy)"
              >
                {{ getAccIcon(move.accuracy) }}</v-icon
              >
              &nbsp;{{ move.notation }}</v-list-item
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
    getAccIcon(acc) {
      if (acc == "Blunder") {
        return "fas fa-minus";
      }
      if (acc == "Inaccuracy") {
        return "fas fa-question";
      }
      if (acc == "Mistake") {
        return "fas fa-exclamation";
      }
      return "";
    },
    getAccColor(acc) {
      if (acc == "Blunder") {
        return "red";
      }
      if (acc == "Inaccuracy") {
        return "blue";
      }
      if (acc == "Mistake") {
        return "orange";
      }
      return "";
    },
  },
  data() {
    return {};
  },
};
</script>

<style>
.v-icon.outlined {
  border: 1px solid currentColor;
  border-radius: 50%;
  height: 30px;
  width: 30px;
}
</style>