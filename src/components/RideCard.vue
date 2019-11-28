<template>
    <v-menu
        v-model="selectedOpen"
        :close-on-content-click="false"
        :close-on-click="true"
        :activator="this.$store.state.ride.selectedElement"
        offset-x
        max-width="400"
        @keydown.esc="selectedOpen = false"
    >
        <v-card color="grey lighten-4" flat>
            <v-toolbar
                :color="selectedEvent.getEventColor"
                :style="{color: selectedEvent.getEventTextColor}"
            >
                <v-toolbar-title>
                    <span>{{selectedEvent.name}}</span>
                    <v-progress-circular color="primary" v-if="status.updating" class="ml-2 mb-1"
                                         size="17" width="2" indeterminate></v-progress-circular>
                </v-toolbar-title>
                <v-spacer/>
                <v-btn @click="selectedOpen = false" icon small>
                    <v-icon>mdi-close</v-icon>
                </v-btn>
            </v-toolbar>
            <v-card-text class="pb-0">
                <span v-html="selectedEvent.details"></span>
            </v-card-text>
            <div v-if="isRideConfirmed && isMyRide && !isAdmin" class="pl-4 pr-4 pt-2 caption">
                <span class="font-weight-bold red--text">Achtung:</span> Das Ändern einer bestätigten Fahrt macht diese
                unbestätigt. Ein Admin muss dir dann erneut ein Auto zuweisen.
            </div>
            <v-card-actions>
                <v-container pt-0 pb-0>
                    <v-row dense v-if="isAdmin">
                        <v-col cols="9">
                            <v-select
                                :disabled="status.updating"
                                v-model="carName"
                                label="Auto"
                                :items="cars"
                                clearable
                                v-on:change="carChangeHandler"
                            ></v-select>
                        </v-col>
                    </v-row>
                    <v-row dense>
                        <v-col>
                            <v-btn v-if="isAdmin || isMyRide" text color="primary" @click="editRide" class="pl-0 pa-1">
                                <v-icon>mdi-pencil</v-icon>
                                Ändern
                            </v-btn>
                            <v-btn v-if="isAdmin || isMyRide" text color="red" @click="deleteRide" class="pa-1">
                                <v-icon>mdi-delete</v-icon>
                                Löschen
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-container>
            </v-card-actions>
        </v-card>
    </v-menu>
</template>

<script>
  import * as helper from '../_services/helper'
  import {mapActions, mapState} from 'vuex'

  const constants = require('../_services/constants')

  export default {
    computed: {
      ...mapState('ride', ['status']),
      selectedOpen: {
        get() {
          return this.$store.state.ride.selectedOpen
        },
        set(value) {
          this.$store.commit('ride/setSelectedOpen', value)
        }
      },
      selectedEvent: {
        get() {
          return this.$store.state.ride.selectedEvent
        }
      },
      isAdmin: function () {
        return this.$store.state.account.status.loggedIn && this.$store.state.account.user.isAdmin
      },
      isMyRide: function () {
        return this.$store.state.account.status.loggedIn &&
          this.$store.state.account.user.id === this.selectedEvent.userId
      },
      isRideConfirmed: function () {
        return this.selectedEvent.hasOwnProperty('carName')
      },
      carName: {
        get() {
          return this.selectedEvent.carName
        },
        set(value) {
          this.$store.commit('ride/setCarName', value)
        }
      },
    },
    methods: {
      ...mapActions('ride', ['delete', 'updateRide']),
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
      carChangeHandler() {
        this.updateRide()
      },
      getEventColor: helper.getEventColor,
      getEventTextColor: helper.getEventTextColor
    },
    data() {
      return {
        cars: constants.cars
      }
    }
  }
</script>
<style scoped>
    >>> td:first-child {
        padding-right: 10px;
    }
</style>