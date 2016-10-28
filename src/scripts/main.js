$(function() {
  // On page load / resize, force now & next header boxes to be equal (sqaure) size.
  var $el = $('.current-next-live');
  var $el2 = $('.current-next-now');
  var $el3 = $('.current-next-next');
  var $window = $(window).on('resize', function() {
    var width = $el.outerWidth();
    $el.height(width);
    $el2.height(width);
    $el3.height(width);
  }).trigger('resize'); //on page load
});


//Used for autoupdating the homepage Now & Next.
function updateShow() {
  $.getJSON("//ury.org.uk/api/Timeslot/getCurrentAndNext/?api_key=VMRSpUZ3uB8q4TnbNhr1oP",
    function (data) {
      if (data.status === "OK") {
        if (typeof data.payload.current.url !== 'undefined') {
          $('#current-next-now h3').html('<a href=' + data.payload.current.url + '>' + data.payload.current.title + '</a>');
          // $('#box-message').html('<h1>Send A Message</h1><span>Via the website, or text the studio</span><form action="https://ury.org.uk/schedule/message-current-show/" name="message" method="post"><textarea id="comments" name="comments" cols="30" rows="3" placeholder="Type your message here"></textarea><input type="submit" value="Send Message"></form><dl><dt>Text</dt><dd>07851 101 313</dd></dl>');
        } else {
          $('#current-next-now h3').text(data.payload.current.title);
         //  $('#box-message').html('<h1>Send A Message</h1><p><i class="icon-ban-circle icon-3x pull-left icon-border"></i>Show messaging will be available when the next show starts.</p>');
        }

        if (typeof data.payload.next.url !== 'undefined') {
          $('#current-next-next h3').html('<a href=' + data.payload.next.url + '>' + data.payload.next.title + '</a>');
        } else {
          $('#current-next-next h3').text(data.payload.next.title);
        }

        $('#current-next-live').attr('background-image',"url('https://ury.org.uk" + data.payload.current.photo) + "')";

      } else {
        console.log("API Error: " + data.status + " - " + data.payload);
      }
    });
}
$(function() {
  updateShow();
  setInterval(updateShow, 300000);
});