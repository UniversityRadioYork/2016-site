/*global google:true placeName*/

$(document).ready(function () {
  //Google Maps JS
  //Set Map
  function initialize() {
    var uryLocation = new google.maps.LatLng(Lat, Lng);
    var mapOptions = {
      zoom: 15,
      center: uryLocation
    };

    var map = new google.maps.Map(document.getElementById("map"), mapOptions);
    //Callout Content
    var contentString = placeName;
    //Set window width + content
    var infowindow = new google.maps.InfoWindow({
      content: contentString,
      maxWidth: 500
    });

    //Add Marker
    var marker = new google.maps.Marker({
      position: uryLocation,
      map,
      title: placeName
    });

    //Resize Function
    google.maps.event.addDomListener(window, "resize", function () {
      var center = map.getCenter();
      google.maps.event.trigger(map, "resize");
      map.setCenter(center);
    });
  }

  google.maps.event.addDomListener(window, "load", initialize);
});