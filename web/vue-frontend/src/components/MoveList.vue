<template>
  <v-card variant="outlined" :min-height="height" :max-height="height" >
    <div id="moveListContainer" class="overflow-y-auto" style="height: 300px">
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
                  :class="move.accuracy && showEvaluation ? 'outlined' : ''"
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
                  :class="move.accuracy && showEvaluation ? 'outlined' : ''"
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
    </div>
  </v-card>
</template>

<script>
export default {
  name: "MoveList",
  props: ["movesBlack", "movesWhite", "showEvaluation", "height"],
  watch: {
    movesWhite: function () {
      var container = this.$el.querySelector("#moveListContainer");
      container.scrollTop = container.scrollHeight;
    },
  },
  methods: {
    getAccIcon(acc) {
      if (!this.showEvaluation) {
        return "";
      }

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
      if (!this.showEvaluation) {
        return "";
      }

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