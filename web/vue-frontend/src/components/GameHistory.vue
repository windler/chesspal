<template>
  <v-container fill-height>
    <v-row justify="center">
      <v-col cols="12" sm="12">
        <v-card>
          <v-data-table
            :headers="headers"
            :items="history"
            item-key="date"
            class="elevation-1"
            :search="search"
            :single-expand="singleExpand"
            show-expand
            :sort-by.sync="sortBy"
            :sort-desc.sync="sortDesc"
            @click:row="importLichess"
          >
            <template v-slot:top>
              <v-toolbar flat>
                <v-toolbar-title>Game History</v-toolbar-title>
                <v-spacer></v-spacer>
                <v-text-field
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search"
                  single-line
                  hide-details
                ></v-text-field>
                <v-btn icon @click="getGames()"
                  ><v-icon>fas fa-refresh</v-icon></v-btn
                >
              </v-toolbar>
            </template>
            <template v-slot:[`item.actions`]="{ item }">
              <v-icon class="mr-2" @click="importLichess(item)">
                fas fa-magnifying-glass-chart
              </v-icon>
              <!-- <v-icon small class="mr-2" @click="archiveGame(item.id)">
                fas fa-box-open
              </v-icon>
              <v-icon small class="mr-2" @click="deleteGame(item.id)">
                fas fa-trash-can
              </v-icon> -->
            </template>
            <template v-slot:expanded-item="{ headers, item }">
              <td :colspan="headers.length">
                <div class="ma-6 text-center board-xsmall" v-html="item.svg" />
              </td>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: "GameHistory",

  data() {
    return {
      history: [],
      search: "",
      singleExpand: true,
      sortBy: "dateTime",
      sortDesc: true,
    };
  },

  computed: {
    headers() {
      return [
        {
          text: "DateTime",
          value: "dateTime",
          align: " d-none",
        },
        {
          text: "Date",
          align: "start",
          sortable: true,
          value: "date",
        },
        {
          text: "White",
          sortable: true,
          value: "white",
        },
        {
          text: "Black",
          sortable: true,
          value: "black",
        },
        {
          text: "Result",
          sortable: true,
          value: "result",
        },
        {
          text: "Archived",
          value: "archived",
          align: " d-none",
          filter: this.archiveFitler,
        },
        {
          text: "Botgame",
          value: "botgame",
          align: " d-none",
          filter: this.botFitler,
        },
        { text: "Actions", value: "actions", sortable: false },
        { text: "", value: "data-table-expand" },
      ];
    },
  },

  methods: {
    archiveFitler: function (item) {
      if (this.showArchived === "true") {
        return item;
      }

      return !item;
    },
    botFitler: function (item) {
      if (this.showBotGames === "false") {
        if (this.showHumanGames === "true") {
          return !item;
        }
      }

      if (this.showBotGames === "true") {
        if (this.showHumanGames === "false") {
          return item;
        }
      }

      return true;
    },
    getHost: function() {
      var host = location.host
      if (process.env.VUE_APP_CHESSPAL_HOST !== undefined) {
        host = process.env.VUE_APP_CHESSPAL_HOST
      }
      return host
    },
    getGames: function () {
      fetch("http://" + this.getHost() + "/history")
        .then((response) => response.json())
        .then((data) => (this.history = data.games));
    },
    deleteGame: function (id) {
      fetch("http://" + this.getHost() + "/history/" + id, { method: "DELETE" });
      this.getGames();
    },
    archiveGame: function (id) {
      fetch("http://" + this.getHost() + "/history/" + id + "/archive", {
        method: "POST",
      });
      this.getGames();
    },
    importLichess: async function (row) {
      var win = window.open('', '_blank');
      const requestOptions = {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: "pgn=" + row.pgn,
      };
      const response = await fetch(
        "https://lichess.org/api/import",
        requestOptions
      );
      const data = await response.json();
      win.location = data.url
    },
  },

  created: function () {
    this.getGames();
  },

  props: ["showArchived", "showBotGames", "showHumanGames"],
};
</script>
<style>
.board-xsmall {
  transform: scale(0.8);
  transform-origin: 0 0;
  width: 280px;
  height: 280px;
}
</style>
