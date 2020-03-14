function displayScreenSaver() {
 $('#april-fool').fadeIn();
}

function hideScreenSaver() {
  $('#april-fool').fadeOut();;
}

const timeout = 30000;

let timeoutHandle = window.setTimeout(displayScreenSaver, timeout);


document.addEventListener('mousemove', () => {
  window.clearTimeout(timeoutHandle);
  timeoutHandle = window.setTimeout(displayScreenSaver, timeout);
  hideScreenSaver();
})
