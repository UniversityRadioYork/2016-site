function jumpToNow(disableMove=false){
	let daysOfWeek = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
	let d = new Date();
	let weekday = daysOfWeek[d.getDay()];
	let hour = d.getHours();
	let cell = $(".day-" + weekday + " .hour-" + hour);
	if(cell.length == 1){
		if(!disableMove){
			$(window).scrollTop(Math.max(cell.offset().top - 100, 0));
		}
	} else {
		if(cell.length == 0){
			$("#jumpToNow").attr('disabled', true)
			$("#jumpToNow").innerHTML("No show on air right now!");
			setTimeout(function(){
				$("#jumpToNow").attr('disabled', false);
				$("#jumpToNow").innerHTML("Jump to current show");
			}, 3000);

		} else {
			console.error("Found multiple cells matching current time");
		}
	}
}

$(function(){
	jumpToNow(true);
})