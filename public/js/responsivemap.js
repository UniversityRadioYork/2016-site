/*global
	google
*/

$(document).ready(function () {

	//Google Maps JS
	//Set Map
	function initialize() {
		var uryLocation = new google.maps.LatLng(53.948193, -1.054030);
		var mapOptions = {
			zoom: 15,
			center: uryLocation
		};

		var map = new google.maps.Map(document.getElementById("map"), mapOptions);
		//Callout Content
		var contentString = "URY Station";
		//Set window width + content
		var infowindow = new google.maps.InfoWindow({
			content: contentString,
			maxWidth: 500
		});

		//Add Marker
		var marker = new google.maps.Marker({
			position: uryLocation,
			map,
			title: "image title"
		});

		//Resize Function
		google.maps.event.addDomListener(window, "resize", function () {
			var center = map.getCenter();
			google.maps.event.trigger(map, "resize");
			map.setCenter(center);
		});

		$("#map_container").css({
			"width": "100%",
			"margin": "0 auto",
			"margin-top": "10px"
		});

		$("#map").css({
			"height": "400",
			"overflow": "hidden",
			"position": "relative",
			"padding-bottom": "22.25%",
			"padding-top": "30px"
		});
	}

	google.maps.event.addDomListener(window, 'load', initialize);

});