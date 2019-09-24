const state = {
    startTime: '12:00',
    endTime: '',
    showAddEventForm: false,
    focus: '',
    selectedOpen: false,
    selectedEvent: {},
    selectedElement: null,
};

const actions = {};

const mutations = {
    setStartTime: (state, v) => state.startTime = v,
    setEndTime: (state, v) => state.endTime = v,
    setShowAddEventForm: (state, v) => state.showAddEventForm = v,
    setFocus: (state, v) => state.focus = v,
    setSelectedOpen: (state, v) => state.selectedOpen = v,
    setSelectedEvent: (state, v) => state.selectedEvent = v,
    setSelectedElement: (state, v) => state.selectedElement = v
};

export const rides = {
    namespaced: true,
    state,
    actions,
    mutations
};