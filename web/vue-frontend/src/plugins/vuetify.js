import '@fortawesome/fontawesome-free/css/all.css' // Ensure you are using css-loader
import Vue from 'vue'
import Vuetify from 'vuetify/lib'

Vue.use(Vuetify)

export default new Vuetify({
  icons: {
    iconfont: 'fa',
  },
  theme: {
    themes: {
      light: {
        primary: '#1867C0',
        secondary: '#5CBBF6',
        tertiary: '#E57373',
        accent: '#005CAF',
      },
    },
  },
})