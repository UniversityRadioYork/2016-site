/* global StartHour */
function zpad(n, width) {
  n = n + "";
  if (n.length >= width) {
    return n;
  } else {
    return new Array(width - n.length + 1).join("0") + n;
  }
}

function jumpToNow(disableMove = false) {
  let daysOfWeek = {
    0: "Sunday",
    1: "Monday",
    2: "Tuesday",
    3: "Wednesday",
    4: "Thursday",
    5: "Friday",
    6: "Saturday"
  };
  let d = new Date();
  let weekdayIndex = d.getDay();
  let hour = d.getHours();
  if (hour < StartHour) {
    weekdayIndex -= 1;
    if (weekdayIndex < 0) {
      weekdayIndex = 6;
    }
  }
  hour = zpad(hour, 2);
  let weekday = daysOfWeek[weekdayIndex];
  let selector = ".day-" + weekday + " .hour-" + hour;
  let cell = $(selector);
  if (cell.length === 1) {
    if (!disableMove) {
      $(window).scrollTop(Math.max(cell.offset().top - 200, 0));
      $(selector).animate({opacity: 0}, 500, "swing", function() {
        $(selector).animate({opacity: 1}, 500, "linear");
      });
    }
  } else if (cell.length === 0) {
    $("#jumpToNow").attr("disabled", true);
    $("#jumpToNow").text("No show on air right now!");
    setTimeout(function() {
      $("#jumpToNow").attr("disabled", false);
      $("#jumpToNow").text("Jump to current show");
    }, 3000);
  }
}

$(function() {
  jumpToNow(true);
});
