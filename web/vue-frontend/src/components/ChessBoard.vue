<template>
  <v-card variant="outlined">
    <v-card-title primary-title class="justify-center">
      <v-icon color="grey">fas fa-chess</v-icon>
    </v-card-title>

    <div v-html="svg">
      {{ svg }}
    </div>
  </v-card>
</template>

<script>
export default {
  name: "ChessBoard",

  props: ["svg"],
  data() {
    return {};
  },
};
</script>
