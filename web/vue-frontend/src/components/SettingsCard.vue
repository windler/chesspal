<template>
  <v-card>
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fat fa-gear</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-container>
        <v-row>
          <v-col col="6">
            <v-switch
              :value="speak"
              label="Speak"
              v-on:change="$emit('speakChange', $event)"
            />
          </v-col>
          <v-col col="6">
            <v-switch
              v-model="upsideDown"
              label="Invert DGT"
              v-on:change="$emit('upsideDownChange', $event)"
              :disabled="locked"
            />
          </v-col> </v-row
        >
      </v-container>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "SettingsCard",
  props: ["locked", "speak"],

  data() {
    return {
      upsideDown: false,
    };
  },
};
</script>
