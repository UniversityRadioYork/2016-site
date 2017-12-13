function getCurrentShow() {
    $.getJSON("https://ury.org.uk/api/v2/timeslot/currentandnext/?api_key=" + myRadioAPIKey, (data) => {
        if (data.status === "OK") {
            return data.payload.current;
        } else {
            console.error("API Error: " + data.status + " - " + data.payload);
            return;
        }
    });
}

function messageSendError() {
    $("#comments").val("");
    $("#comments").attr("placeholder", "Message cannot be sent at this time. Is Jukebox on Air?");
}


$(function() {
    $("#studiomessage").submit(function(){
        var currentShow = getCurrentShow();
        var msg = $("#comments:first").val();
        if (msg !== "" && typeof currentShow != "undefined" && typeof currentShow.id == "number") {
            $("#comments").val("Sending...");
            $.ajax({
                type: "PUT",
                url: "https://ury.org.uk/api/v2/timeslot/" + currentShow.id +"/sendmessage/?api_key=" + myRadioAPIKey,
                data: { message: msg },
                error: () => {
                    messageSendError();
                },
                success: () => {
                    $("#comments").val("");
                    $("#comments").attr("placeholder","Message has been sent");
                },
                complete: () => {
                    window.setTimeout(function() {
                        $("#comments").val("");
                        $("#comments").attr("placeholder","Why not send a message to the studio presenters? Have a bit of banter, share your opinions or simply request your favourite song?");
                    }, 3000);
                }
            });
        } else {
            messageSendError();
        }
        return false;
    });
});
