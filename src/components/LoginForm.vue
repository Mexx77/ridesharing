<template>
    <v-row justify="center">
        <v-dialog
            v-model="showLoginForm"
            :fullscreen="$vuetify.breakpoint.smAndDown"
            hide-overlay max-width="400px"
        >
            <v-form
                ref="form"
                v-model="formIsValid"
                lazy-validation
                @keyup.native.enter="handleLogin"
            >
                <v-card>
                    <v-toolbar color="primary" dark>
                        <v-toolbar-title>
                            <v-icon class="pb-1">mdi-account</v-icon>
                            Login
                        </v-toolbar-title>
                    </v-toolbar>
                    <v-card-text class="pb-0">
                        <v-container pa-0>
                            <v-row dense>
                                <v-col cols="12">
                                    <v-text-field
                                        v-model="usernamePhone"
                                        prepend-icon="mdi-account"
                                        label="Handy-Nr. oder Benutzername*"
                                        :rules="[v => !!v || 'benötigt']"
                                        required
                                        data-cy="username-phone"
                                    ></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field
                                        v-model="password"
                                        prepend-icon="mdi-lock-question"
                                        :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                                        :rules="[v => !!v || 'benötigt']"
                                        :type="showPassword ? 'text' : 'password'"
                                        label="Passwort*"
                                        @click:append="showPassword = !showPassword"
                                        data-cy="password"
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions class="mr-2 pb-4 pt-0">
                        <v-container class="pt-0 pr-1 pl-1">
                            <v-row dense>
                                <v-col>
                                    <v-btn
                                        text
                                        color="primary"
                                        @click="handleRegister"
                                    >Registrieren
                                    </v-btn>
                                </v-col>
                                <v-spacer></v-spacer>
                                <v-col>
                                    <v-btn @click="showLoginForm = false">Abbrechen</v-btn>
                                </v-col>
                                <v-col>
                                    <v-btn @click="handleLogin" :disabled="status.loggingIn" data-cy="login-btn">
                                        <v-progress-circular color="primary" v-if="status.loggingIn" class="mr-1"
                                                             size="12" width="2" indeterminate></v-progress-circular>
                                        Login
                                    </v-btn>
                                </v-col>
                            </v-row>
                        </v-container>

                    </v-card-actions>
                </v-card>
            </v-form>
        </v-dialog>
    </v-row>
</template>

<script>
  import {mapState, mapActions} from 'vuex'

  export default {
    data() {
      return {
        usernamePhone: '',
        password: '',
        showPassword: false,
        formIsValid: false
      }
    },
    computed: {
      ...mapState('account', ['status']),
      showLoginForm: {
        get() {
          return this.$store.state.account.showLoginForm
        },
        set(v) {
          this.$store.commit('account/showLoginForm', v)
        }
      }
    },
    methods: {
      ...mapActions('account', ['login']),
      handleLogin() {
        if (this.$refs.form.validate()) {
          const {usernamePhone, password} = this;
          this.login({usernamePhone, password}).then(this.resetData)
        }
      },
      resetData() {
        Object.assign(this.$data, this.$options.data.apply(this))
      },
      handleRegister() {
        this.showLoginForm = false;
        this.$store.dispatch('account/showRegisterForm', true)
      }
    }
  }
</script>