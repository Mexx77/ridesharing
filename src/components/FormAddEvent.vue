<template>
    <v-row justify="center">
        <v-dialog :value="isActive" persistent max-width="600px">
            <v-card>
                <v-card-title>
                    <span class="headline">Neue Reservierungsanfrage</span>
                </v-card-title>
                <v-card-text>
                    <v-container>
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
                        <v-row justify="space-around" align="center">
                            <v-col cols="12" sm="6">
                                <v-time-picker :width="260" format="24hr" v-model="start" :max="end"></v-time-picker>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-time-picker :width="260" format="24hr" v-model="end" :min="start"></v-time-picker>
                            </v-col>
                        </v-row>
                    </v-container>
                    <small>*indicates required field</small>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" text @click="close">Close</v-btn>
                    <v-btn color="blue darken-1" text @click="close">Save</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>
    </v-row>
</template>

<script>
    export default {
        props: ['isActive'],
        methods: {
            close: function() {
                this.$emit('update:isActive', false)
            }
        },
        data: function () {
            return {
            }
        }
    }
</script>