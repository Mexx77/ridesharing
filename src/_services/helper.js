export function getEventColor(event) {
    return event.confirmed && event.carColor ? event.carColor : 'grey';
}

export function getEventTextColor(event) {
    if (event.carColor) {
        return event.carColor === 'white' ? 'secondary' : 'white'
    } else {
        return 'white';
    }
}