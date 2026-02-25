let lastTimestamp
let progress
let duration
let playing

document.addEventListener('DOMContentLoaded', function() {

async function updateTrack() {
    const res = await fetch('/api/current');
    const data = await res.json();

    if (!data.item) {
        document.querySelector('.player').classList.add('hidden');
        return;
    }

    document.querySelector('.player').classList.remove('hidden');

    lastTimestamp = Date.now()
    progress = data.progress_ms;
    duration = data.item.duration_ms;
    playing = data.is_playing;

    const s = data.item.artists
    .map(artist => artist.name)
    .join(", ");
    
    console.log('Updating...', data);
    document.getElementById('artist').textContent = s;
    document.getElementById('progress').style.background = "red" //data.accent_color; next iter
    document.getElementById('track-name').textContent = data.item.name;
    document.getElementById('cover').src = data.item.album.images[0].url;

    const total =  Math.floor(duration / 1000);
    document.getElementById('total-time').textContent = formatTime(total);    
}

function updateProgressBar() {
    const now = Date.now();
    const elapsed = playing ? now - lastTimestamp : 0;

    const rawProgress = progress + elapsed;
    const clampedProgress = Math.min(rawProgress, duration);

    const percent = (clampedProgress / duration) * 100;
    const timeSeconds = Math.floor(clampedProgress / 1000);

    document.getElementById('current-time').textContent =
        formatTime(timeSeconds);

    document.getElementById('progress').style.width =
        percent + '%';
}

function formatTime(seconds) {
    const mins = Math.floor(seconds/60); 
    const secs = seconds % 60;
    return mins + ':' + (secs < 10 ? "0" : '') + secs;
}

function tick() {
    updateProgressBar();
    requestAnimationFrame(tick);
}

    updateTrack();
    setInterval(updateTrack, 5000);
    requestAnimationFrame(tick);

    console.log("updated")
});
