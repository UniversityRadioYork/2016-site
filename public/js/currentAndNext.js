/* global MyRadioAPIKey */
$(document).ready(function() {

  function scheduleUpdate() {
    // Call the function 10 seconds past every half hour
    let nextCall = new Date();

    if (nextCall.getMinutes() >= 30) {
      nextCall.setHours(nextCall.getHours() + 1);
      nextCall.setMinutes(0);
    } else {
      nextCall.setMinutes(30);
    }
    nextCall.setSeconds(10);

    let difference = nextCall - new Date();
    setTimeout(updateShow, difference);
  }

  // Used for autoupdating the now and next.
  function updateShow() {
    var data;
    $.ajax({
      dataType: "json",
      url: "//ury.org.uk/api/v2/timeslot/currentandnext/?api_key=" + MyRadioAPIKey,
      data,
      success(data) {
        var calcTime = function(timestamp) {
          var date = new Date(timestamp * 1000);
          // Hours part from the timestamp
          var hours = "0" + date.getHours();
          // Minutes part from the timestamp
          var minutes = "0" + date.getMinutes();
          // Use substr to remove the extra 0 if 2 digit hour/min
          return hours.substr(-2) + ":" + minutes.substr(-2);
        };

        var makeContent = function(show) {
          if (typeof show.title !== "undefined") {
            // If we're off air
            if (show.end_time === "The End of Time") {
              return "<h5 class=\"ellipsis\">We're off air right now</h5>" +
                "<h6>Check back next term</h6>";
            }
            // Use "Now" as start time if it's missing
            let startTimeString = calcTime(show.start_time);
            if (typeof show.start_time === "undefined") {
              startTimeString = "Now";
            }
            // Default case (regular show)
            return "<h5 class=\"ellipsis\">" +
              show.title +
              "</h5>" +
              "<h6>" + startTimeString + " - " + calcTime(show.end_time) + "</h6>";
          } else {
            return "<span>There's nothing on here</span>";
          }
        };

        // Current show
        if (typeof data.payload.current === "undefined") {
          // There is no current show; Something is probably very wrong...
          $(".current-and-next-now").replaceWith(
            "<div class=\"current-and-next-now p-2 pt-3 px-3 p-md-3 p-lg-4 \" title=\"View the show now on air.\">" +
            "<h2>Now</h2>" +
            "<h5 class=\"ellipsis\">There's nothing on right now</h5>");
          $("#studiomessage *").attr("disabled", true);
        } else if (typeof data.payload.current.url !== "undefined") {
          $(".current-and-next-now").replaceWith(
            "<a class=\"current-and-next-now p-2 pt-3 px-3 p-md-3 p-lg-4 \" href=" +
            data.payload.current.url +
            " title=\"View the show now on air.\">" +
            "<h2>Now</h2>" +
            makeContent(data.payload.current) +
            "</a>");
          $("#studiomessage *").attr("disabled", false);
        } else {
          $(".current-and-next-now").replaceWith(
            "<div class=\"current-and-next-now p-2 pt-3 px-3 p-md-3 p-lg-4 \" title=\"View the show now on air.\">" +
            "<h2>Now</h2>" +
            makeContent(data.payload.current) +
            "</a>");
          $("#studiomessage *").attr("disabled", true);
        }

        // Next show
        if (typeof data.payload.next === "undefined") {
          // There is no next show (e.g. we're off air)
          $(".current-and-next-next").replaceWith(
            "<div class=\"current-and-next-next p-2 pt-3 px-3 p-md-3 p-lg-4 \" title=\"View the show up next.\">" +
            "<h2>Next</h2>" +
            "<h5 class=\"ellipsis\">There's nothing up next yet</h5>" +
            "</a>");
        } else if (typeof data.payload.next.url !== "undefined") {
          $(".current-and-next-next").replaceWith(
            "<a class=\"current-and-next-next p-2 pt-3 px-3 p-md-3 p-lg-4 \" href=" +
            data.payload.next.url +
            " title=\"View the show up next.\">" +
            "<h2>Next</h2>" +
            makeContent(data.payload.next) +
            "</a>");
        } else {
          $(".current-and-next-next").replaceWith(
            "<div class=\"current-and-next-next p-2 pt-3 px-3 p-md-3 p-lg-4 \" title=\"View the show up next.\">" +
            "<h2>Next</h2>" +
            makeContent(data.payload.next) +
            "</a>");
        }

        $(".current-and-next-img img").attr("src", "//ury.org.uk" + data.payload.current.photo);

        //Schedule when the next update will happen.
        scheduleUpdate();
      }
    });
  }

  // Call on startup too, mainly to schedule next update.
  updateShow();

});