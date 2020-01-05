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
        $(htmlid).append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<div class="thumbnail">' +
            '<a href="//youtube.com/watch?v=' + response.items[i].snippet.resourceId.videoId + '" target="_blank">' +
              '<img src="' + response.items[i].snippet.thumbnails.maxres.url +
              '" alt="' + response.items[i].snippet.title + '" class="img-fluid">' +
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

//Comments Character Count
document.getElementById("comments").onkeyup = function(){
  var remaining = 1000 - document.getElementById("comments").value.length;
  document.getElementById("charcount").style.display = (remaining<=100) ? "block":"none";
  document.getElementById("charcount").style.color = (remaining<=10) ? "red":"black";
  document.getElementById("charcount").innerHTML = remaining + " Characters Remaining";
};
