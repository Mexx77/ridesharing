import {authHeader} from './auth-header';
import * as constants from "./constants"

export const userService = {
  login,
  logout,
  register,
  delete: _delete,
  refreshToken
};

function login(usernamePhone, password) {
  const requestOptions = {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({usernamePhone, password})
  };

  return fetch(`${constants.hostname}/users/authenticate`, requestOptions)
    .then(handleResponse)
    .then(user => {
      // login successful if there's a jwt token in the response
      if (user.token) {
        // store user details and jwt token in local storage to keep user logged in between page refreshes
        localStorage.setItem('user', JSON.stringify(user));
      }
      return user;
    });
}

function refreshToken(token) {
  const requestOptions = {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: token
  };

  return fetch(`${constants.hostname}/users/refreshToken`, requestOptions)
    .then(handleResponse)
    .then(data => {
      if (!data.token) {
        return Promise.reject('no token in response');
      }
      const userString = localStorage.getItem('user');
      let user = JSON.parse(userString)
      user.token = data.token
      user.expires = data.expires
      localStorage.setItem('user', JSON.stringify(user));
    })
}

function logout() {
  // remove user from local storage to log user out
  localStorage.removeItem('user');
}

function register(user) {
  const requestOptions = {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(user)
  };

  return fetch(`${constants.hostname}/users/register`, requestOptions).then(handleResponse);
}

// prefixed function name with underscore because delete is a reserved word in javascript
function _delete(id) {
  const requestOptions = {
    method: 'DELETE',
    headers: authHeader()
  };

  return fetch(`${constants.hostname}/users/${id}`, requestOptions).then(handleResponse);
}

function handleResponse(response) {
  return response.text().then(text => {
    const data = text && JSON.parse(text);
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        //location.reload(true);
      }

      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }

    return data;
  });
}