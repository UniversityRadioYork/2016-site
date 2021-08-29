function hideParentVideoContainer(htmlid) {
  $(htmlid).closest(".container-fluid").remove(); // Delete the box if we couldn't load.
}


function getYoutubeFeed(playlistid, results, htmlid) {
  if (!playlistid || !youtubeAPIKey) {
    console.error(`Failed to get YouTube videos for ${htmlid}, missing playlistid (${playlistid}) or youtubeAPIKey (${youtubeAPIKey}).`);
    hideParentVideoContainer(htmlid);
    return;
  }
  console.log(`Getting YouTube videos for ${htmlid}, with playlistid ${playlistid}.`);
  gapi.client.setApiKey(youtubeAPIKey);
  gapi.client.load("youtube", "v3", function() {
    if (!gapi.client.youtube) {
      console.error("Failed to init YouTube API.");
      hideParentVideoContainer(htmlid);
      return;
    }
    var request = gapi.client.youtube.playlistItems.list({
      part: "snippet",
      playlistId: playlistid,
      maxResults: results
    });

    request.execute(function(response) {
      if (!response || !response.items) {
        console.log(`Failed to get YouTube videos for ${htmlid} with playlistid ${playlistid}.`);
        hideParentVideoContainer(htmlid);
        return;
      }
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

if (window.isCountdown) {
  if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
    document.getElementById("index-countdown-video").autoplay = false;
    document.getElementById("index-countdown-video").pause();
  }
}

// Countdown Timer
function countdown() {
  if (window.isCountdown) {
    const now = new Date();
    const countTo = new Date("2021-03-02T18:00:00Z");

    const diffSeconds = (countTo - now) / 1000;
    var timerSeconds = (diffSeconds % 60);
    var timerMinutes = Math.floor((diffSeconds % 3600) / 60)
    var timerHours = Math.floor((diffSeconds % 86400) / 3600)
    var timerDays = Math.floor(diffSeconds / 86400);

    if (timerSeconds.toFixed(0) == "60") {
      timerSeconds = 0;
      timerMinutes++;
    }

    if (timerMinutes.toFixed(0) == "60") {
      timerMinutes = 0;
      timerHours++;
    }

    if (timerDays < 0) {
      timerDays = 0;
      timerMinutes = 0;
      timerHours = 0;
      timerSeconds = 0;
    }

    document.getElementById("countdownDays").innerText = timerDays.toFixed(0);
    document.getElementById("countdownHours").innerText = timerHours.toFixed(0).padStart(2, "0");
    document.getElementById("countdownMinutes").innerText = timerMinutes.toFixed(0).padStart(2, "0");
    document.getElementById("countdownSeconds").innerText = timerSeconds.toFixed(0).padStart(2, "0");

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

//Comments Character Count

function updateMessageboxCharacterCount() {
  var remaining =
    document.getElementById("comments").maxLength -
    document.getElementById("comments").value.length;
  document.getElementById("charcount").style.color =
    remaining <= 10 ? "red" : "black";
  document.getElementById(
    "charcount"
  ).innerHTML = `${remaining} characters remaining`;
}

document.getElementById("comments").onkeyup = () => {
  updateMessageboxCharacterCount();
};

$(document).ready(() => {
  updateMessageboxCharacterCount();
});
