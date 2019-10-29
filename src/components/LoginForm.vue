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
                                        v-model="username"
                                        prepend-icon="mdi-account"
                                        label="Benutzername*"
                                        :rules="[v => !!v || 'Benutzername benötigt']"
                                        required
                                    ></v-text-field>
                                </v-col>
                                <v-col cols="12">
                                    <v-text-field
                                        v-model="password"
                                        prepend-icon="mdi-lock-question"
                                        :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                                        :rules="[v => !!v || 'Passwort benötigt']"
                                        :type="showPassword ? 'text' : 'password'"
                                        label="Passwort*"
                                        @click:append="showPassword = !showPassword"
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions class="mr-2 pb-4 pt-0">
                        <v-container pt-0>
                            <v-row dense>
                                <v-col>
                                    <v-btn
                                        outlined
                                        @click="handleRegister"
                                    >Registrieren</v-btn>
                                </v-col>
                                <v-spacer></v-spacer>
                                <v-col>
                                    <v-btn @click="showLoginForm = false">Abbrechen</v-btn>
                                </v-col>
                                <v-col>
                                    <v-btn @click="handleLogin">Login</v-btn>
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
    import { mapState, mapActions } from 'vuex'

    export default {
        data() {
            return {
                username: '',
                password: '',
                showPassword: false,
                submitted: false,
                formIsValid: false
            }
        },
        computed: {
            ...mapState('account', ['status']),
            showLoginForm: {
                get () {
                    return this.$store.state.user.showLoginForm
                },
                set (v) {
                    this.$store.commit('user/showLoginForm', v)
                }
            }
        },
        methods: {
            ...mapActions('account', ['login']),
            handleLogin() {
                if (this.$refs.form.validate()) {
                    const { username, password } = this;
                    this.login({ username, password })
                }
            },
            handleRegister() {
                this.showLoginForm = false;
                this.$store.dispatch('user/showRegisterForm', true)
            }
        }
    }
</script>