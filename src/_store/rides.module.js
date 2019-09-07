const state = {
    startTime: '12:00',
    endTime: '',
    showAddEventForm: false,
    focus: ''
};

const actions = {};

const mutations = {
    setStartTime: (state, v) => state.startTime = v,
    setEndTime: (state, v) => state.endTime = v,
    setShowAddEventForm: (state, v) => state.showAddEventForm = v,
    setFocus: (state, v) => state.focus = v
};

export const rides = {
    namespaced: true,
    state,
    actions,
    mutations
};