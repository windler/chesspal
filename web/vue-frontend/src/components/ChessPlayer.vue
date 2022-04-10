<template>
  <v-card shaped>
    <v-card-title primary-title class="justify-center">
      <v-icon>{{ icon }}</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-container fluid>
        <v-row>
          <v-text-field
            v-on:input="$emit('nameChange', $event)"
            label="Name"
            :value="name"
            :disabled="locked"
          ></v-text-field>

          <v-select
            :items="players"
            item-text="name"
            item-value="value"
            v-model="defaultVal"
            v-on:input="$emit('modeChange', $event.value)"
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
      players.push({
        name: "Human",
        value: 0,
      });
      for (var i = 0; i < this.bots.length; i++) {
        players.push({
          name: this.bots[i].name,
          value: i + 1,
        });
      }
      return players;
    },
  },
  data() {
    return {
      defaultVal: {
        name: "Human",
        value: 0,
      },
    };
  },

  props: ["color", "locked", "bots", "name"],
};
</script>
