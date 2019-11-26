import {rideService} from "../_services/ride.service";
import {carProperties} from "../_services/constants";

const defaultStartTime = '12:00'
const getDefaultRideState = () => {
  return {
    driver: '',
    destination: '',
    startTime: defaultStartTime,
    endTime: '',
    bigCarNeeded: false,
    isUpdate: false,
    carName: '',
    id: '',
    date: ''
  }
}

const state = {
  showAddEventForm: false,
  selectedOpen: false,
  selectedEvent: {},
  selectedElement: null,
  rides: [],
  ride: getDefaultRideState(),
  unconfirmedRides: 0
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
      start: `${state.ride.date}T${state.ride.startTime}:00`,
      end: `${state.ride.date}T${state.ride.endTime}:00`,
      startTime: state.ride.startTime,
      endTime: state.ride.endTime,
      bigCarNeeded: state.ride.bigCarNeeded,
      date: state.ride.date,
      userId: state.ride.userId
    }
    if (state.ride.carName !== "" && state.ride.carName !== undefined) {
      ride.carName = state.ride.carName
      ride.carColor = carProperties[ride.carName].color
    }
    rideService.update(ride).then(
      data => {
        const newRides = state.rides.filter(r => r.id !== state.ride.id).concat([data])
        commit('setRides', newRides)
        commit('showAddEventForm', false)
        dispatch('refreshUnconfirmedRides')
        dispatch('alert/success', {
          message: 'Fahrt erfolgreich aktualisiert',
          timeout: 6000
        }, {root: true});
      },
      () => {
        dispatch('alert/error', {
          message: 'Ups, da ist was fehlgeschlagen - sorry',
          timeout: 6000
        }, {root: true});
      }
    )
  },
  addRide({commit, dispatch, rootState}) {
    const ride = {
      driver: state.ride.driver,
      destination: state.ride.destination,
      start: `${state.ride.date}T${state.ride.startTime}:00`,
      end: `${state.ride.date}T${state.ride.endTime}:00`,
      startTime: state.ride.startTime,
      endTime: state.ride.endTime,
      bigCarNeeded: state.ride.bigCarNeeded,
      date: state.ride.date
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
        commit('showAddEventForm', false)
        const newRides = state.rides.concat([data])
        commit('setRides', newRides)
        dispatch('alert/success', {
          message: msg,
          timeout: 6000
        }, {root: true});
      },
      () => {
        dispatch('alert/error', {
          message: 'Ups, da ist was fehlgeschlagen - sorry',
          timeout: 6000
        }, {root: true});
      }
    )
  },
  showAddUpdateRideForm: ({commit, rootState}, {visible, isUpdate, startTime, date}) => {
    if (visible) {
      if (isUpdate) {
        commit('setRide', {
          driver: state.selectedEvent.driver,
          destination: state.selectedEvent.destination,
          startTime: state.selectedEvent.startTime,
          endTime: state.selectedEvent.endTime,
          bigCarNeeded: state.selectedEvent.bigCarNeeded,
          isUpdate: true,
          carName: state.selectedEvent.carName,
          id: state.selectedEvent.id,
          date: state.selectedEvent.date,
          userId: state.selectedEvent.userId
        })
      } else {
        let ride = getDefaultRideState()
        ride.date = date
        ride.startTime = startTime === '' ? defaultStartTime : startTime
        if (rootState.account.status.loggedIn) {
          ride.driver = rootState.account.user.firstName + ' ';
          ride.driver += rootState.account.user.lastName.substr(0, 1) + '.'
        }
        commit('setRide', ride)
      }
    }
    commit('showAddEventForm', visible)
  },
  refreshUnconfirmedRides({commit}) {
    rideService.unconfirmedRides().then(
      count => {
        commit('setUnconfirmedRides', count)
      }
    )
  }
}


const mutations = {
  setUnconfirmedRides: (state, v) => state.unconfirmedRides = v,
  reduceUnconfirmedRides: (state) => --state.unconfirmedRides,
  showAddEventForm: (state, v) => state.showAddEventForm = v,
  setDriver: (state, v) => state.ride.driver = v,
  setDestination: (state, v) => state.ride.destination = v,
  setBigCarNeeded: (state, v) => state.ride.bigCarNeeded = v,
  setStartTime: (state, v) => state.ride.startTime = v,
  setEndTime: (state, v) => state.ride.endTime = v,
  setCarName: (state, v) => state.ride.carName = v,
  setSelectedOpen: (state, v) => {
    if (!v) {
      state.selectedEvent = {}
    }
    state.selectedOpen = v
  },
  setSelectedEvent: (state, v) => state.selectedEvent = v,
  setSelectedElement: (state, v) => state.selectedElement = v,
  setRides: (state, v) => state.rides = v,
  setRide: (state, ride) => state.ride = ride,
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