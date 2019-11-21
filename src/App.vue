<template>
  <v-app>
    <v-app-bar app>
      <v-icon class="ml-3" large>mdi-bus-school</v-icon>
      <v-toolbar-title class="headline text-uppercase mr-4 ml-3" v-if="$vuetify.breakpoint.mdAndUp">
        <span class="font-weight-light">RIDESHARING</span>
      </v-toolbar-title>
      <v-menu bottom right>
        <template v-slot:activator="{ on }">
          <v-btn text v-on="on">
            <span>{{ typeToLabel[type] }}</span>
            <v-icon right>mdi-menu-down</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="type = 'day'">
            <v-list-item-title>Tag</v-list-item-title>
          </v-list-item>
          <v-list-item @click="type = '4day'">
            <v-list-item-title>4 Tage</v-list-item-title>
          </v-list-item>
          <v-list-item @click="type = 'week'" v-if="$vuetify.breakpoint.mdAndUp">
            <v-list-item-title>Woche</v-list-item-title>
          </v-list-item>
          <v-list-item @click="type = 'month'" v-if="$vuetify.breakpoint.mdAndUp">
            <v-list-item-title>Monat</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-spacer></v-spacer>
      <v-btn fab text @click="prev">
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>
      <v-btn text width="30" @click="focusToday">
        Heute
      </v-btn>
      <v-btn fab text @click="next">
        <v-icon>mdi-chevron-right</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content v-touch:swipe.left="next" v-touch:swipe.right="prev">
      <Calendar ref="calendar"></Calendar>
      <LoginForm/>
      <RegisterForm/>
      <v-snackbar top v-model="snackbar" :timeout="timeout" :color="alert.type">
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
    <v-footer fixed>
      <v-spacer></v-spacer>
      <div v-if="$store.state.account.status.loggedIn">
          <span style="vertical-align: text-top">Hi, {{$store.state.account.user.firstName}}!</span>
          <v-btn class="pr-2 pl-2 pt-0" text color="secondary" @click="logout">
            <v-icon>mdi-logout-variant</v-icon> Logout
          </v-btn>
      </div>
      <div v-else>
        <v-btn class="pa-2" text color="secondary" @click="$store.dispatch('account/showLoginForm', true)">
          <v-icon>mdi-login-variant</v-icon>&nbsp;Anmelden
        </v-btn>
      </div>
    </v-footer>
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
      typeToLabel: {
        month: 'Monat',
        week: 'Woche',
        day: 'Tag',
        '4day': '4 Tage',
      },
    }
  },
  computed: {
    ...mapState('calendar', ['today']),
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
    },
    type: {
      get () {
        return this.$store.state.calendar.type
      },
      set (value) {
        this.$store.commit('calendar/setType', value)
      }
    },
    focus: {
      get () {
        return this.$store.state.calendar.focus
      },
      set (value) {
        this.$store.commit('calendar/setFocus', value)
      }
    }
  },
  methods: {
    ...mapActions('account', ['logout']),
    prev() {
      this.$refs.calendar.prev()
    },
    next() {
      this.$refs.calendar.next()
    },
    focusToday() {
      this.focus = this.today
    },
  },
  mounted: function() {
    if (this.$store.state.account.status.loggedIn) {
      this.$store.dispatch('account/refreshToken')
    }
    if (!this.$store.state.account.user) {
      this.$store.dispatch('alert/info', {
        message: '💡 Um eine Fahrt hinzuzufügen, tippe neben die ungefähre Startzeit (dafür musst du angemeldet sein)',
        timeout: 21000
      })
    }

  }
};
</script>
<style>
  .v-application--wrap .v-toolbar__content{
    padding: 0 5px 0 15px;
    background: #fafafa;
  }
  .theme--light.v-application.v-application--is-ltr {
    background: white;
  }
  body {
    overflow: hidden;
    overflow-y: auto;
  }
</style>
