<template>
    <v-row justify="center">
        <v-dialog
            v-model="showAddEventForm"
            :fullscreen="$vuetify.breakpoint.smAndDown"
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
                            <span v-if="isUpdate"> Fahrt {{selectedEvent.name}}</span>
                            <span v-else> Reservierung am {{germanDate}}</span>
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
                            <v-row v-if="isAdmin">
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-select
                                        v-model="carName"
                                        prepend-icon="mdi-car"
                                        label="Auto*"
                                        :items="cars"
                                    ></v-select>
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
                                    <v-btn @click="validateAndSubmitForm">{{isAdmin ? 'Speichern' : 'Anfragen'}}</v-btn>
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
    import {mapActions} from 'vuex'

    export default {
        data() {
            return {
                formIsValid: false,
                menuStartTime: false,
                menuEndTime: false,
                cars: ['Red Bus', 'White Bus', 'e-Auto', 'Little Red']
            }
        },
        computed: {
            focus () {
                return this.$store.state.ride.focus
            },
            germanDate() {
                if (this.focus == undefined) return undefined
                const p = this.focus.split(/\D/g)
                return [p[2], p[1], p[0]].join(".")
            },
            startTime: {
                get () {
                    return this.$store.state.ride.startTime
                },
                set (value) {
                    this.$store.commit('ride/setStartTime', value)
                }
            },
            endTime: {
                get () {
                    return this.$store.state.ride.endTime
                },
                set (value) {
                    this.$store.commit('ride/setEndTime', value)
                }
            },
            showAddEventForm: {
                get () {
                    return this.$store.state.ride.showAddEventForm
                },
                set (value) {
                    this.$store.commit('ride/showAddUpdateRideForm', value)
                }
            },
            driver: {
                get () {
                    return this.$store.state.ride.driver
                },
                set (value) {
                    this.$store.commit('ride/setDriver', value)
                }
            },
            destination: {
                get () {
                    return this.$store.state.ride.destination
                },
                set (value) {
                    this.$store.commit('ride/setDestination', value)
                }
            },
            bigCarNeeded: {
                get () {
                    return this.$store.state.ride.bigCarNeeded
                },
                set (value) {
                    this.$store.commit('ride/setBigCarNeeded', value)
                }
            },
            selectedEvent: {
                get () {
                    return this.$store.state.ride.selectedEvent
                }
            },
            isUpdate: {
                get () {
                    return this.$store.state.ride.isUpdate
                }
            },
            isAdmin: function () {
                return this.$store.state.account.status.loggedIn && this.$store.state.account.user.isAdmin
            },
            carName: {
                get () {
                    return this.$store.state.ride.carName
                },
                set (value) {
                    this.$store.commit('ride/setCarName', value)
                }
            },
        },
        methods: {
            ...mapActions('ride', ['addRide', 'updateRide']),
            validateAndSubmitForm() {
                if (this.$refs.form.validate()) {
                    if (this.isUpdate) {
                        this.updateRide()
                    } else {
                        this.addRide()
                    }
                }
            },
            allowedMinutes: m => m % 15 == 0,
        }
    }
</script>