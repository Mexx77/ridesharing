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
        <v-icon>mdi-logout-variant</v-icon>
      </div>
      <div v-else>
        <v-btn
                icon
                text
                @click="$store.commit('users/setShowLoginForm', true)"
        >
          <v-icon>mdi-login-variant</v-icon>
        </v-btn>
      </div>
    </v-app-bar>

    <v-content>
      <Calendar/>
      <LoginForm/>
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

export default {
  name: 'App',
  components: {
    Calendar,
    LoginForm
  },
  data: () => ({
    //
  }),
};
</script>
