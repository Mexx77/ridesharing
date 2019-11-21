<template>
    <v-row justify="center">
        <v-dialog
            v-model="showRegisterForm"
            :fullscreen="$vuetify.breakpoint.smAndDown"
            hide-overlay max-width="700px"
        >
            <v-form
                ref="form"
                v-model="formIsValid"
                lazy-validation
                @keyup.native.enter="handleRegister"
            >
                <v-card>
                    <v-toolbar color="primary" dark>
                        <v-toolbar-title>
                            <v-icon class="pb-1">mdi-account</v-icon>
                            Registrieren
                        </v-toolbar-title>
                    </v-toolbar>
                    <v-card-text class="pb-0">
                        <v-container pa-0>
                            <v-row dense>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="user.firstName"
                                        prepend-icon="mdi-account-card-details-outline"
                                        label="Vorname*"
                                        :rules="[rules.min2]"
                                    ></v-text-field>
                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="user.lastName"
                                        prepend-icon="mdi-account-card-details-outline"
                                        label="Nachname*"
                                        :rules="[rules.min2]"
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row dense>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="user.phone"
                                        prepend-icon="mdi-cellphone-iphone"
                                        label="dt. Handy-Nr.*"
                                        hint="Für Bestätigungen deiner Auto-Anfragen per SMS"
                                        persistent-hint
                                        :rules="[rules.phone]"
                                    ></v-text-field>
                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="user.password"
                                        prepend-icon="mdi-lock-question"
                                        :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                                        :hint="passwordHint"
                                        persistent-hint
                                        :rules="[rules.password]"
                                        :type="showPassword ? 'text' : 'password'"
                                        label="Passwort*"
                                        @click:append="showPassword = !showPassword"
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row dense>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="user.username"
                                        prepend-icon="mdi-account"
                                        label="Benutzername"
                                        hint="Kann beim Login statt der Handy-Nr. verwendet werden"
                                        persistent-hint
                                        :rules="[rules.min3orEmpty]"
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions class="mr-2 pb-4 pt-3">
                        <v-container pt-0>
                            <v-row dense>
                                <v-col>
                                    <v-btn @click="showRegisterForm = false">Abbrechen</v-btn>
                                    <v-btn @click="handleRegister" :disabled="status.registering">
                                        <v-progress-circular color="primary" v-if="status.registering"
                                                             size="12" width="2" indeterminate></v-progress-circular>
                                        <span v-if="status.registering">&nbsp;</span>
                                        Registrieren
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
        user: {
          firstName: '',
          lastName: '',
          username: '',
          password: '',
          phone: ''
        },
        showPassword: false,
        formIsValid: false,
        rules: {
          required: value => !!value || 'Benötigt',
          min2: value => value.length >= 2 || 'mind. 2 Zeichen',
          min3orEmpty: value => (value === '' || value.length >= 3) || 'leer oder mind. 3 Zeichen',
          password: value => {
            const pattern = new RegExp("^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9]).{8,}$");
            return pattern.test(value) || this.passwordHint
          },
          phone: value => {
            const pattern = new RegExp("^01[567][0-9]{8,11}$");
            return pattern.test(value) || 'dt. Handynummer im Format 01712345678'
          },
        },
        passwordHint: '1 Kleinbuchstabe, 1 Großbuchstabe, 1 Zahl, Länge ≥ 8'
      }
    },
    computed: {
      ...mapState('account', ['status']),
      showRegisterForm: {
        get() {
          return this.$store.state.account.showRegisterForm
        },
        set(v) {
          this.$store.commit('account/showRegisterForm', v)
        }
      }
    },
    methods: {
      ...mapActions('account', ['register']),
      handleRegister() {
        if (this.$refs.form.validate()) {
          this.register(this.user)
        }
      },
    }
  }
</script>