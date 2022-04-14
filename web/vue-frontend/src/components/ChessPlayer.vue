<template>
  <v-card>
    <v-card-title primary-title class="justify-center">
      <v-icon>{{ icon }}</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-container fluid>
        <v-row>
          <v-select
            :items="players"
            item-text="name"
            item-value="value"
            v-on:input="$emit('modeChange', $event)"
            :disabled="locked"
            return-object
          ></v-select>
        </v-row>
      </v-container>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "ChessPlayer",
  computed: {
    icon() {
      return this.color == "black" ? "fa fa-chess-king" : "far fa-chess-king";
    },
    players() {
      var players = [];

      var val = 0
      for (var i = 0; i < this.humans.length; i++) {
        players.push({
          name: this.humans[i].name,
          isHuman: true,
          value: val,
          mode: i,
        });
        val++
      }
      for (var j = 0; j < this.bots.length; j++) {
        players.push({
          name: this.bots[j].name,
          isHuman: false,
          value: val,
          mode: j,
        });
        val++
      }
      return players;
    },
  },
  data() {
    return {};
  },

  props: ["color", "locked", "bots", "humans", "name"],
};
</script>
