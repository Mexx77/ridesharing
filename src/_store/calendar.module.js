const formatDate = (date) => {
    let month = '' + (date.getMonth() + 1),
      day = '' + date.getDate(),
      year = date.getFullYear();

    if (month.length < 2) month = '0' + month;
    if (day.length < 2) day = '0' + day;

    return [year, month, day].join('-');
};

const state = {
    type: 'day',
    focus: '',
    today: formatDate(new Date()),
};

const actions = {};

const mutations = {
    setType: (state, v) => state.type = v,
    setFocus: (state, v) => state.focus = v,
    setToday: (state, v) => state.today = v,
};

export const calendar = {
    namespaced: true,
    state,
    actions,
    mutations
};

