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
      $('#youtube-videos').append('<div class="thumbnail-container col-10 col-sm-7 col-md-4 col-lg-3">' +
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
$(document).ready(function () {
//Used for autoupdating the now and next.
function updateShow() {
  var data;
  $.ajax({
    dataType: "json",
    url: "//ury.org.uk/api/v2/timeslot/currentandnext/?api_key=" + MyRadioAPIKey,
    data: data,
    success: function(data) {
      var calcTime = function(timestamp) {
        var date = new Date(timestamp * 1000);
        // Hours part from the timestamp
        var hours = date.getHours();
        // Minutes part from the timestamp
        var minutes = "0" + date.getMinutes();

        return hours + ':' + minutes.substr(-2)
      }
      var makeContent = function(show) {
        if (typeof show.title !== 'undefined') {
          return "<h5 class=\"ellipsis\">"
          + show.title
          + "</h5>"
          + "<h6>" + calcTime(show.start_time) +" - " + calcTime(show.end_time) + "</h6>";
        } else {
          return "<span>Looks like there is nothing on here.</span>"
        }
      }

      if (typeof data.payload.current.url !== 'undefined') {

        $(".current-and-next-now").replaceWith(
          "<a class=\"current-and-next-now p-2 pt-3 px-3 p-md-3 p-lg-4 \" href="
            + data.payload.current.url
            + " title=\"View the show now on air.\">"
            + "<h2>Now</h2>"
            + makeContent(data.payload.current)
            + "</a>");

        $("#studiomessage *").attr('disabled', false)
      } else {
        $(".current-and-next-now").replaceWith(
          "<div class=\"current-and-next-now p-2 pt-3 px-3 p-md-3 p-lg-4 \" title=\"View the show now on air.\">"
            + "<h2>Now</h2>"
            + makeContent(data.payload.current)
            + "</a>");
        $("#studiomessage *").attr('disabled', true)
      }
      // Next show
      if (typeof data.payload.next.url !== 'undefined') {
        $(".current-and-next-next").replaceWith(
          "<a class=\"current-and-next-next p-2 pt-3 px-3 p-md-3 p-lg-4 \" href="
            + data.payload.nexturl
            + " title=\"View the show up next.\">"
            + "<h2>Next</h2>"
            + makeContent(data.payload.next)
            + "</a>");
      } else {
        $(".current-and-next-next").replaceWith(
          "<div class=\"current-and-next-next p-2 pt-3 px-3 p-md-3 p-lg-4 \" title=\"View the show up next.\">"
            + "<h2>Next</h2>"
            + makeContent(data.payload.next)
            + "</a>");
      }

      $(".current-and-next-img img").attr('src', "//ury.org.uk" + data.payload.current.photo);
    }
  });
  setInterval(updateShow, 300000);
}

updateShow();
});
