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