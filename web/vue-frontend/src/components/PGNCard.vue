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
        <v-btn class="ma-2" icon @click="importLichess()"
          ><v-icon>fas fa-magnifying-glass-chart</v-icon>
        </v-btn>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script>
export default {
  name: "PGNCard",

  props: ["pgn"],
  methods: {
    importLichess: async function () {
      const requestOptions = {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: "pgn=" + this.pgn,
      };
      const response = await fetch(
        "https://lichess.org/api/import",
        requestOptions
      );
      const data = await response.json()
      window.open(data.url, "_blank");
    },
  },
  data() {
    return {
      show: false,
    };
  },
};
</script>
