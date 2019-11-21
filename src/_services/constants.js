import colors from 'vuetify/lib/util/colors';

export const hostname = process.env.NODE_ENV === 'development' ? 'http://localhost:8090': '';
export const cars = ['Red Bus', 'White Bus', 'e-Auto', 'Little Red'];
export const carProperties = {
    'Red Bus': { color: colors.red.darken4, isBig: true },
    'White Bus': { color: colors.white, isBig: true },
    'e-Auto': { color: colors.blueGrey.darken4, isBig: false },
    'Little Red': { color: colors.pink.darken4, isBig: false }
};
export const brandName = 'Sammatz';