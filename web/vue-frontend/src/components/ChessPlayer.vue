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
            :readonly="locked"
          ></v-text-field>
        </v-row>
        <v-row>
          <v-select
            :items="players"
            item-text="name"
            item-value="value"
            v-model="defaultVal"
            v-on:input="$emit('modeChange', $event.value)"
            :readonly="locked"
            return-object
          ></v-select>
        </v-row>
        <v-row>
          <v-switch
            v-model="speak"
            label="Speak?"
            v-on:change="$emit('speakChange', $event)"
          />
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
      speak: false,
      defaultVal: {
        name: "Human",
        value: "0",
      },
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
  props: ["color", "locked"],
};
</script>
