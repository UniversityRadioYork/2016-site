//Youtube slideshow for index page
function onGoogleLoad() {
  gapi.client.setApiKey(youtubeAPIKey);
  gapi.client.load('youtube', 'v3', function() {
  var request = gapi.client.youtube.playlistItems.list({
    part: 'snippet',
    playlistId: youtubePlaylistID,
    maxResults: 12
  });

  request.execute(function(response) {
    for (var i = 0; i < response.items.length; i++) {
      $('#youtube-video-slider').append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
          '<div class="thumbnail">' +
            '<a href="//youtube.com/watch?v=' + response.items[i].snippet.resourceId.videoId + '" target="_blank">' +
              '<img src="' + response.items[i].snippet.thumbnails.maxres.url +
              '" alt="' + response.items[i].snippet.title + '" class="img-fluid">' +
            '</a>' +
          '</div>' +
        '</div>');
      }
    });
  });
}
