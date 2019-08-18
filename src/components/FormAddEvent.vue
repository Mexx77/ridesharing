<template>
    <v-row justify="center">
        <v-dialog :value="isActive" @update:return-value="close" max-width="600px">
            <v-card>
                <v-toolbar color="indigo" dark>
                    <v-toolbar-title><v-icon>mdi-car</v-icon> Reservierung am {{date}}</v-toolbar-title>
                </v-toolbar>
                <v-card-text>
                    <v-container pa-0>
                        <v-row>
                            <v-col cols="12" sm="6" md="4">
                                <v-text-field label="Name*" required></v-text-field>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-select
                                        prepend-icon="mdi-car"
                                        :items="['Red Bus', 'White Bus', 'e-Auto', 'Little Red']"
                                        label="Auto*"
                                        required
                                ></v-select>
                            </v-col>
                        </v-row>
                        <v-row justify="space-between" align="center">
                            <v-col cols="12" sm="6">
                                <v-label mb-10>Uhrzeit Start</v-label>
                                <v-time-picker
                                        :value="startTime"
                                        color="indigo"
                                        :width="265"
                                        format="24hr"
                                        v-model="startTime"
                                ></v-time-picker>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-label>Uhrzeit Rückgabe</v-label>
                                <v-time-picker
                                        :value="end"
                                        color="indigo"
                                        :width="265"
                                        format="24hr"
                                        :min="startTime"
                                ></v-time-picker>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" text @click="close">Abbrechen</v-btn>
                    <v-btn color="blue darken-1" text @click="close">Anfragen</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </v-row>
</template>

<script>
    export default {
        props: ['isActive', 'date', 'startTime'],
        methods: {
            close: function() {
                this.$emit('update:is-active', false)
            }
        },
        data: () => ({
            end: ''
        })
    }
</script>