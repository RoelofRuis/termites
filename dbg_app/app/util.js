export const msToTime = (duration) => {
    let millis = Math.floor(duration % 1000),
        seconds = Math.floor((duration / 1000) % 60),
        minutes = Math.floor((duration / (1000 * 60)) % 60),
        hours = Math.floor((duration / (1000 * 60 * 60)) % 24)

    hours = (hours < 10) ? '0' + hours : hours;
    minutes = (minutes < 10) ? '0' + minutes : minutes;
    seconds = (seconds < 10) ? '0' + seconds : seconds;
    millis = (millis < 10) ? '00' + millis : (millis < 100 ? '0' + millis : millis)

    return hours + ':' + minutes + ':' + seconds + "." + millis;
}