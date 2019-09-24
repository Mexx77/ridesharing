<template>
    <v-menu
            v-model="selectedOpen"
            :close-on-content-click="false"
            :activator="this.$store.state.rides.selectedElement"
            offset-x
    >
        <v-card
                color="grey lighten-4"
                min-width="300px"
                flat
        >
            <v-toolbar
                    :color="selectedEvent.getEventColor"
                    :style="{color: selectedEvent.getEventTextColor}"
            >

                <v-toolbar-title v-html="selectedEvent.name"></v-toolbar-title>
                <v-spacer></v-spacer>
                <template v-if="$store.state.account.status.loggedIn && $store.state.account.user.isAdmin">
                    <v-btn icon small>
                        <v-icon>mdi-pencil</v-icon>
                    </v-btn>
                    <v-btn icon small @click="deleteRide">
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </template>
            </v-toolbar>
            <v-card-text>
                <span v-html="selectedEvent.details"></span>
            </v-card-text>
            <v-card-actions>
                <v-btn
                        text
                        color="secondary"
                        @click="selectedOpen = false"
                >
                    Abbrechen
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-menu>
</template>

<script>
    let rideService = require("../_services/ride.service");

    export default {
        computed: {
            selectedOpen: {
                get () {
                  return this.$store.state.rides.selectedOpen
                },
                set (value) {
                  this.$store.commit('rides/setSelectedOpen', value)
                }
            },
            selectedEvent: {
                get () {
                  return this.$store.state.rides.selectedEvent
                }
            },

        },
        methods: {
            deleteRide() {
                rideService.rideService.delete(this.selectedEvent.id, this.selectedEvent.name)
                this.selectedOpen = false
                // TODO: delete event from events in store (needs refactoring)
            },
            getEventColor(event) {
                return event.confirmed && event.carColor ? event.carColor : 'grey';
            },
            getEventTextColor(event) {
                if (event.carColor) {
                    return event.carColor === 'white' ? 'secondary' : 'white'
                } else {
                    return 'white';
                }
            }
        },
        data() {
            return {

            }
        }
    }
</script>

<style>
    .v-event-timed {
        font-size: 16px !important;
    }
</style>