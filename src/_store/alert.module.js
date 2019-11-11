const state = {
    type: null,
    message: null,
    visible: false,
    timeout: 0,
};

const actions = {
    success({ commit }, {message, timeout}) {
        commit('success', {message, timeout});
    },
    error({ commit }, {message, timeout}) {
        commit('error', {message, timeout});
    },
    info({ commit }, {message, timeout}) {
        commit('info', {message, timeout});
    },
    clear({ commit }) {
        commit('clear');
    },
    setVisibility({ commit }, value) {
        commit('setVisibility', value);
    },
    setTimeout({ commit }, value) {
        commit('setTimeout', value);
    }
};

const mutations = {
    success(state, {message, timeout}) {
        state.type = 'success';
        state.message = message;
        state.timeout = timeout;
        state.visible = true;
    },
    error(state, {message, timeout}) {
        state.type = 'error';
        state.message = message;
        state.timeout = timeout;
        state.visible = true;
    },
    info(state, {message, timeout}) {
        state.type = 'info';
        state.message = message;
        state.timeout = timeout;
        state.visible = true;
    },
    clear(state) {
        state.type = null;
        state.message = null;
        state.visible = false;
        state.timeout = 0;
    },
    setVisibility(state, value) {
        state.visible = value;
    },
    setTimeout(state, value) {
        state.timeout = value;
    }
};

export const alert = {
    namespaced: true,
    state,
    actions,
    mutations
};