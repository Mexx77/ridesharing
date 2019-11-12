<template>
    <v-menu
        v-model="selectedOpen"
        :close-on-content-click="false"
        :close-on-click="true"
        :activator="this.$store.state.ride.selectedElement"
        offset-x
        min-width="270"
        max-width="400"
        @keydown.esc="selectedOpen = false"
    >
        <v-card color="grey lighten-4" flat>
            <v-toolbar
                :color="selectedEvent.getEventColor"
                :style="{color: selectedEvent.getEventTextColor}"
            >
                <v-toolbar-title v-html="selectedEvent.name"></v-toolbar-title>
                <v-spacer/>
                <v-btn @click="selectedOpen = false" icon small><v-icon>mdi-close</v-icon></v-btn>
            </v-toolbar>
            <v-card-text class="pb-0">
                <span v-html="selectedEvent.details"></span>
            </v-card-text>
            <v-card-actions>
                <v-btn v-if="isAdmin" text color="primary" @click="editRide">
                    <v-icon>mdi-pencil</v-icon> Ändern
                </v-btn>
                <v-btn v-if="isAdmin" text color="red" @click="deleteRide">
                    <v-icon>mdi-delete</v-icon> Löschen
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-menu>
</template>

<script>
    import * as helper from '../_services/helper'
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
                this.$store.dispatch('ride/showAddUpdateRideForm', {visible: true, isUpdate: true})
            },
            getEventColor: helper.getEventColor,
            getEventTextColor: helper.getEventTextColor
        },
        data() {
            return {

            }
        }
    }
</script>
<style scoped>
    >>> td:first-child {
        padding-right: 10px;
    }
</style>