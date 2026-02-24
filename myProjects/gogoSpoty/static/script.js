async function updateTrack() {
    const res = await fetch('/api/current');
    const data = res.json();

    document.getElementById('cover').src = data.Item.Album.Images[0].URL;
    document.getElementById('track-name').textContent = data.Item.Name;
    // add loop for multiple singers
    document.getElementById('artist').textContent = data.Item.Artists[0].Name;

    const current = Math.floor(data.Progress / 1000);
    const total =  Math.floor(data.Timestamp / 1000);
    const percent = (data.Progress / data.Timestamp) * 100;

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