<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fa fa-code</v-icon>
    </v-card-title>

    <v-card-actions>
      <v-btn color="primary" text> EXPLORE </v-btn>

      <v-spacer></v-spacer>

      <v-btn icon @click="show = !show">
        <v-icon>{{ show ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
      </v-btn>
    </v-card-actions>

    <v-expand-transition>
      <div v-show="show">
        <v-divider></v-divider>

        <v-card-text>
          {{ pgn }}
        </v-card-text>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script>
export default {
  name: "PGNCard",

  props: ["pgn"],
  data() {
    return {
      show: false
    };
  },
};
</script>
