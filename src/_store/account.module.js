import {userService} from '../_services/user.service';

const user = JSON.parse(localStorage.getItem('user'));
const state = {
  status: user ? {loggedIn: true} : {},
  user: user ? user : null,
  showLoginForm: false,
  showRegisterForm: false,
};

const actions = {
  login({dispatch, commit}, {usernamePhone, password}) {
    commit('loginRequest', {usernamePhone});

    userService.login(usernamePhone, password)
      .then(
        user => {
          commit('loginSuccess', user);
          dispatch('account/showLoginForm', false, {root: true})
          dispatch('alert/success', {
            message: `Erfolgreich angemeldet. Schön dich wiederzusehen, ${user.firstName}! :)`,
            timeout: 6000
          }, {root: true});
        },
        error => {
          commit('loginFailure', error);
          dispatch('alert/error', {
            message: error,
            timeout: 6000
          }, {root: true});
        }
      );
  },
  logout({commit}) {
    userService.logout();
    commit('logout');
  },
  refreshToken({commit}) {
    userService.refreshToken(state.user.token)
      .catch(() => commit('logout'))
  },
  register({dispatch, commit}, user) {
    commit('registerRequest', user);
    const firstName = user.firstName;

    userService.register(user)
      .then(
        () => {
          commit('registerSuccess');
          dispatch('alert/success',
            {message: `Registrierung erfolgreich! Bitte melde dich an, ${firstName} :)`, timeout: 6000},
            {root: true}
          );
          dispatch('account/showRegisterForm', false, {root: true});
          dispatch('account/showLoginForm', true, {root: true});
        },
        error => {
          commit('registerFailure', error);
          dispatch('alert/error', {message: error, timeout: 6000}, {root: true});
        }
      );
  },
  showLoginForm({commit}, v) {
    commit('showLoginForm', v);
  },
  showRegisterForm({commit}, v) {
    commit('showRegisterForm', v);
  },
};

const mutations = {
  loginRequest(state, user) {
    state.status = {loggingIn: true};
    state.user = user;
  },
  loginSuccess(state, user) {
    state.status = {loggedIn: true};
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
    state.status = {registering: true};
  },
  registerSuccess(state) {
    state.status = {};
  },
  registerFailure(state) {
    state.status = {};
  },
  showLoginForm: (state, v) => state.showLoginForm = v,
  showRegisterForm: (state, v) => state.showRegisterForm = v,
};

export const account = {
  namespaced: true,
  state,
  actions,
  mutations
};