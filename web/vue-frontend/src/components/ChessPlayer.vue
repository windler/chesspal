<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon>{{ icon }}</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-container fluid>
        <v-row>
          <v-text-field
            v-on:input="$emit('nameChange', $event)"
            label="Name"
            v-bind:value="color.toUpperCase()"
          ></v-text-field>
        </v-row>
        <v-row>
          <v-select
            :items="players"
            item-text="name"
            item-value="value"
            v-model="players[0]"
            v-on:input="$emit('modeChange', $event)"
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
  },
  data() {
    return {
      players: [
        {
          name: "Human",
          value: "0",
        },
      ],
    };
  },
  created: function () {
    for (var i = 0; i < 20; i++) {
      this.players.push({
        name: "AI skill " + (i + 1),
        value: i + 1,
      });
    }
  },
  props: ["color"],
};
</script>
