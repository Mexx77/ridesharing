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
                            <span v-if="isUpdate"> Fahrt {{driver}} ↦ {{destination}}</span>
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
                                    <v-text-field
                                        prepend-icon="mdi-calendar-clock"
                                        label="Startzeit"
                                        v-model="startTime"
                                        color="primary"
                                        type="time"
                                        step="600"
                                        suffix="Uhr"
                                        format="24hr"
                                        :rules="[v => !!v || 'Startzeit benötigt']"
                                    ></v-text-field>
                                </v-col>
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-text-field
                                        prepend-icon="mdi-calendar-clock"
                                        label="Zeit der Rückgabe"
                                        v-model="endTime"
                                        color="primary"
                                        type="time"
                                        step="600"
                                        suffix="Uhr"
                                        format="24hr"
                                        :rules="[v => !!v || 'Zeit der Rückgabe benötigt']"
                                    ></v-text-field>
                                </v-col>
                            </v-row>
                            <v-row v-if="isAdmin">
                                <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                    <v-select
                                        v-model="carName"
                                        prepend-icon="mdi-car"
                                        label="Auto"
                                        :items="cars"
                                        clearable
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
                                    <v-btn @click="closeForm">Abbrechen</v-btn>
                                    <v-btn @click="validateAndSubmitForm">{{saveButtonText}}</v-btn>
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
    const constants = require('../_services/constants')

    export default {
        data() {
            return {
                formIsValid: false,
                cars: constants.cars
            }
        },
        computed: {
            saveButtonText() {
                if(this.isUpdate){
                    return 'Fahrt aktualisieren'
                } else if (this.isAdmin){
                    return 'hinzufügen'
                } else {
                    return 'anfragen'
                }
            },
            rideDate() {
                return this.$store.state.ride.ride.date
            },
            germanDate() {
                if (this.rideDate == undefined) return undefined
                const p = this.rideDate.split(/\D/g)
                return [p[2], p[1], p[0]].join(".")
            },
            startTime: {
                get () {
                    return this.$store.state.ride.ride.startTime
                },
                set (value) {
                    this.$store.commit('ride/setStartTime', value)
                }
            },
            endTime: {
                get () {
                    return this.$store.state.ride.ride.endTime
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
                    this.$store.commit('ride/showAddEventForm', value)
                }
            },
            driver: {
                get () {
                    return this.$store.state.ride.ride.driver
                },
                set (value) {
                    this.$store.commit('ride/setDriver', value)
                }
            },
            destination: {
                get () {
                    return this.$store.state.ride.ride.destination
                },
                set (value) {
                    this.$store.commit('ride/setDestination', value)
                }
            },
            bigCarNeeded: {
                get () {
                    return this.$store.state.ride.ride.bigCarNeeded
                },
                set (value) {
                    this.$store.commit('ride/setBigCarNeeded', value)
                }
            },
            isUpdate: {
                get () {
                    return this.$store.state.ride.ride.isUpdate
                }
            },
            isAdmin: function () {
                return this.$store.state.account.status.loggedIn && this.$store.state.account.user.isAdmin
            },
            carName: {
                get () {
                    return this.$store.state.ride.ride.carName
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
            closeForm() {
              this.$refs.form.reset()
              this.showAddEventForm = false
            },
            allowedMinutes: m => m % 15 == 0,
        }
    }
</script>