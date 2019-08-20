<template>
    <v-row class="fill-height">
        <v-col>
            <v-sheet height="64">
                <v-toolbar flat color="white">
                    <v-btn outlined class="mr-4" @click="setToday">
                        Today
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
                                <v-list-item-title>Day</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="type = 'week'">
                                <v-list-item-title>Week</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="type = 'month'">
                                <v-list-item-title>Month</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="type = '4day'">
                                <v-list-item-title>4 days</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-toolbar>
            </v-sheet>
            <v-sheet height="600">
                <v-calendar
                        ref="calendar"
                        v-model="focus"
                        color="primary"
                        :events="events"
                        :event-color="getEventColor"
                        :event-margin-bottom="3"
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
                <v-menu
                        v-model="selectedOpen"
                        :close-on-content-click="false"
                        :activator="selectedElement"
                        full-width
                        offset-x
                >
                    <v-card
                            color="grey lighten-4"
                            min-width="350px"
                            flat
                    >
                        <v-toolbar
                                :color="selectedEvent.color"
                                dark
                        >
                            <v-btn icon small>
                                <v-icon>mdi-pencil</v-icon>
                            </v-btn>
                            <v-toolbar-title v-html="selectedEvent.name"></v-toolbar-title>
                            <v-spacer></v-spacer>
                            <v-btn icon small>
                                <v-icon>mdi-heart-outline</v-icon>
                            </v-btn>
                            <v-btn icon small>
                                <v-icon>mdi-dots-vertical</v-icon>
                            </v-btn>
                        </v-toolbar>
                        <v-card-text>
                            <span v-html="selectedEvent.details"></span>
                        </v-card-text>
                        <v-card-actions>
                            <v-btn
                                    text
                                    color="secondary"
                                    @click="selectedOpen = false"
                            >
                                Cancel
                            </v-btn>
                        </v-card-actions>
                    </v-card>
                </v-menu>
                <v-row justify="center">
                    <v-dialog
                            v-model="showAddEventForm"
                            :fullscreen="$vuetify.breakpoint.smAndDown ? true : false"
                            hide-overlay max-width="600px"
                    >
                        <v-card>
                            <v-toolbar color="indigo" dark>
                                <v-toolbar-title>
                                    <v-icon>mdi-car</v-icon>
                                    Reservierung am {{focus}}
                                </v-toolbar-title>
                            </v-toolbar>
                            <v-card-text>
                                <v-container pa-0>
                                    <v-row>
                                        <v-col cols="6">
                                            <v-select
                                                    prepend-icon="mdi-car"
                                                    :items="['Red Bus', 'White Bus', 'e-Auto', 'Little Red']"
                                                    label="Auto*"
                                                    required
                                            ></v-select>
                                        </v-col>
                                        <v-col cols="6">
                                            <v-text-field
                                                    prepend-icon="mdi-account"
                                                    label="Name des Fahrers*"
                                                    required
                                            ></v-text-field>
                                        </v-col>
                                        <v-col cols="6">
                                            <v-combobox
                                                    prepend-icon="mdi-city"
                                                    :items="['Lüneburg', 'Dannenberg', 'Hitzacker']"
                                                    label="Fahrtziel*"
                                                    required
                                            ></v-combobox>
                                        </v-col>
                                    </v-row>
                                    <v-row>
                                        <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                            <v-time-picker
                                                    v-model="startTime"
                                                    color="indigo"
                                                    :width="272"
                                                    format="24hr"
                                            ></v-time-picker>
                                        </v-col>
                                        <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                            <v-time-picker
                                                    v-model="endTime"
                                                    color="indigo"
                                                    :width="272"
                                                    format="24hr"
                                                    :min="startTime"
                                            ></v-time-picker>
                                        </v-col>
                                    </v-row>
                                </v-container>
                            </v-card-text>
                            <v-card-actions class="mr-2 pb-4 pt-0">
                                <v-spacer></v-spacer>
                                <v-btn @click="showAddEventForm = false">Abbrechen</v-btn>
                                <v-btn @click="showAddEventForm = false">Anfragen</v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </v-row>
            </v-sheet>
        </v-col>
    </v-row>
</template>

<script>
    export default {
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
        },
        methods: {
            roundMinutes(hour, minute) {
                const m = (((minute + 7.5) / 15 | 0) * 15) % 60
                const h = ((((minute / 105) + .5) | 0) + hour) % 24
                return h + ':' + m
            },
            addEvent(time) {
                // eslint-disable-next-line
                console.log(time)
                this.focus = time.date
                this.startTime = this.roundMinutes(time.hour, time.minute)
                this.time = time
                this.showAddEventForm = true
            },
            viewDay({date}) {
                this.focus = date
                this.type = 'day'
            },
            getEventColor(event) {
                return event.color
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
                    this.selectedEvent = event
                    this.selectedElement = nativeEvent.target
                    setTimeout(() => this.selectedOpen = true, 10)
                }

                if (this.selectedOpen) {
                    this.selectedOpen = false
                    setTimeout(open, 10)
                } else {
                    open()
                }

                nativeEvent.stopPropagation()
            },
            updateRange({start, end}) {
                // You could load events from an outside source (like database) now that we have the start and end dates on the calendar
                this.start = start
                this.end = end
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
        data() {
            return {
                showAddEventForm: false,
                startTime: '12:00',
                endTime: null,
                time: null,
                today: this.formatDate(new Date()),
                focus: this.today,
                type: '4day',
                typeToLabel: {
                    month: 'Month',
                    week: 'Week',
                    day: 'Day',
                    '4day': '4 Days',
                },
                selectedEvent: {},
                selectedElement: null,
                selectedOpen: false,
                events: [
                    {
                        name: 'Vacation',
                        details: 'Going to the beach!',
                        start: '2019-08-17',
                        end: '2019-08-18',
                        color: 'blue',
                    },
                    {
                        name: 'Meeting',
                        details: 'Spending time on how we do not have enough time',
                        start: '2019-08-17 09:00',
                        end: '2019-08-17 09:30',
                        color: 'indigo',
                    },
                    {
                        name: 'Big Meeting',
                        details: 'A very important meeting about nothing',
                        start: '2019-08-19 08:00',
                        end: '2019-08-19 11:30',
                        color: 'red',
                    }
                ],
            }
        }
    }
</script>