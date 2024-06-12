var map = L.map('map').setView([Lat, Lng], 15);
L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map);
var marker = L.marker([Lat, Lng]).addTo(map);
marker.bindPopup("University Radio York").openPopup();
