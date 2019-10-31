export const hostname = process.env.NODE_ENV === 'development' ? 'http://localhost:8090': '';
export const cars = ['Red Bus', 'White Bus', 'e-Auto', 'Little Red'];
export const carProperties = {
    'Red Bus': { color: '#F44336', isBig: true },
    'White Bus': { color: 'white', isBig: true },
    'e-Auto': { color: 'black', isBig: false },
    'Little Red': { color: '#EF5350', isBig: false }
};
export const brandName = 'Sammatz';