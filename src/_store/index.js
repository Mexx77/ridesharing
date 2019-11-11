import Vue from 'vue';
import Vuex from 'vuex';

import { alert } from './alert.module';
import { account } from './account.module';
import { user } from './user.module';
import { ride } from './ride.module';
import { calendar } from './calendar.module';

Vue.use(Vuex);

export const store = new Vuex.Store({
    modules: {
        alert,
        account,
        user,
        ride,
        calendar
    }
});