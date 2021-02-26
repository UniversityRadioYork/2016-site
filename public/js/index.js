function getYoutubeFeed(playlistid, results, htmlid) {
  gapi.client.setApiKey(youtubeAPIKey);
  gapi.client.load("youtube", "v3", function() {
    var request = gapi.client.youtube.playlistItems.list({
      part: "snippet",
      playlistId: playlistid,
      maxResults: results
    });

    request.execute(function(response) {
      for (var i = 0; i < response.items.length; i++) {
        var thumb;
        if (response.items[i].snippet.thumbnails.maxres != undefined) {
          thumb = response.items[i].snippet.thumbnails.maxres.url;
        } else {
          thumb = response.items[i].snippet.thumbnails.standard.url;
        }
        $(htmlid).append(
          '<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
            '<div class="thumbnail">' +
            '<a href="//youtube.com/watch?v=' +
            response.items[i].snippet.resourceId.videoId +
            '" target="_blank" rel="noopener noreferrer">' +
            '<img src="' +
            thumb +
            '" alt="' +
            response.items[i].snippet.title +
            '" class="img-fluid">' +
            "</a>" +
            "</div>" +
            "</div>"
        );
      }
      if (isIndex && htmlid === "#sessions-videos") {
        $(htmlid).append(
          '<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
            "<a class=\"ury-card sessions link\" href='/ontap/'>" +
            '<div class="ury-card-body">' +
            '<div class="ury-card-lg-title">See more...</div>' +
            "</div>" +
            "</a>" +
            "</div>"
        );
      }
      if (isOD && htmlid === "#sessions-videos") {
        $(htmlid).append(
          '<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
            '<a class="ury-card sessions link" href="' +
            youtubeLink +
            '">' +
            '<div class="ury-card-body">' +
            '<div class="ury-card-lg-title">View more on Youtube...</div>' +
            "</div>" +
            "</a>" +
            "</div>"
        );
      }
      if (isIndex && htmlid === "#cin-videos") {
        $(htmlid).append(
          '<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
            "<a class=\"ury-card cin link\" href='/cin/'>" +
            '<div class="ury-card-body">' +
            '<div class="ury-card-lg-title">See more...</div>' +
            "</div>" +
            "</a>" +
            "</div>"
        );
      }
      if (isCIN && !isIndex && htmlid === "#cin-videos") {
        $(htmlid).append(
          '<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
            '<a class="ury-card cin link" href="' +
            youtubeLink +
            '">' +
            '<div class="ury-card-body">' +
            '<div class="ury-card-lg-title">View more on Youtube...</div>' +
            "</div>" +
            "</a>" +
            "</div>"
        );
      }
    });
  });
}

// Countdown Timer
function countdown() {
  if (window.isCountdown) {
    const now = new Date();
    const istorn2020 = new Date("2021-03-02T18:00:00Z");

    const diffSeconds = (istorn2020 - now) / 1000;
    var timerSeconds = (diffSeconds % 60).toFixed(0).padStart(2, "0");
    var timerMinutes = Math.floor((diffSeconds % 3600) / 60)
      .toFixed(0)
      .padStart(2, "0");
    var timerHours = Math.floor((diffSeconds % 86400) / 3600)
      .toFixed(0)
      .padStart(2, "0");
    var timerDays = Math.floor(diffSeconds / 86400).toFixed(0);

    if (timerSeconds == 60) {
      timerSeconds = 0;
      timerMinutes++;
    }

    if (timerMinutes == 60) {
      timerMinutes = 0;
      timerHours++;
    }

    if (timerDays < 0) {
      timerDays = 0;
      timerMinutes = 0;
      timerHours = 0;
      timerSeconds = 0;
    }

    document.getElementById("countdownDays").innerText = timerDays;
    document.getElementById("countdownHours").innerText = timerHours;
    document.getElementById("countdownMinutes").innerText = timerMinutes;
    document.getElementById("countdownSeconds").innerText = timerSeconds;

    window.setTimeout(countdown, 1000);
  }
}
countdown();

//Youtube slideshow for index page
function onGoogleLoad() {
  if (isIndex) {
    getYoutubeFeed(youtubeSessionsPlaylistID, 7, "#sessions-videos");
  }
  if (isOD) {
    getYoutubeFeed(youtubeSessionsPlaylistID, 15, "#sessions-videos");
  }
  if (isIndex) {
    getYoutubeFeed(youtubeCINPlaylistID, 7, "#cin-videos");
  }
  if (isCIN && !isIndex) {
    getYoutubeFeed(youtubeCINPlaylistID, 15, "#cin-videos");
  }
}
