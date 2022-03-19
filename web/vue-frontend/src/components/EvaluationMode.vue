<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fat fa-magnifying-glass</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-container fluid>
        <v-row>
          <v-select
            :items="modes"
            item-text="name"
            item-value="value"
            v-model="modes[0]"
            v-on:input="$emit('changeMode', $event)"
          ></v-select>
        </v-row>
        <v-row>
          <v-switch
            v-model="speak"
            label="Speak?"
            v-on:input="$emit('changeSpeech', $event)"
          />
        </v-row>
      </v-container>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "EvaluationMode",

  data() {
    return {
      speak: true,
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
