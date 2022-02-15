/**
 * University Radio York Candidate Interview Night
 */

import { h, render } from "https://unpkg.com/preact@latest?module";
import {
    useState,
    useEffect,
    useRef,
} from "https://unpkg.com/preact@latest/hooks/dist/hooks.module.js?module";
import htm from "https://unpkg.com/htm?module";

// Initialize htm with Preact
const html = htm.bind(h);

let interviews = [];

let refreshTime = 3000;
const longTermRefreshTime = 20000;
const defaultYouTube = "dtRgUJHNHII";

function LiveCard(props) {
    if (props.show) {
        return html `
      <div
        class="card bg-cin-card mx-auto m-2 mt-4 mb-4"
        style="width: 35em; max-width: 90%;"
      >
        <div class="card-body">
          <div class="card-title">
            <h1 class="cin-text-2"><b>${props.live}</b></h1>
          </div>
          <div class="card-text"><h2>${props.position}</h2></div>
          <div class="card-text"><h3>${props.candidate}</h3></div>
          <div class="card-text">with <b>${props.interviewer}</b></div>
          <div class="card-text">${props.time}</div>
        </div>
      </div>
    `;
    }
}

const FutureScheduleCard = (props) => {
    return html `
    <div
      class="card bg-cin-card mx-auto m-2 mt-4 mb-4"
      style="width: 35em; max-width: 90%;"
    >
      <div class="card-body">
        <div class="card-text"><h2>${props.position}</h2></div>
        <div class="card-text"><h3>${props.candidate}</h3></div>
        <div class="card-text">with <b>${props.interviewer}</b></div>
        <div class="card-title">${props.time}</div>
      </div>
    </div>
  `;
};

const LiveScheduleCard = (props) => {
    return html `
    <div
      class="card bg-cin-card mx-auto m-2 mt-4 mb-4"
      style="width: 35em;  max-width: 90%;"
    >
      <div class="card-body">
        <div class="card-title">
          <h1><a href="#liveStream" class="text-danger">Live</a></h1>
        </div>
        <div class="card-text"><h2>${props.position}</h2></div>
        <div class="row">
          <div class="col">
            <div class="card-text"><h3>${props.candidate}</h3></div>
            <div class="card-text">with <b>${props.interviewer}</b></div>
            <div class="card-text">${props.time}</div>
          </div>
          <div class="col">
            <a
              href="#liveStream"
              class="ml-5 pl-5 fa fa-play-circle youtubePlay"
              style="font-size:5em;color:white"
            ></a>
          </div>
        </div>
      </div>
    </div>
  `;
};

const PastScheduleCard = (props) => {
    let playButton = "";
    if (props.youtubeID != null) {
        playButton = html `<a
      href="javascript:void(0)"
      onClick=${() => props.callback(props.youtubeID)}
      class="ml-5 pl-5 fa fa-play-circle youtubePlay"
      style="font-size:5em;color:white"
    ></a>`;
    }

    let youtubeStatus = "";
    if (props.youtubeStatus == "WATCH") {
        youtubeStatus = html `<a class="cin-text-2" href="javascript:void(0)"
        >Watch on YouTube</a
      >`;
    } else if (props.youtubeStatus == "AVAILABLE SOON") {
        youtubeStatus = "Available on YouTube Soon";
    }

    return html `
    <div
      class="card bg-cin-card mx-auto m-2 mt-4 mb-4"
      style="width: 35em;  max-width: 90%;"
    >
      <div class="card-body">
        <div class="card-text"><h2>${props.position}</h2></div>
        <div class="row">
          <div class="col">
            <div class="card-text"><h3>${props.candidate}</h3></div>
            <div class="card-text">with <b>${props.interviewer}</b></div>
            <div
              class="card-text"
              onClick=${() => props.callback(props.youtubeID)}
            >
              ${youtubeStatus}
            </div>
          </div>
          <div class="col">${playButton}</div>
        </div>
      </div>
    </div>
  `;
};

function prettifyCandidates(candidates) {
    let names = [];
    candidates.forEach((candidate) => {
        names.push(candidate.full_name);
    });
    return names.join(", ");
}

const getInterviewers = (event) => {
    let names = [];
    event.user_roles.forEach((user) => {
        if (user.role.name == "Interviewer") {
            names.push(user.user.name + (user.user.postnom.length > 0 ? (" (" + user.user.postnom + ")") : ""));
        }
    });
    return names.join(", ");
};

const getLatestODVideo = () => {
    // return "7lUz0xU5d9g";
    for (let i = interviews.length - 1; i >= 0; i--) {
        if (interviews[i].interview.youtube_id != null) {
            return interviews[i].interview.youtube_id;
        }
    }
    return defaultYouTube; // URY Ad (there's nothing else to show :()
};

function escapeRegExp(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

const ScheduleArea = () => {
        const [slots, setSlots] = useState([{ "type": "LOADING" }]);
        const searchTerm = useRef("");
        const [searched, setSearched] = useState(true);
        const [youtubeVid, setYoutubeVid] = useState(getLatestODVideo());
        const [youtubeTitle, setYoutubeTitle] = useState(getLatestODVideo() == defaultYouTube ? "On Demand Videos Available Soon!" : "You Just Missed");

        const handleSearch = (event) => {
            searchTerm.current = event.target.value;
            setSearched(false);
        };

        const updateYoutube = (id) => {
            setYoutubeVid(id);

            if (id == defaultYouTube) {
                setYoutubeTitle("On Demand Videos Available Soon!");
            } else {
                setYoutubeTitle("");
            }
        };

        const updateSchedule = (auto) => {
            if (!auto || searchTerm.current == "") {
                console.log("Update Schedule");
                let scheduleData = [{
                    "type": "SEARCH"
                }];
                interviews.forEach((event) => {
                    // Spaces seem to break the search, so just yeet the space characters
                    if (
                        searchTerm.current == "" ||
                        event.interview.position.full_name
                        .toLowerCase()
                        .replace(/\s/g, "")
                        .search(
                            escapeRegExp(searchTerm.current).toLowerCase().replace(/\s/g, "")
                        ) != -1 ||
                        prettifyCandidates(event.interview.candidates)
                        .toLowerCase()
                        .replace(/\s/g, "")
                        .search(
                            escapeRegExp(searchTerm.current).toLowerCase().replace(/\s/g, "")
                        ) != -1
                    ) {
                        if (new Date(event.end_time).getTime() < Date.now()) {
                            let youtube = event.interview.youtube_id;
                            if (typeof youtube === "string" && youtube === "") {
                              youtube = null;
                            }

                            scheduleData.push({
                                "type": "PAST",
                                "position": event.interview.position.full_name,
                                "candidate": prettifyCandidates(event.interview.candidates),
                                "interviewer": getInterviewers(event),
                                "youtubeStatus": youtube != null ? "WATCH" : "AVAILABLE SOON",
                                "callback": updateYoutube,
                                "youtubeID": youtube
                            });

                        } else if (new Date(event.start_time).getTime() > Date.now()) {
                            let startTime = new Date(event.start_time);
                            let endTime = new Date(event.end_time);
                            let time = startTime.getHours().toString().padStart(2, "0") + ":" + startTime.getMinutes().toString().padStart(2, "0") +
                                " - " + endTime.getHours().toString().padStart(2, "0") + ":" + endTime.getMinutes().toString().padStart(2, "0");

                            scheduleData.push({
                                "type": "FUTURE",
                                "position": event.interview.position.full_name,
                                "candidate": prettifyCandidates(event.interview.candidates),
                                "interviewer": getInterviewers(event),
                                "time": time
                            });

                        } else {
                            let endTime = new Date(event.end_time);
                            let time =
                                "Now - " +
                                endTime.getHours().toString().padStart(2, "0") + ":" + endTime.getMinutes().toString().padStart(2, "0");

                            scheduleData.push({
                                "type": "LIVE",
                                "position": event.interview.position.full_name,
                                "candidate": prettifyCandidates(event.interview.candidates),
                                "interviewer": getInterviewers(event),
                                "time": time
                            });

                        }
                    }
                });
                if (scheduleData.length == 1) {
                    // No Interviews, Only Search
                    if (searchTerm.current == "") {
                        setSlots([{ "type": "COMING SOON" }]);
                    } else {
                        setSlots(scheduleData.concat({ "type": "NO RESULTS" }));
                    }
                } else {
                    setSlots(scheduleData);

                    if (youtubeVid == defaultYouTube) {
                        setYoutubeVid(getLatestODVideo())
                        setYoutubeTitle("You Just Missed")
                    }
                }
            }
        };

        useEffect(() => {
            if (!searched) {
                updateSchedule(false);
            } else if (searchTerm.current == "") {
                // Saves generating loads of refreshes by searching
                setTimeout(() => {
                    updateSchedule(true);
                }, refreshTime);
            }
            setSearched(true);
        });

        let youtubeColumn = "";
        if (isCINlive) {
            youtubeColumn = html `
      <div class="col">
        <div
          style="display: flex; position: -webkit-sticky;position: sticky;top: 33vh;"
        >
          <div>
            <h3 class="text-center">${youtubeTitle}</h3>
            <iframe
              src="https://www.youtube.com/embed/${youtubeVid}"
              width="600"
              height="338"
              style="border:none;overflow:hidden; max-width: 90%;"
              scrolling="no"
              frameborder="0"
              allowTransparency="true"
              allow="encrypted-media"
              allowfullscreen="true"
            ></iframe>
            <br /><a
              href="https://www.youtube.com/watch?v=${youtubeVid}"
              class="cin-text-2"
              >[External Link]</a
            >
          </div>
        </div>
      </div>
    `;
        }

        return html `
    <div>
      <h1 class="display-3 cin-text text-center">All Interviews</h1>
      <div class="row">
        <div class="col">${
          slots.map(item => {
            switch(item.type){
              case "SEARCH":
                return html `<input
                type="search"
                key="search"
                id="search"
                class="form-control mx-auto bg-cin"
                placeholder="Search"
                aria-label="Search"
                onKeyUp=${handleSearch}
                style="width: 25em;  max-width: 90%;"
              />`;

              case "PAST":
                return html `<${PastScheduleCard}
                position=${item.position}
                key=${item.candidate.replace(/\s+/g, '')}
                candidate=${item.candidate}
                interviewer=${item.interviewer}
                youtubeStatus=${item.youtubeStatus}
                callback=${item.callback}
                youtubeID=${item.youtubeID}
              />`;

              case "FUTURE":
                return html `<${FutureScheduleCard}
                time=${item.time}
                key=${item.candidate.replace(/\s+/g, '')}
                position=${item.position}
                candidate=${item.candidate}
                interviewer=${item.interviewer}
              />`;

              case "LIVE":
                return html `<${LiveScheduleCard}
                position=${item.position}
                key=${item.candidate.replace(/\s+/g, '')}
                candidate=${item.candidate}
                interviewer=${item.interviewer}
                time=${item.time}
              />`;

              case "COMING SOON":
                return html `<h2 class="text-center" key="soon">Coming Soon...</h2>`;

              case "NO RESULTS":
                  return html `<div key="noresults"><br />
                  <h2 class="text-center">No Results</h2></div>`;

              case "LOADING":
                  return html `<h2 class="text-center" key="load">Loading...</h2>`
            }
          })
        }</div>
        ${youtubeColumn}
      </div>
    </div>
  `;
};

const LiveArea = () => {
  const [keepAlive, setKeepAlive] = useState(false);
  const [positions, setPositions] = useState([
    "Live Position",
    "Next Position",
  ]);
  const [candidates, setCandidates] = useState([
    "Live Candidate",
    "Next Candidate",
  ]);
  const [interviewers, setInterviewers] = useState([
    "Live Interviewer",
    "Next Interviewer",
  ]);
  const [times, setTimes] = useState(["Live Time", "Next Time"]);
  const [showLive, setShowLive] = useState(false);
  const [showNext, setShowNext] = useState(false);

  const updateLives = async () => {
    setKeepAlive(!keepAlive);
    console.log("Updating Live Tiles");
    for (let i = 0; i < interviews.length; i++) {
      if (
        new Date(interviews[i].start_time).getTime() < Date.now() &&
        new Date(interviews[i].end_time).getTime() > Date.now()
      ) {
        setShowLive(true);
        if (i + 1 != interviews.length) {
          setShowNext(true);
          setPositions([
            interviews[i].interview.position.full_name,
            interviews[i + 1].interview.position.full_name,
          ]);

          setCandidates([
            prettifyCandidates(interviews[i].interview.candidates),
            prettifyCandidates(interviews[i + 1].interview.candidates),
          ]);

          setInterviewers([
            getInterviewers(interviews[i]),
            getInterviewers(interviews[i + 1]),
          ]);

          let currentEnd = new Date(interviews[i].end_time);
          let nextStart = new Date(interviews[i+1].start_time);
          let nextEnd = new Date(interviews[i+1].end_time);

          setTimes([
            "Now - " +
              currentEnd.getHours().toString().padStart(2, "0") + ":" + currentEnd.getMinutes().toString().padStart(2, "0"),

              nextStart.getHours().toString().padStart(2, "0") + ":" + nextStart.getMinutes().toString().padStart(2, "0") + " - "
              + nextEnd.getHours().toString().padStart(2, "0") + ":" + nextEnd.getMinutes().toString().padStart(2, "0")
            ]);
          break;
        } else {
          setShowNext(false);
          setPositions([interviews[i].interview.position.full_name, ""]);
          setCandidates([
            prettifyCandidates(interviews[i].interview.candidates),
            "",
          ]);
          setInterviewers([getInterviewers(interviews[i]), ""]);

          let endTime = new Date(interviews[i].end_time);
          setTimes([
            "Now - " +
              endTime.getHours().toString().padStart(2, "0") + ":" + endTime.getMinutes().toString().padStart(2, "0"),
            "",
          ]);
          break;
        }
      } else {
        setShowLive(false);
        let intDate = new Date(interviews[i].start_time);
        if (
          intDate.getTime() > Date.now() &&
          intDate.getTime() < Date.now() + 900000
        ) {
          setShowNext(true);
          setPositions(["", interviews[i].interview.position.full_name]);
          setCandidates([
            "",
            prettifyCandidates(interviews[i].interview.candidates),
          ]);
          setInterviewers(["", getInterviewers(interviews[i])]);

          let nextEndTime = new Date(interviews[i].end_time);
          setTimes([
            "",
            intDate.getHours().toString().padStart(2, "0") + ":" + intDate.getMinutes().toString().padStart(2, "0") +
              " - " +
              nextEndTime.getHours().toString().padStart(2, "0") + ":" + nextEndTime.getMinutes().toString().padStart(2, "0")
          ]);
          break;
        } else {
          setShowNext(false);
        }
      }
    }
  };

  useEffect(() => {
    setTimeout(() => {
      updateLives();
    }, refreshTime);
  });

  return html`
    <div class="row">
      <${LiveCard}
        class="col"
        live="Live Now"
        show=${showLive}
        position=${positions[0]}
        candidate=${candidates[0]}
        interviewer=${interviewers[0]}
        time=${times[0]}
      />

      <${LiveCard}
        class="col"
        live="Next Up"
        show=${showNext}
        position=${positions[1]}
        candidate=${candidates[1]}
        interviewer=${interviewers[1]}
        time=${times[1]}
      />
    </div>
  `;
};

const quickSortTime = (data) => {
  if (data.length === 0) {
    return [];
  } else {
    let before = [];
    let after = [];
    let pivot = data[0];
    for (let i = 1; i < data.length; i++) {
      if (
        new Date(data[i].start_time).getTime() <
        new Date(pivot.start_time).getTime()
      ) {
        before.push(data[i]);
      } else {
        after.push(data[i]);
      }
    }
    return quickSortTime(before).concat([pivot]).concat(quickSortTime(after));
  }
};

const getData = async () => {
  console.log("Updating API Data");
  fetch(cinAPI + "/events/")
    .then((r) => r.json())
    .then((data) => {
      let usableEvents = [];
      data.data.forEach((x) => {
        if (x.interview != null) {
          usableEvents.push(x);
        }
      });
      interviews = quickSortTime(usableEvents);
    });
};

const App = () => {
  useEffect(() => {
    setTimeout(() => {
      getData();
    }, refreshTime);
  });

  return html`
    <${LiveArea} />
    <${ScheduleArea} />
  `;
};


getData();
render(html`<${App} />`, document.getElementById("interactive"));

// This stuff is fun. It (mostly) works though.
setInterval(() => {
  getData();
}, longTermRefreshTime);

setTimeout(() => {refreshTime = longTermRefreshTime}, refreshTime * 20);