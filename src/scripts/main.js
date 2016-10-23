   
    	$(function(){
    		// On page load / resize, force now & next header boxes to be equal (sqaure) size.
			var $el = $('.now-next-live');
			var $el2 = $('.now-next-now');
			var $el3 = $('.now-next-next');
			var $window = $(window).on('resize', function(){
			   var width = $el.outerWidth();
			   $el.height(width);
			   $el2.height(width);
			   $el3.height(width);
			}).trigger('resize'); //on page load

		});

