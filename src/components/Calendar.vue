<template>
    <v-row class="fill-height">
        <v-col>
            <v-sheet height="64">
                <v-toolbar flat color="white">
                    <v-btn outlined class="mr-4" @click="setToday">
                        Heute
                    </v-btn>
                    <v-btn fab text small @click="prev">
                        <v-icon>mdi-chevron-left</v-icon>
                    </v-btn>
                    <v-btn fab text small @click="next" class="mr-4">
                        <v-icon>mdi-chevron-right</v-icon>
                    </v-btn>
                    <v-toolbar-title>{{ title }}</v-toolbar-title>
                    <v-spacer></v-spacer>
                    <v-menu bottom right>
                        <template v-slot:activator="{ on }">
                            <v-btn
                                    outlined
                                    v-on="on"
                            >
                                <span>{{ typeToLabel[type] }}</span>
                                <v-icon right>mdi-menu-down</v-icon>
                            </v-btn>
                        </template>
                        <v-list>
                            <v-list-item @click="type = 'day'">
                                <v-list-item-title>Tag</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="type = 'week'">
                                <v-list-item-title>Woche</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="type = 'month'">
                                <v-list-item-title>Monat</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="type = '4day'">
                                <v-list-item-title>4 Tage</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-toolbar>
            </v-sheet>
            <v-sheet>
                <v-calendar
                    ref="calendar"
                    v-model="focus"
                    color="primary"
                    :events="events"
                    :event-color="getEventColor"
                    :event-text-color="getEventTextColor"
                    :event-margin-bottom="3"
                    :event-overlap-threshold=60
                    :now="today"
                    :type="type"
                    :first-interval="7"
                    :interval-count="15"
                    @click:event="showEvent"
                    @click:more="viewDay"
                    @click:date="viewDay"
                    @click:time="addEvent"
                    @change="updateRange"
            ></v-calendar>
                <RideCard/>
                <AddEventForm/>
            </v-sheet>
        </v-col>
    </v-row>
</template>

<script>
    import AddEventForm from "./AddEventForm";
    import RideCard from "./RideCard";
    import * as constants from "../_services/constants"

    export default {
        components: {
            AddEventForm,
            RideCard
        },
        computed: {
            title() {
                const {start, end} = this
                if (!start || !end) {
                    return ''
                }

                const startMonth = this.monthFormatter(start)
                const endMonth = this.monthFormatter(end)
                const suffixMonth = startMonth === endMonth ? '' : endMonth

                const startYear = start.year
                const endYear = end.year
                const suffixYear = startYear === endYear ? '' : endYear

                const startDay = start.day + this.nth(start.day)
                const endDay = end.day + this.nth(end.day)

                switch (this.type) {
                    case 'month':
                        return `${startMonth} ${startYear}`
                    case 'week':
                    case '4day':
                        return `${startMonth} ${startDay} ${startYear} - ${suffixMonth} ${endDay} ${suffixYear}`
                    case 'day':
                        return `${startMonth} ${startDay} ${startYear}`
                }
                return ''
            },
            monthFormatter() {
                return this.$refs.calendar.getFormatter({
                    timeZone: 'UTC', month: 'long',
                })
            },
            focus: {
                get () {
                    return this.$store.state.rides.focus
                },
                set (value) {
                    this.$store.commit('rides/setFocus', value)
                }
            },
        },
        methods: {
            fourDaysFromNow() {
                return this.dateAddDays(this.today, 3)
            },
            dateAddDays(dateStr, nDays){
                let date =  new Date(dateStr);
                date.setDate(date.getDate() + nDays || 1);
                return [
                    date.getFullYear(),
                    this.zeroPad(date.getMonth()+1, 10),
                    this.zeroPad(date.getDate(), 10)
                ].join('-');
            },
            zeroPad(nr, base){
                const len = (String(base).length - String(nr).length) + 1;
                return len > 0? new Array(len).join('0') + nr : nr;
            },
            roundMinutes(hour, minute) {
                const m = (((minute + 7.5) / 15 | 0) * 15) % 60
                const h = ((((minute / 105) + .5) | 0) + hour) % 24
                const twoDigitM = m === 0 ? '00' : m
                return h + ':' + twoDigitM
            },
            addEvent(time) {
                this.focus = time.date
                this.$store.commit('rides/setStartTime', this.roundMinutes(time.hour, time.minute))
                this.time = time
                this.$store.commit('rides/setShowAddEventForm', true)
            },
            viewDay({date}) {
                this.focus = date
                this.type = 'day'
            },
            getEventColor(event) {
                return event.confirmed && event.carColor ? event.carColor : 'grey';
            },
            getEventTextColor(event) {
                if (event.carColor) {
                    return event.carColor === 'white' ? 'secondary' : 'white'
                } else {
                    return 'white';
                }
            },
            setToday() {
                this.focus = this.today
            },
            prev() {
                this.$refs.calendar.prev()
            },
            next() {
                this.$refs.calendar.next()
            },
            showEvent({nativeEvent, event}) {
                const open = () => {
                    this.$store.commit('rides/setSelectedEvent', event)
                    this.$store.commit('rides/setSelectedElement', nativeEvent.target)
                    setTimeout(() => this.$store.commit('rides/setSelectedOpen', true), 10)
                }

                if (this.$store.state.rides.selectedOpen) {
                    this.$store.commit('rides/setSelectedOpen', false)
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
                        this.events = response.data
                    });
            },
            nth(d) {
                return d > 3 && d < 21
                    ? 'th'
                    : ['th', 'st', 'nd', 'rd', 'th', 'th', 'th', 'th', 'th', 'th'][d % 10]
            },
            formatDate(date) {
                let month = '' + (date.getMonth() + 1),
                    day = '' + date.getDate(),
                    year = date.getFullYear();

                if (month.length < 2) month = '0' + month;
                if (day.length < 2) day = '0' + day;

                return [year, month, day].join('-');
            }
        },
        mounted: function () {
            this.$http
                .get(constants.hostname + '/rides?start=' + this.today  + '&end=' + this.fourDaysFromNow() )
                .then((response) => {
                    this.events = response.data
                });
        },
        data() {
            return {
                time: null,
                today: this.formatDate(new Date()),
                type: this.$vuetify.breakpoint.smAndDown ? 'day' : '4day',
                typeToLabel: {
                    month: 'Monat',
                    week: 'Woche',
                    day: 'Tag',
                    '4day': '4 Tage',
                },
                events: [],
            }
        }
    }
</script>