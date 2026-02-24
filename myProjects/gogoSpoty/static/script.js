async function updateTrack() {
    const res = await fetch('/api/current');
    const data = res.json();

    document.getElementById('cover').src = data.item.album.images[0].url;
    document.getElementById('track-name').textContent = data.item.name;
    // add loop for multiple singers
    document.getElementById('artist').textContent = data.item.artists[0].name;

    const current = Math.floor(data.progress_ms / 1000);
    const total =  Math.floor(data.item.duration_ms / 1000);
    const percent = (data.progress_ms / data.item.duration_ms) * 100;

    document.getElementById('current-time').textContent = formatTime(current);
    document.getElementById('total-time').textContent = formatTime(total);
    document.getElementById('progress').style.width = percent + '%';
}

function formatTime(seconds) {
    const mins = Math.floor(secs/60);
    const secs = seconds % 60;
    return mins + ':' + (secs < 10 ? "0" : '') + secs;
}

updateTrack();
setInterval(updateTrack, 5000);