<template>
  <v-card shaped>
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fas fa-forward</v-icon>
    </v-card-title>
    <v-container>
      <v-row class="justify-center">
        <v-btn class="ma-4" @click="$emit('showHint', true)" outlined
          ><v-icon left>fas fa-circle-question</v-icon> Hint
        </v-btn>
      </v-row>
      <v-row class="justify-center">
        <v-btn class="ma-4" @click="$emit('undoMoves', 2)" outlined
          ><v-icon left> fas fa-rotate-left </v-icon> Undo
        </v-btn>
        <v-btn class="ma-4" @click="$emit('resign', 2)" outlined
          ><v-icon left> fas fa-flag-checkered </v-icon> Resign
        </v-btn>
        <v-btn class="ma-4" @click="$emit('draw', 2)" outlined
          >1/2 - 1/2
        </v-btn>
      </v-row>
      <v-row class="justify-center">
        <v-select
          class="ma-4"
          label="Evaluation mode"
          persistent-hint
          :items="modes"
          item-text="name"
          item-value="value"
          v-model="defaultVal"
          v-on:input="$emit('changeMode', $event.value)"
          return-object
        ></v-select>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>
export default {
  name: "GameActions",

  data() {
    return {
      defaultVal: {
        name: "None",
        value: "0",
      },
      modes: [
        {
          name: "None",
          value: "0",
        },
        {
          name: "Move Feedback",
          value: "1",
        },
      ],
    };
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