let lastTimestamp
let progress
let duration

document.addEventListener('DOMContentLoaded', function() {

async function updateTrack() {
    const res = await fetch('/api/current');
    const data = await res.json();

    if (!data.item) {
        return;
    }

    lastTimestamp = Date.now()
    progress = data.progress_ms;
    duration = data.item.duration_ms;

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

    const elapsed = now - lastTimestamp;
    const percent = ((progress + elapsed) / duration) * 100;
    const current = Math.floor((progress + elapsed) / 1000);

    document.getElementById('current-time').textContent = formatTime(current);
    document.getElementById('progress').style.width = percent + '%';
}

function formatTime(seconds) {
    const mins = Math.floor(seconds/60); 
    const secs = seconds % 60;
    return mins + ':' + (secs < 10 ? "0" : '') + secs;
}


    updateTrack();
    setInterval(updateTrack, 5000);
    setInterval(updateProgressBar, 100);
    console.log("updated")
});
