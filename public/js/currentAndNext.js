/* global MyRadioAPIKey */
function updateShow() {

  console.log("Starting to update now and next.");

  var request = new XMLHttpRequest()
  request.open(
    'GET',
    '//ury.org.uk/api/v2/timeslot/currentandnext/?api_key=' + MyRadioAPIKey,
    true,
  )

  request.onload = function () {
    if (this.status >= 200 && this.status < 400) {
      success(JSON.parse(this.response))
    } else {
      fail(this.response)
    }
  }
  request.send()

  //Schedule when the next update will happen.
  scheduleUpdate()
}

function scheduleUpdate() {
  // Call the function 15 seconds past every half hour (to line up with live stream)
  let nextCall = new Date()

  if (nextCall.getMinutes() >= 30) {
    nextCall.setHours(nextCall.getHours() + 1)
    nextCall.setMinutes(0)
  } else {
    nextCall.setMinutes(30)
  }
  nextCall.setSeconds(15)

  let difference = nextCall - new Date()
  setTimeout(updateShow, difference)
}

function calcTime(timestamp) {
  var date = new Date(timestamp * 1000)
  // Hours part from the timestamp
  var hours = '0' + date.getHours()
  // Minutes part from the timestamp
  var minutes = '0' + date.getMinutes()
  // Use substr to remove the extra 0 if 2 digit hour/min
  return hours.substr(-2) + ':' + minutes.substr(-2)
}

function makeContent(show) {
  if (typeof show.title !== 'undefined') {
    // If we're off air
    if (show.end_time === 'The End of Time') {
      return (
        '<div class="h3 show-title ellipsis">We\'re off air right now.</div>' +
        '<div class="show-time h4">Check back next term.</div>'
      )
    }
    // Use "Now" as start time if it's missing
    let startTimeString = calcTime(show.start_time)
    if (typeof show.start_time === 'undefined') {
      startTimeString = 'Now'
    }
    // Default case (regular show)
    return (
      '<div class="h3 show-title ellipsis">' +
      show.title +
      '</div>' +
      '<div class="show-time h4">' +
      startTimeString +
      ' - ' +
      calcTime(show.end_time) +
      '</div>'
    )
  } else {
    return '<span>Looks like there is nothing on here.</span>'
  }
}

function fail(data) {
  console.log("Failed to update now and next.", data)
}

function success(data) {
  console.log("Successfully retrieved now and next.")
  // Current show
  if (typeof data.payload.current === 'undefined') {
    // There is no current show; Something is probably very wrong...
    $('.current-and-next-now').replaceWith(
      '<div class="current-and-next-now p-2 pt-3 px-3 p-md-3 p-lg-4 " title="View the show now on air.">' +
        '<h2 class="font-weight-bold">Now</h2>' +
        '<div class="h3 show-title ellipsis">There\'s nothing on right now.</div>',
    )
    $('#studiomessage *').attr('disabled', true)
  } else if (typeof data.payload.current.url !== 'undefined') {
    $('.current-and-next-now').replaceWith(
      '<a class="current-and-next-now p-2 pt-3 px-3 p-sm-3 p-lg-4 " href=' +
        data.payload.current.url +
        ' title="Show currently on air: ' +
        data.payload.current.title +
        '">' +
        '<h2 class="font-weight-bold">Now</h2>' +
        makeContent(data.payload.current) +
        '</a>',
    )
    $('#studiomessage *').attr('disabled', false)
  } else {
    $('.current-and-next-now').replaceWith(
      '<div class="current-and-next-now p-2 pt-3 px-3 p-sm-3 p-lg-4 " title="View the show now on air.">' +
        '<h2 class="font-weight-bold">Now</h2>' +
        makeContent(data.payload.current) +
        '</a>',
    )
    $('#studiomessage *').attr('disabled', true)
  }

  // Next show
  if (typeof data.payload.next === 'undefined') {
    // There is no next show (e.g. we're off air)
    $('.current-and-next-next').replaceWith(
      '<div class="current-and-next-next p-2 pt-3 px-3 p-sm-3 p-lg-4 " title="View the show up next.">' +
        '<h2 class="font-weight-bold">Next</h2>' +
        '<div class="h3 show-title ellipsis">There\'s nothing up next yet.</div>' +
        '</a>',
    )
  } else if (typeof data.payload.next.url !== 'undefined') {
    $('.current-and-next-next').replaceWith(
      '<a class="current-and-next-next p-2 pt-3 px-3 p-sm-3 p-lg-4 " href=' +
        data.payload.next.url +
        ' title="Show on air next: ' +
        data.payload.next.title +
        '.">' +
        '<h2 class="font-weight-bold">Next</h2>' +
        makeContent(data.payload.next) +
        '</a>',
    )
  } else {
    $('.current-and-next-next').replaceWith(
      '<div class="current-and-next-next p-2 pt-3 px-3 p-sm-3 p-lg-4 " title="View the show up next.">' +
        '<h2 class="font-weight-bold">Next</h2>' +
        makeContent(data.payload.next) +
        '</a>',
    )
  }

  $('.current-and-next-img img').attr(
    'src',
    '//ury.org.uk' + data.payload.current.photo,
  )

  console.log("Finished rendering now and next.")
}

$(document).ready(function () {
  // Call on startup too, this will schedule the next update.
  updateShow()
})
