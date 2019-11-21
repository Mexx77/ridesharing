import {userService} from '../_services/user.service';

const user = JSON.parse(localStorage.getItem('user'));
const state = {
  status: user ? {loggedIn: true} : {},
  user: user ? user : null,
  showLoginForm: false,
  showRegisterForm: false,
};

const actions = {
  login({dispatch, commit}, {usernamePhone, password, message}) {
    return new Promise((resolve, reject) => {

      commit('loginRequest', {usernamePhone});

      userService.login(usernamePhone, password)
        .then(
          user => {
            commit('loginSuccess', user);
            dispatch('account/showLoginForm', false, {root: true})
            dispatch('alert/success', {
              message: message ? message : `Erfolgreich angemeldet. Schön dich wiederzusehen, ${user.firstName}! :)`,
              timeout: 10000
            }, {root: true});
            resolve('loginSuccess');
          },
          error => {
            commit('loginFailure', error);
            dispatch('alert/error', {
              message: error,
              timeout: 6000
            }, {root: true});
            reject('loginFailure', error);
          }
        );
    })
  },
  logout({commit, dispatch}) {
    const firstName = state.user.firstName
    userService.logout();
    commit('logout');
    dispatch('alert/success', {
      message: `Erfolgreich abgemeldet. Bis bald, ${firstName}! :)`,
      timeout: 6000
    }, {root: true});
  },
  refreshToken({commit}) {
    userService.refreshToken(state.user.token)
      .catch(() => commit('logout'))
  },
  register({dispatch, commit}, user) {
    return new Promise((resolve, reject) => {
      commit('registerRequest', user);

      userService.register(user)
        .then(
          () => {
            commit('registerSuccess');
            dispatch('account/showRegisterForm', false, {root: true});
            dispatch('login', {
              usernamePhone: user.phone,
              password: user.password,
              message: `Registrierung erfolgreich. Wir haben dich auch gleich eingeloggt, ${user.firstName}! :)`
            });
            resolve('registerSuccess');
          },
          error => {
            commit('registerFailure', error);
            dispatch('alert/error', {message: error, timeout: 6000}, {root: true});
            reject('registerFailure', error)
          }
        );
    })
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