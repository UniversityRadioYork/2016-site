const sendButton = document.getElementById("messagesend");
const messageBox = document.getElementById("comments");

sendButton.onclick = () => {
    if (messageBox.value != "") {
        fetch(
            "https://ury.org.uk/api/v2/timeslot/currentandnext/?api_key=" + MyRadioAPIKey
        ).then(res => res.json()).then(data => {
            var currentTimeslot = data.payload.current.id;
            fetch(
                `https://ury.org.uk/api/v2/timeslot/${currentTimeslot}/sendmessage?api_key=${MyRadioAPIKey}`, {
                    method: 'post',
                    body: messageBox.value
                }
            ).then(() => {
                messageBox.value = "";
                updateMessageboxCharacterCount();
            })

            // Update the Button
            sendButton.innerText = "Message Sent!";
            sendButton.classList.remove("btn-primary");
            sendButton.classList.add("btn-success");
            setTimeout(() => {
                sendButton.innerText = "Send Message";
                sendButton.classList.remove("btn-success");
                sendButton.classList.add("btn-primary");
            }, 3000);
        })

    }
}