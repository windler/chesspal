<template>
  <v-card variant="outlined" class="overflow-y-auto">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fas fa-forward</v-icon>
    </v-card-title>
    <v-container>      
      <v-row class="justify-center">
        <v-btn class="ma-2" @click="$emit('showHint', true)" outlined
          ><v-icon left>fas fa-circle-question</v-icon> Hint
        </v-btn>
      </v-row>
      <v-divider class="ma-4"></v-divider>
      <v-row class="justify-center">
        <v-btn class="ma-2" @click="$emit('undoMoves', 2)" outlined
          ><v-icon left> fas fa-rotate-left </v-icon> Undo
        </v-btn>
        <v-btn class="ma-2" @click="$emit('resign', 2)" outlined
          ><v-icon left> fas fa-flag-checkered </v-icon> Resign
        </v-btn>
        <v-btn class="ma-2" @click="$emit('draw', 2)" outlined
          >1/2 - 1/2
        </v-btn>
      </v-row>
      <v-divider class="ma-4"></v-divider>
      <v-row class="justify-center">
        <v-col cols="12" sm="9">
          <v-select
            class="ma-2"
            hint="Evaluation mode"
            persistent-hint
            :items="modes"
            item-text="name"
            item-value="value"
            v-model="defaultVal"
            v-on:input="$emit('changeMode', $event.value)"
            :readonly="locked"
            return-object
          ></v-select>
        </v-col>
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