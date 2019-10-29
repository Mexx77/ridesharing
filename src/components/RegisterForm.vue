<template>
    <v-row justify="center">
        <v-dialog
            v-model="showRegisterForm"
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
                            Registrieren
                        </v-toolbar-title>
                    </v-toolbar>
                    <v-card-text class="pb-0">
                        <v-container pa-0>
                            <v-row dense>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="vorname"
                                        prepend-icon="mdi-account-card-details-outline"
                                        label="Vorname*"
                                        :rules="[rules.required]"
                                        required
                                    ></v-text-field>
                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="nachname"
                                        prepend-icon="mdi-account-card-details-outline"
                                        label="Nachname*"
                                        :rules="[rules.required]"
                                        required
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row dense>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="username"
                                        prepend-icon="mdi-account"
                                        label="Benutzername*"
                                        :rules="[rules.required,rules.min3]"
                                        required
                                    ></v-text-field>
                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        v-model="password"
                                        prepend-icon="mdi-lock-question"
                                        :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                                        hint="1 Kleinbuchstaben, 1 Großbuchstaben, 1 Zahl, Länge 8"
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
                                        v-model="handy"
                                        prepend-icon="mdi-cellphone-iphone"
                                        label="Handy-Nr.*"
                                        hint="Für Fahrt-Bestätigungen per SMS"
                                        persistent-hint
                                        :rules="[rules.required]"
                                        required
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions class="mr-2 pb-4 pt-3">
                        <v-container pt-0>
                            <v-row dense>
                                <v-spacer></v-spacer>
                                <v-col>
                                    <v-btn @click="showRegisterForm = false">Abbrechen</v-btn>
                                </v-col>
                                <v-col>
                                    <v-btn >Registrieren</v-btn>
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
    import { mapState } from 'vuex'

    export default {
        data() {
            return {
                vorname: '',
                nachname: '',
                username: '',
                password: '',
                handy: '',
                showPassword: false,
                submitted: false,
                formIsValid: false,
                rules: {
                    required: value => !!value || 'Benötigt',
                    min3: value => value.length >= 3 || 'mind. 3 Zeichen',
                    password: value => {
                        const pattern = new RegExp("^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.{8,})");
                        return pattern.test(value) || '1 Kleinbuchstaben, 1 Großbuchstaben, 1 Zahl, Länge 8'
                    },
                },
            }
        },
        computed: {
            ...mapState('account', ['status']),
            showRegisterForm: {
                get () {
                    return this.$store.state.user.showRegisterForm
                },
                set (v) {
                    this.$store.commit('user/showRegisterForm', v)
                }
            }
        },
        methods: {

        }
    }
</script>