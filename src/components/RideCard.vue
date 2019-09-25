<template>
    <v-menu
        v-model="selectedOpen"
        :close-on-content-click="false"
        :activator="this.$store.state.ride.selectedElement"
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
                <template v-if="isAdmin">
                    <v-btn icon small @click="editRide">
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
    import {mapActions} from 'vuex'

    export default {
        computed: {
            selectedOpen: {
                get () {
                  return this.$store.state.ride.selectedOpen
                },
                set (value) {
                  this.$store.commit('ride/setSelectedOpen', value)
                }
            },
            selectedEvent: {
                get () {
                  return this.$store.state.ride.selectedEvent
                }
            },
            isAdmin: function () {
                return this.$store.state.account.status.loggedIn && this.$store.state.account.user.isAdmin
            }
        },
        methods: {
            ...mapActions('ride', ['delete']),
            deleteRide() {
                const confirmed = confirm(`Die Fahrt ${this.selectedEvent.name} wirklich löschen?`)
                if (!confirmed) {
                    return
                }
                this.delete(this.selectedEvent.id)
            },
            editRide() {
                this.$store.commit('ride/setShowAddEventForm', true)
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