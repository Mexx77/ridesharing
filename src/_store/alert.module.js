const state = {
    type: null,
    message: null,
    visible: false
};

const actions = {
    success({ commit }, {message, visible}) {
        commit('success', {message, visible});
    },
    error({ commit }, {message, visible}) {
        commit('error', {message, visible});
    },
    info({ commit }, {message, visible}) {
        commit('info', {message, visible});
    },
    clear({ commit }) {
        commit('clear');
    },
    setVisibility({ commit }, value) {
        commit('setVisibility', value);
    }
};

const mutations = {
    success(state, {message, visible}) {
        state.type = 'success';
        state.message = message;
        state.visible = visible;
    },
    error(state, {message, visible}) {
        state.type = 'error';
        state.message = message;
        state.visible = visible;
    },
    info(state, {message, visible}) {
        state.type = 'info';
        state.message = message;
        state.visible = visible;
    },
    clear(state) {
        state.type = null;
        state.message = null;
        state.visible = false;
    },
    setVisibility(state, value) {
        state.visible = value;
    }
};

export const alert = {
    namespaced: true,
    state,
    actions,
    mutations
};