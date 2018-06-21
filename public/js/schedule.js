// global StartHour
function zpad(n, width) {
  n = n + "";
  if (n.length >= width){
    return n
  } else {
    return new Array(width - n.length + 1).join("0") + n;
  }
}

function jumpToNow(disableMove=false){
  let daysOfWeek = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
  let d = new Date();
  let weekday = d.getDay();
  let hour = d.getHours();
  if (hour < StartHour){
    weekday -= 1;
    if (weekday < 0){
      weekday = 6;
    }
  }
  hour = zpad(hour,2);
  weekday = daysOfWeek[weekday]
  let cell = $(".day-" + weekday + " .hour-" + hour);
  console.log(".day-" + weekday + " .hour-" + hour, cell)
  if(cell.length == 1){
    if(!disableMove){
      $(window).scrollTop(Math.max(cell.offset().top - 100, 0));
    }
  } else {
    if(cell.length == 0){
      $("#jumpToNow").attr('disabled', true)
      $("#jumpToNow").text("No show on air right now!");
      setTimeout(function(){
        $("#jumpToNow").attr('disabled', false);
        $("#jumpToNow").text("Jump to current show");
      }, 3000);

    } else {
      console.error("Found multiple cells matching current time");
    }
  }
}

$(function(){
  jumpToNow(true);
})