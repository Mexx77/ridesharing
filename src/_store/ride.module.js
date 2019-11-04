import {rideService} from "../_services/ride.service";
import {carProperties} from "../_services/constants";

const getDefaultRideState = () => {
    return {
        driver: '',
        destination: '',
        startTime: '12:00',
        endTime: '',
        bigCarNeeded: false,
        isUpdate: false,
        carName: '',
        id: '',
    }
}

const state = {
    showAddEventForm: false,
    focus: '',
    selectedOpen: false,
    selectedEvent: {},
    selectedElement: null,
    rides: [],
    ride: getDefaultRideState()
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
            id: state.ride.id,
            driver: state.ride.driver,
            destination: state.ride.destination,
            start: `${state.focus}T${state.ride.startTime}:00`,
            end: `${state.focus}T${state.ride.endTime}:00`,
            startTime: state.ride.startTime,
            endTime: state.ride.endTime,
            bigCarNeeded: state.ride.bigCarNeeded,
        }
        if (state.ride.carName !== "" && state.ride.carName !== undefined) {
            ride.carName = state.ride.carName
            ride.carColor = carProperties[ride.carName].color
        }
        rideService.update(ride).then(
            data => {
                const newRides = state.rides.filter(r => r.id !== state.ride.id).concat([data])
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
            driver: state.ride.driver,
            destination: state.ride.destination,
            start: `${state.focus}T${state.ride.startTime}:00`,
            end: `${state.focus}T${state.ride.endTime}:00`,
            startTime: state.ride.startTime,
            endTime: state.ride.endTime,
            bigCarNeeded: state.ride.bigCarNeeded,
        }
        if (state.ride.carName !== "" && state.ride.carName !== undefined) {
            ride.carName = state.ride.carName
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
    setDriver: (state, v) => state.ride.driver = v,
    setDestination: (state, v) => state.ride.destination = v,
    setBigCarNeeded: (state, v) => state.ride.bigCarNeeded = v,
    setStartTime: (state, v) => state.ride.startTime = v,
    setEndTime: (state, v) => state.ride.endTime = v,
    setCarName: (state, v) => state.ride.carName = v,
    showAddUpdateRideForm: (state, v) => {
        if (v) {
            // show form
            if (Object.keys(state.selectedEvent).length !== 0) {
                // update ride
                state.ride.isUpdate = true
                state.ride.driver = state.selectedEvent.driver
                state.ride.destination = state.selectedEvent.destination
                state.ride.startTime = state.selectedEvent.startTime
                state.ride.endTime = state.selectedEvent.endTime
                state.ride.bigCarNeeded = state.selectedEvent.bigCarNeeded
                state.ride.carName = state.selectedEvent.carName
                state.ride.id = state.selectedEvent.id
            } else {
                // new ride
                // state.ride = getDefaultRideState()
            }
        } else {
            // hide form
            state.ride = getDefaultRideState()
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
        state.ride = getDefaultRideState()
    },
};

export const ride = {
    namespaced: true,
    state,
    actions,
    mutations
};