<template>
    <v-row class="fill-height">
        <v-col>
            <v-sheet>
                <v-calendar
                    ref="calendar"
                    v-model="focus"
                    color="primary"
                    :events="$store.state.ride.rides"
                    :event-color="getEventColor"
                    :event-text-color="getEventTextColor"
                    :event-margin-bottom="3"
                    :event-overlap-threshold=60
                    :now="today"
                    :type="type"
                    :first-interval="7"
                    :interval-count="17"
                    :interval-format="(i) => i.time"
                    @click:event="showEvent"
                    @click:more="viewDay"
                    @click:date="viewDay"
                    @click:time="addEvent"
                    @change="updateRange"
                ></v-calendar>
                <RideCard/>
                <AddUpdateRideForm/>
            </v-sheet>
        </v-col>
    </v-row>
</template>

<script>
  import AddUpdateRideForm from "./AddUpdateRideForm";
  import RideCard from "./RideCard";
  import * as constants from "../_services/constants"
  import * as helper from '../_services/helper'
  import {mapState} from 'vuex'

  export default {
    components: {
      AddUpdateRideForm,
      RideCard
    },
    computed: {
      ...mapState('calendar', ['today']),
      focus: {
        get() {
          return this.$store.state.calendar.focus
        },
        set(value) {
          this.$store.commit('calendar/setFocus', value)
        }
      },
      type: {
        get() {
          return this.$store.state.calendar.type
        },
        set(value) {
          this.$store.commit('calendar/setType', value)
        }
      },
    },
    methods: {
      fourDaysFromNow() {
        return this.dateAddDays(this.today, 3)
      },
      dateAddDays(dateStr, nDays) {
        let date = new Date(dateStr);
        date.setDate(date.getDate() + nDays || 1);
        return [
          date.getFullYear(),
          this.zeroPad(date.getMonth() + 1, 10),
          this.zeroPad(date.getDate(), 10)
        ].join('-');
      },
      zeroPad(nr, base) {
        const len = (String(base).length - String(nr).length) + 1;
        return len > 0 ? new Array(len).join('0') + nr : nr;
      },
      roundMinutesAndPadZeros(hour, minute) {
        const m = (((minute + 7.5) / 15 | 0) * 15) % 60
        const h = ((((minute / 105) + .5) | 0) + hour) % 24
        return ('0' + h).slice(-2) + ':' + ('0' + m).slice(-2)
      },
      addEvent(time) {
        if (this.$store.state.account.status.loggedIn) {
          this.time = time
          this.$store.dispatch('ride/showAddUpdateRideForm', {
            visible: true,
            isUpdate: false,
            startTime: this.roundMinutesAndPadZeros(time.hour, time.minute),
            date: time.date
          })
        } else {
          this.$store.dispatch('alert/error', {
            message: '💡 Bitte melde dich an, um Fahrten hinzuzufügen',
            timeout: 6000
          })
        }
      },
      viewDay({date}) {
        this.focus = date
        this.type = 'day'
      },
      getEventColor: helper.getEventColor,
      getEventTextColor: helper.getEventTextColor,
      prev() {
        this.$refs.calendar.prev()
      },
      next() {
        this.$refs.calendar.next()
      },
      showEvent({nativeEvent, event}) {
        const open = () => {
          this.$store.commit('ride/setSelectedEvent', event)
          this.$store.commit('ride/setSelectedElement', nativeEvent.target)
          setTimeout(() => this.$store.commit('ride/setSelectedOpen', true), 10)
        }

        if (this.$store.state.ride.selectedOpen) {
          this.$store.commit('ride/setSelectedOpen', false)
          setTimeout(open, 10)
        } else {
          open()
        }

        nativeEvent.stopPropagation()
      },
      updateRange({start, end}) {
        this.$http
          .get(constants.hostname + '/rides?start=' + start.date + '&end=' + end.date)
          .then((response) => {
            this.$store.commit('ride/setRides', response.data)
          });
      }
    },
    mounted: function () {
      this.focus = this.today;
      if (this.$vuetify.breakpoint.mdAndUp) {
        this.type = '4day'
      }
      this.$http
        .get(constants.hostname + '/rides?start=' + this.today + '&end=' + this.fourDaysFromNow())
        .then((response) => {
          this.$store.commit('ride/setRides', response.data)
        });
    },
    data() {
      return {
        time: null,
      }
    }
  }
</script>

<style scoped>
    >>> .v-calendar-daily__day {
        cursor: pointer;
    }

    >>> .v-calendar .v-event-timed {
        font-size: 1em;
        padding: 5px;
    }

    >>> .v-calendar-daily__scroll-area {
        overflow-y: auto;
    }
</style>