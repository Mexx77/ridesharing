<template>
  <v-app>
    <v-app-bar app>
      <v-toolbar-title class="headline text-uppercase">
        <span v-if="$vuetify.breakpoint.mdAndUp">{{brandName}} </span>
        <span class="font-weight-light">RIDESHARING</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <div v-if="$store.state.account.status.loggedIn">
        Hi, {{$store.state.account.user.firstName}}!
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
            text
            color="secondary"
            @click="$store.dispatch('user/showLoginForm', true)"
        >
          <v-icon>mdi-login-variant</v-icon>&nbsp;Anmelden
        </v-btn>
      </div>
    </v-app-bar>

    <v-content>
      <Calendar/>
      <LoginForm/>
      <RegisterForm/>
      <v-snackbar v-model="snackbar" :timeout="timeout" :color="alert.type">
        {{ alert.message }}
        <v-btn
            dark
            text
            @click="snackbar = false"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-snackbar>
    </v-content>
<!--    <v-footer>-->
<!--      <v-spacer></v-spacer>-->
<!--      <div>&copy; {{ new Date().getFullYear() }}</div>-->
<!--    </v-footer>-->
  </v-app>
</template>

<script>
import Calendar from "./components/Calendar";
import LoginForm from "./components/LoginForm";
import RegisterForm from "./components/RegisterForm";
import {mapState,mapActions} from 'vuex';
import * as constants from "./_services/constants";

export default {
  name: 'App',
  components: {
    Calendar,
    LoginForm,
    RegisterForm
  },
  data() {
    return {
      brandName: constants.brandName,
    }
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
    },
    timeout: {
      get () {
        return this.$store.state.alert.timeout
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
    this.$store.dispatch('alert/info', {
      message: '💡 Um eine Fahrt hinzuzufügen, klicke neben die ungefähre Startzeit',
      timeout: 0
    })
  }
};
</script>
<style>
  .theme--light.v-application.v-application--is-ltr {
    background: white;
  }
  body {
    overflow: hidden;
    overflow-y: auto;
  }
</style>