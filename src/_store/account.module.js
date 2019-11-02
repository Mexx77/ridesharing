import { userService } from '../_services/user.service';

const user = JSON.parse(localStorage.getItem('user'));
const state = user
    ? { status: { loggedIn: true }, user }
    : { status: {}, user: null };

const actions = {
    login({ dispatch, commit }, { username, password }) {
        commit('loginRequest', { username });

        userService.login(username, password)
            .then(
                user => {
                    commit('loginSuccess', user);
                    dispatch('user/showLoginForm', false, {root: true})
                    dispatch('alert/success', {
                        message: `Erfolgreich angemeldet. Schön dich wiederzusehen, ${user.firstName}! :)`,
                        visible: true
                    }, {root: true});
                },
                error => {
                    commit('loginFailure', error);
                    dispatch('alert/error', {
                        message: error,
                        visible: true
                    }, {root: true});
                }
            );
    },
    logout({ commit }) {
        userService.logout();
        commit('logout');
    },
    refreshToken({ commit }) {
        userService.refreshToken(state.user.token)
            .catch(() => commit('logout'))
    },
    register({ dispatch, commit }, user) {
        commit('registerRequest', user);
        const firstName = user.firstName;

        userService.register(user)
            .then(
              () => {
                    commit('registerSuccess');
                    dispatch('alert/success',
                        { message: `Registrierung erfolgreich! Bitte melde dich an, ${firstName} :)`, visible: true },
                        { root: true }
                    );
                    dispatch('user/showRegisterForm', false, { root: true });
                    dispatch('user/showLoginForm', true, { root: true });
                },
                error => {
                    commit('registerFailure', error);
                    dispatch('alert/error', {message: error, visible: true}, { root: true });
                }
            );
    }
};

const mutations = {
    loginRequest(state, user) {
        state.status = { loggingIn: true };
        state.user = user;
    },
    loginSuccess(state, user) {
        state.status = { loggedIn: true };
        state.user = user;
    },
    loginFailure(state) {
        state.status = {};
        state.user = null;
    },
    logout(state) {
        state.status = {};
        state.user = null;
    },
    registerRequest(state) {
        state.status = { registering: true };
    },
    registerSuccess(state) {
        state.status = {};
    },
    registerFailure(state) {
        state.status = {};
    }
};

export const account = {
    namespaced: true,
    state,
    actions,
    mutations
};