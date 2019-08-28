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
                                :color="selectedEvent.getEventColor"
                                :style="{color: selectedEvent.getEventTextColor}"
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
                        <v-form
                                ref="form"
                                v-model="formIsValid"
                                lazy-validation
                        >
                            <v-card>
                                <v-toolbar color="primary" dark>
                                    <v-toolbar-title>
                                        <v-icon class="pb-1">mdi-car</v-icon>
                                        Reservierung am {{germanDate}}
                                    </v-toolbar-title>
                                </v-toolbar>
                                <v-card-text class="pb-0">
                                    <v-container pa-0>
                                        <v-row>
                                            <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                                <v-text-field
                                                        v-model="driver"
                                                        prepend-icon="mdi-account"
                                                        label="Fahrer*"
                                                        :rules="[v => !!v || 'Name des Fahrer benötigt']"
                                                        required
                                                ></v-text-field>
                                            </v-col>
                                            <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                                <v-combobox
                                                        v-model="destination"
                                                        prepend-icon="mdi-city"
                                                        :items="['Lüneburg', 'Dannenberg', 'Hitzacker']"
                                                        label="Fahrtziel*"
                                                        :rules="[v => !!v || 'Fahrtziel benötigt']"
                                                        required
                                                ></v-combobox>
                                            </v-col>
                                        </v-row>
                                        <v-row>
                                            <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                                <v-menu
                                                        v-model="menuStartTime"
                                                        :close-on-content-click="true"
                                                        :nudge-right="40"
                                                        transition="scale-transition"
                                                        offset-y
                                                        full-width
                                                        min-width="272px"
                                                >
                                                    <template v-slot:activator="{ on }">
                                                        <v-text-field
                                                                v-model="startTime"
                                                                label="Startzeit"
                                                                prepend-icon="mdi-calendar-clock"
                                                                readonly
                                                                v-on="on"
                                                                :rules="[v => !!v || 'Startzeit benötigt']"
                                                        ></v-text-field>
                                                    </template>
                                                    <v-time-picker
                                                            v-model="startTime"
                                                            color="primary"
                                                            :width="272"
                                                            format="24hr"
                                                            required
                                                            :allowed-minutes="allowedMinutes"
                                                    ></v-time-picker>
                                                </v-menu>

                                            </v-col>
                                            <v-col :cols="$vuetify.breakpoint.mdAndUp ? 6 : 12">
                                                <v-menu
                                                        v-model="menuEndTime"
                                                        :close-on-content-click="true"
                                                        :nudge-right="40"
                                                        transition="scale-transition"
                                                        offset-y
                                                        full-width
                                                        min-width="272px"
                                                >
                                                    <template v-slot:activator="{ on }">
                                                        <v-text-field
                                                                v-model="endTime"
                                                                label="Zeit der Rückgabe"
                                                                prepend-icon="mdi-calendar-clock"
                                                                readonly
                                                                v-on="on"
                                                                :rules="[v => !!v || 'Zeit der Rückgabe benötigt']"
                                                        ></v-text-field>
                                                    </template>
                                                    <v-time-picker
                                                            v-model="endTime"
                                                            color="primary"
                                                            :width="272"
                                                            format="24hr"
                                                            :min="startTime"
                                                            required
                                                            :allowed-minutes="allowedMinutes"
                                                    ></v-time-picker>
                                                </v-menu>
                                            </v-col>
                                        </v-row>
                                    </v-container>
                                </v-card-text>
                                <v-card-actions class="mr-2 pb-4 pt-0">
                                    <v-container pt-0>
                                        <v-row dense>
                                            <v-col>
                                                <v-switch
                                                        class="mt-0"
                                                        v-model="bigCarNeeded"
                                                        label="Ich brauche ein großes Auto"
                                                ></v-switch>
                                                <v-btn @click="showAddEventForm = false">Abbrechen</v-btn>
                                                <v-btn @click="validateAndSubmitForm">Anfragen</v-btn>
                                            </v-col>
                                        </v-row>
                                    </v-container>

                                </v-card-actions>
                            </v-card>
                        </v-form>
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
            germanDate() {
                if (this.focus == undefined) return undefined
                const p = this.focus.split(/\D/g)
                return [p[2], p[1], p[0]].join(".")
            }
        },
        methods: {
            allowedMinutes: m => m % 15 === 0,
            roundMinutes(hour, minute) {
                const m = (((minute + 7.5) / 15 | 0) * 15) % 60
                const h = ((((minute / 105) + .5) | 0) + hour) % 24
                return h + ':' + m
            },
            addEvent(time) {
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
                return event.carColor
            },
            getEventTextColor(event) {
                return event.carColor == 'white' ? 'secondary' : 'white'
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
            },
            validateAndSubmitForm() {
                if (this.$refs.form.validate()) {
                    this.showAddEventForm = false;

                }
            }
        },
        mounted: function () {
            this.$http
                //.get('https://ridesharing-0df0.restdb.io/rest/rides')
                .get('http://localhost:8090/rides')
                .then((response) => {
                    this.events = response.data
                });
        },
        data() {
            return {
                formIsValid: false,
                driver: '',
                destination: '',
                bigCarNeeded: false,
                showAddEventForm: false,
                startTime: '12:00',
                menuStartTime: false,
                endTime: null,
                menuEndTime: false,
                time: null,
                today: this.formatDate(new Date()),
                focus: this.today,
                type: this.$vuetify.breakpoint.smAndDown ? 'day' : '4day',
                typeToLabel: {
                    month: 'Monat',
                    week: 'Woche',
                    day: 'Tag',
                    '4day': '4 Tage',
                },
                selectedEvent: {},
                selectedElement: null,
                selectedOpen: false,
                events: []
            }
        }
    }
</script>

<style>
    .v-event-timed {
        font-size: 16px !important;
    }
</style>