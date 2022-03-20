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
            v-model="defaultVal"
            v-on:input="$emit('changeMode', $event.value)"
            :readonly="locked"
            return-object
          ></v-select>
        </v-row>
      </v-container>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "EvaluationMode",
  props: ["locked"],

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
