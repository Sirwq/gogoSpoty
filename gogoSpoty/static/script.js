const state = { lastTimestamp: 0, progress: 0, duration: 0, playing: false };
const UPDATE_INTERVAL = 5000;

function tick() {
    updateProgressBar();
    requestAnimationFrame(tick);
}

document.addEventListener('DOMContentLoaded', function() {
    updateTrack();
    setInterval(updateTrack, UPDATE_INTERVAL);
    requestAnimationFrame(tick);
});

async function updateTrack() {
    const data = await fetchData();
    updateUI(data);
}

function updateUI(data) {

    if (!data?.item) {
        handleVisibility(false);
        return;
    }

    state.lastTimestamp = Date.now()
    state.progress = data.progress_ms;
    state.duration = data.item.duration_ms;
    state.playing = data.is_playing;
    
    handleVisibility(state.playing);
    const s = data.item.artists
    .map(artist => artist.name)
    .join(", ");
    
    document.getElementById('artist').textContent = s;
    document.getElementById('progress').style.background = "red" //data.accent_color; next iter
    document.getElementById('track-name').textContent = data.item.name;

    if (data.item?.album?.images?.length > 0) {
        document.getElementById('cover').src = data.item.album.images[0].url;
    }

    const total =  Math.floor(state.duration / 1000);
    document.getElementById('total-time').textContent = formatTime(total);   
    
    console.log('Updating...', data);
}

function handleVisibility(shouldShow) {
    const player = document.querySelector('.player');
    player.classList.toggle('hidden', !shouldShow);
}

function updateProgressBar() {
    const now = Date.now();
    const elapsed = state.playing ? now - state.lastTimestamp : 0;
    const rawProgress = state.progress + elapsed;
    const clampedProgress = Math.min(rawProgress, state.duration);
    const percent = (clampedProgress / state.duration) * 100;
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

async function fetchData() {
    try {
        const res = await fetch('/api/current');

        if (!res.ok) {
            throw new Error (`HTTP error ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error('Fetch failed:', error);
        return null;
    }
}


