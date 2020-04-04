function getYoutubeFeed(playlistid, results, htmlid)
{
  gapi.client.setApiKey(youtubeAPIKey);
  gapi.client.load('youtube', 'v3', function() {

    var request = gapi.client.youtube.playlistItems.list({
      part: 'snippet',
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
        $(htmlid).append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<div class="thumbnail">' +
            '<a href="//youtube.com/watch?v=' + response.items[i].snippet.resourceId.videoId + '" target="_blank">' +
              '<img src="' + thumb + '" alt="' + response.items[i].snippet.title + '" class="img-fluid">' +
            '</a>' +
          '</div>' +
        '</div>');
      }
      if(isIndex && htmlid === "#sessions-videos") {
        $(htmlid).append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<a class="ury-card sessions link" href=\'/ontap/\'>' +
            '<div class="ury-card-body">' +
              '<div class="ury-card-lg-title">See more...</div>' +
            '</div>' +
          '</a>' +
        '</div>');
      }
      if(isOD && htmlid === "#sessions-videos") {
        $(htmlid).append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<a class="ury-card sessions link" href=\"' + youtubeLink + '\">' +
            '<div class="ury-card-body">' +
              '<div class="ury-card-lg-title">View more on Youtube...</div>' +
            '</div>' +
          '</a>' +
        '</div>');
      }
      if(isIndex && htmlid === "#cin-videos") {
        $(htmlid).append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<a class="ury-card cin link" href=\'/cin/\'>' +
            '<div class="ury-card-body">' +
              '<div class="ury-card-lg-title">See more...</div>' +
            '</div>' +
          '</a>' +
        '</div>');
      }
      if(isCIN && htmlid === "#cin-videos") {
        $(htmlid).append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<a class="ury-card cin link" href=\"' + youtubeLink + '\">' +
            '<div class="ury-card-body">' +
              '<div class="ury-card-lg-title">View more on Youtube...</div>' +
            '</div>' +
          '</a>' +
        '</div>');
      }
    });

  });
}

const isToday = () => {
  // 12 feb 2020
  someDate = new Date('2020-02-12')
  let today = new Date();
  today = someDate.getDate() == today.getDate() &&
    someDate.getMonth() == today.getMonth() &&
    someDate.getFullYear() == today.getFullYear();
  return today;
};

// CIN countdown
function cinCounter() {
  if (isCIN) {
    const now = new Date();
    const cin = new Date("2020-02-12T19:00:00Z");

    const diffSeconds = (cin - now) / 1000;
    const timerSeconds = (diffSeconds % 60).toFixed(0).padStart(2, "0");
    const timerMinutes = Math.floor(diffSeconds % 3600 / 60).toFixed(0).padStart(2, "0");
    const timerHours = Math.floor(diffSeconds % 86400 / 3600).toFixed(0).padStart(2, "0");
    const timerDays = Math.floor(diffSeconds / 86400).toFixed(0);

    if (false) {
      document.getElementById("cinCountdown").innerText = "" + timerDays + " days " + timerHours + ":" + timerMinutes + ":" + timerSeconds;
    } else {
      document.getElementById("cinCountdownDays").innerText = timerDays;
      document.getElementById("cinCountdownHours").innerText = timerHours;
      document.getElementById("cinCountdownMinutes").innerText = timerMinutes;
      document.getElementById("cinCountdownSeconds").innerText = timerSeconds;
    }

    window.setTimeout(cinCounter, 1000);
  }
}
cinCounter();

// ISTORN 2020 countdown
function istorn2020Counter() {
  if (isISTORN2020) {
    const now = new Date();
    const istorn2020 = new Date("2020-04-13T11:00:00Z");

    const diffSeconds = (istorn2020 - now) / 1000;
    const timerSeconds = (diffSeconds % 60).toFixed(0).padStart(2, "0");
    const timerMinutes = Math.floor(diffSeconds % 3600 / 60).toFixed(0).padStart(2, "0");
    const timerHours = Math.floor(diffSeconds % 86400 / 3600).toFixed(0).padStart(2, "0");
    const timerDays = Math.floor(diffSeconds / 86400).toFixed(0);

    document.getElementById("istornCountdownDays").innerText = timerDays;
    document.getElementById("istornCountdownHours").innerText = timerHours;
    document.getElementById("istornCountdownMinutes").innerText = timerMinutes;
    document.getElementById("istornCountdownSeconds").innerText = timerSeconds;

    window.setTimeout(istorn2020Counter, 1000);
  }
}
istorn2020Counter();

//Youtube slideshow for index page
function onGoogleLoad() {

  if(isIndex) {
    getYoutubeFeed(youtubeSessionsPlaylistID, 7, "#sessions-videos");
  }
  if(isOD) {
    getYoutubeFeed(youtubeSessionsPlaylistID, 15, "#sessions-videos");
  }
  if(isIndex) {
    getYoutubeFeed(youtubeCINPlaylistID, 7, "#cin-videos");
  }
  if(isCIN) {
    getYoutubeFeed(youtubeCINPlaylistID, 15, "#cin-videos");
  }
}
