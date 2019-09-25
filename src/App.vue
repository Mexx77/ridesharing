<template>
  <v-app>
    <v-app-bar app>
      <v-toolbar-title class="headline text-uppercase">
        <span v-if="$vuetify.breakpoint.mdAndUp">Sammatz </span>
        <span class="font-weight-light">RIDESHARING</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <div v-if="$store.state.account.status.loggedIn">
        Hi, {{$store.state.account.user.username}}!
        <v-btn
                icon
                text
                @click="logout"
        >
          <v-icon>mdi-logout-variant</v-icon>
        </v-btn>
      </div>
      <div v-else>
        <v-btn
                icon
                text
                @click="$store.commit('user/setShowLoginForm', true)"
        >
          <v-icon>mdi-login-variant</v-icon>
        </v-btn>
      </div>
    </v-app-bar>

    <v-content>
      <Calendar/>
      <LoginForm/>
      <v-snackbar v-model="snackbar" :color="alert.type">{{ alert.message }}</v-snackbar>
    </v-content>
    <v-footer>
      <v-spacer></v-spacer>
      <div>&copy; {{ new Date().getFullYear() }}</div>
    </v-footer>
  </v-app>
</template>

<script>
import Calendar from "./components/Calendar";
import LoginForm from "./components/LoginForm";
import {mapState,mapActions} from 'vuex';

export default {
  name: 'App',
  components: {
    Calendar,
    LoginForm
  },
  computed: {
    ...mapState({
      alert: state => state.alert
    }),
    snackbar: {
      get () {
        return this.$store.state.alert.visible
      },
      set (value) {
        this.$store.dispatch('alert/setVisibility', value)
      }
    }
  },
  methods: {
    ...mapActions('account', ['logout'])
  },
  mounted: function() {
    if (this.$store.state.account.status.loggedIn) {
      this.$store.dispatch('account/refreshToken')
    }
  }
};
</script>
