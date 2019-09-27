import {rideService} from "../_services/ride.service";
import {carProperties} from "../_services/constants";

const state = {
    driver: '',
    destination: '',
    startTime: '12:00',
    endTime: '',
    bigCarNeeded: false,
    showAddEventForm: false,
    focus: '',
    selectedOpen: false,
    selectedEvent: {},
    selectedElement: null,
    rides: [],
    isUpdate: false,
    carName: '',
    id: '',
};

const actions = {
    delete({commit}, id) {
        rideService.delete(id)
            .then(
                () => commit('deleteSuccess', id),
            );
    },
    updateRide({commit, dispatch}) {
        const ride = {
            id: state.id,
            driver: state.driver,
            destination: state.destination,
            start: `${state.focus}T${state.startTime}:00`,
            end: `${state.focus}T${state.endTime}:00`,
            startTime: state.startTime,
            endTime: state.endTime,
            bigCarNeeded: state.bigCarNeeded,
        }
        if(state.carName !== "" && state.carName !== undefined){
            ride.carName = state.carName
            ride.carColor = carProperties[ride.carName].color
        }
        rideService.update(ride).then(
          data => {
              const newRides = state.rides.filter(r => r.id !== state.id).concat([data])
              commit('setRides', newRides)
              commit('showAddUpdateRideForm', false)
              dispatch('alert/success', {
                  message: 'Fahrt erfolgreich aktualisiert',
                  visible: true
              }, {root: true});
          },
          () => {
              dispatch('alert/error', {
                  message: 'Ups, da ist was fehlgeschlagen - sorry',
                  visible: true
              }, {root: true});
          }
        )
    },
    addRide({commit, dispatch, rootState}) {
        const ride = {
            driver: state.driver,
            destination: state.destination,
            start: `${state.focus}T${state.startTime}:00`,
            end: `${state.focus}T${state.endTime}:00`,
            startTime: state.startTime,
            endTime: state.endTime,
            bigCarNeeded: state.bigCarNeeded,
        }
        if(state.carName !== "" && state.carName !== undefined){
            ride.carName = state.carName
            ride.carColor = carProperties[ride.carName].color
        }
        rideService.add(ride).then(
            data => {
                let msg = 'Danke, deine Reservierungsanfrage wurde entgegengenommen'
                if (rootState.account.status.loggedIn && rootState.account.user.isAdmin) {
                    msg = 'Fahrt gespeichert'
                }
                commit('showAddUpdateRideForm', false)
                const newRides = state.rides.concat([data])
                commit('setRides', newRides)
                dispatch('alert/success', {
                    message: msg,
                    visible: true
                }, {root: true});
            },
            () => {
                dispatch('alert/error', {
                    message: 'Ups, da ist was fehlgeschlagen - sorry',
                    visible: true
                }, {root: true});
            }
        )
    }
}


const mutations = {
    setDriver: (state, v) => state.driver = v,
    setDestination: (state, v) => state.destination = v,
    setBigCarNeeded: (state, v) => state.bigCarNeeded = v,
    setStartTime: (state, v) => state.startTime = v,
    setEndTime: (state, v) => state.endTime = v,
    setCarName: (state, v) => state.carName = v,
    showAddUpdateRideForm: (state, v) => {
        if (v) {
            if (Object.keys(state.selectedEvent).length !== 0) {
                state.isUpdate = true
                state.driver = state.selectedEvent.driver
                state.destination = state.selectedEvent.destination
                state.startTime = state.selectedEvent.startTime
                state.endTime = state.selectedEvent.endTime
                state.bigCarNeeded = state.selectedEvent.bigCarNeeded
                state.carName = state.selectedEvent.carName
                state.id = state.selectedEvent.id
            }
        } else {
            state.isUpdate = false
            state.driver = ""
            state.destination = ""
            state.startTime = ""
            state.endTime = ""
            state.bigCarNeeded = false
            state.carName = ""
            state.carColor = ""
            state.id = ""
        }
        state.showAddEventForm = v
    },
    setFocus: (state, v) => state.focus = v,
    setSelectedOpen: (state, v) => {
        if (!v) {
            state.selectedEvent = {}
        }
        state.selectedOpen = v
    },
    setSelectedEvent: (state, v) => state.selectedEvent = v,
    setSelectedElement: (state, v) => state.selectedElement = v,
    setRides: (state, v) => state.rides = v,
    deleteSuccess(state, id) {
        state.selectedOpen = false
        state.rides = state.rides.filter(ride => ride.id !== id)
    },
};

export const ride = {
    namespaced: true,
    state,
    actions,
    mutations
};