<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fat fa-gear</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-container fluid>
        <v-row>
          <v-switch
            v-model="upsideDown"
            label="Invert sides of DGT board?"
            v-on:change="$emit('upsideDownChange', $event)"
            :readonly="locked"
          />
        </v-row>
      </v-container>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: "SettingsCard",
  props: ["locked"],

  data() {
    return {
      upsideDown: false,
    };
  },
};
</script>
