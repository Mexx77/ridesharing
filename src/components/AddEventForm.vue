<template>
    <v-row justify="center">
        <v-dialog
                v-model="showAddEventForm"
                :fullscreen="$vuetify.breakpoint.smAndDown ? true : false"
                hide-overlay max-width="600px"
        >
            <v-form
                    ref="form"
                    v-model="formIsValid"
                    lazy-validation
            >
                <v-card>
                    <v-toolbar color="primary" dark>
                        <v-toolbar-title>
                            <v-icon class="pb-1">mdi-car</v-icon>
                            Reservierung am {{germanDate}}
                        </v-toolbar-title>
                    </v-toolbar>
                    <v-card-text class="pb-0">
                        <v-container pa-0>
                            <v-row>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                            v-model="driver"
                                            prepend-icon="mdi-account"
                                            label="Fahrer*"
                                            :rules="[v => !!v || 'Name des Fahrers benötigt']"
                                            required
                                    ></v-text-field>
                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-combobox
                                            v-model="destination"
                                            prepend-icon="mdi-city"
                                            :items="['Lüneburg', 'Dannenberg', 'Hitzacker']"
                                            label="Fahrtziel*"
                                            :rules="[v => !!v || 'Fahrtziel benötigt']"
                                            required
                                    ></v-combobox>
                                </v-col>
                            </v-row>
                            <v-row>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-menu
                                            v-model="menuStartTime"
                                            :close-on-content-click="false"
                                            :nudge-right="40"
                                            transition="scale-transition"
                                            offset-y
                                            full-width
                                            min-width="272px"
                                    >
                                        <template v-slot:activator="{ on }">
                                            <v-text-field
                                                    v-model="startTime"
                                                    label="Startzeit"
                                                    prepend-icon="mdi-calendar-clock"
                                                    readonly
                                                    v-on="on"
                                                    :rules="[v => !!v || 'Startzeit benötigt']"
                                            ></v-text-field>
                                        </template>
                                        <v-time-picker
                                                v-model="startTime"
                                                color="primary"
                                                :width="272"
                                                format="24hr"
                                                required
                                                :allowed-minutes="allowedMinutes"
                                        ></v-time-picker>
                                    </v-menu>

                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-menu
                                            v-model="menuEndTime"
                                            :close-on-content-click="false"
                                            :nudge-right="40"
                                            transition="scale-transition"
                                            offset-y
                                            full-width
                                            min-width="272px"
                                    >
                                        <template v-slot:activator="{ on }">
                                            <v-text-field
                                                    v-model="endTime"
                                                    label="Zeit der Rückgabe"
                                                    prepend-icon="mdi-calendar-clock"
                                                    readonly
                                                    v-on="on"
                                                    :rules="[v => !!v || 'Zeit der Rückgabe benötigt']"
                                            ></v-text-field>
                                        </template>
                                        <v-time-picker
                                                v-model="endTime"
                                                color="primary"
                                                :width="272"
                                                format="24hr"
                                                :min="startTime"
                                                required
                                                :allowed-minutes="allowedMinutes"
                                        ></v-time-picker>
                                    </v-menu>
                                </v-col>
                            </v-row>
                        </v-container>
                    </v-card-text>
                    <v-card-actions class="mr-2 pb-4 pt-0">
                        <v-container pt-0>
                            <v-row dense>
                                <v-col>
                                    <v-switch
                                            class="mt-0"
                                            v-model="bigCarNeeded"
                                            label="Ich brauche ein großes Auto"
                                    ></v-switch>
                                    <v-btn @click="showAddEventForm = false">Abbrechen</v-btn>
                                    <v-btn @click="validateAndSubmitForm">Anfragen</v-btn>
                                </v-col>
                            </v-row>
                        </v-container>

                    </v-card-actions>
                </v-card>
            </v-form>
        </v-dialog>
        <v-snackbar v-model="snackbar" :color="snackbarColor">{{ snackbarText }}</v-snackbar>
    </v-row>
</template>

<script>
    export default {
        data() {
            return {
                driver: '',
                destination: '',
                formIsValid: false,
                menuStartTime: false,
                menuEndTime: false,
                bigCarNeeded: false,
                snackbar: false,
                snackbarText: '',
                snackbarColor: 'success'
            }
        },
        computed: {
            focus () {
                return this.$store.state.focus
            },
            germanDate() {
                if (this.focus == undefined) return undefined
                const p = this.focus.split(/\D/g)
                return [p[2], p[1], p[0]].join(".")
            },
            startTime: {
                get () {
                    return this.$store.state.startTime
                },
                set (value) {
                    this.$store.commit('setStartTime', value)
                }
            },
            endTime: {
                get () {
                    return this.$store.state.endTime
                },
                set (value) {
                    this.$store.commit('setEndTime', value)
                }
            },
            showAddEventForm: {
                get () {
                    return this.$store.state.showAddEventForm
                },
                set (value) {
                    this.$store.commit('setShowAddEventForm', value)
                }
            }
        },
        methods: {
            validateAndSubmitForm() {
                if (this.$refs.form.validate()) {
                    this.showAddEventForm = false;
                    this.$http
                        .post(this.$hostname + '/ride', {
                            driver: this.driver,
                            destination: this.destination,
                            start: `${this.focus}T${this.startTime}:00`,
                            end: `${this.focus}T${this.endTime}:00`,
                            bigCarNeeded: this.bigCarNeeded
                        })
                        .then(() => {
                            this.snackbarText = 'Danke, deine Reservierungsanfrage wurde entgegengenommen';
                            this.snackbarColor = 'success';
                            this.snackbar = true;
                        })
                        .catch((error) => {
                            // Error 😨
                            if (error.response) {
                                /*
                                 * The request was made and the server responded with a
                                 * status code that falls out of the range of 2xx
                                 */
                                this.snackbarText = 'Ups, der Server hat deine Anfrage verweigert';
                                this.snackbarColor = 'error';
                                this.snackbar = true;
                            } else if (error.request) {
                                /*
                                 * The request was made but no response was received, `error.request`
                                 * is an instance of XMLHttpRequest in the browser and an instance
                                 * of http.ClientRequest in Node.js
                                 */
                                this.snackbarText = 'Ups, keine Antort vom Server erhalten. Netz?';
                                this.snackbarColor = 'error';
                                this.snackbar = true;
                            } else {
                                // Something happened in setting up the request and triggered an Error
                                // eslint-disable-next-line
                                this.snackbarText = 'Ups, konnte die Anfrage nicht schicken :(';
                                this.snackbarColor = 'error';
                                this.snackbar = true;
                            }
                        });
                }
            },
            allowedMinutes: m => m % 15 == 0,
        }
    }
</script>