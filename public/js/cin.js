/**
 * University Radio York Candidate Interview Night
 */

/*

   TODO:
   * Interviewer Names (see if we keep this idea)

*/

import { h, render } from 'https://unpkg.com/preact@latest?module';
import { useState, useEffect, useRef } from 'https://unpkg.com/preact@latest/hooks/dist/hooks.module.js?module';
import htm from 'https://unpkg.com/htm?module';

// Initialize htm with Preact
const html = htm.bind(h);

var api = "";
var interviews = [];

var refreshTime = 100;
const longTermRefreshTime = 5000;


function LiveCard(props) {
    if (props.show) {
        return html `
        <div class="card bg-cin-card mx-auto m-2 mt-4 mb-4" style="width: 35em;">
            <div class="card-body">
                <div class="card-title"><h1 class="cin-text-2"><b>${props.live}</b></h1></div>
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
    <div class="card bg-cin-card mx-auto m-2 mt-4 mb-4" style="width: 35em";>
        <div class="card-body">
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h3>${props.candidate}</h3></div>
            <div class="card-text">with <b>${props.interviewer}</b></div>
            <div class="card-title">${props.time}</div>
        </div>
    </div>
    `;
}

const LiveScheduleCard = (props) => {
    return html `
    <div class="card bg-cin-card mx-auto m-2 mt-4 mb-4" style="width: 35em";>
        <div class="card-body">
            <div class="card-title"><h1 class="text-danger">Live</h1></div>
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h3>${props.candidate}</h3></div>
            <div class="card-text">with <b>${props.interviewer}</b></div>
            <div class="card-text">${props.time}</div>
        </div>
    </div>
    `;
}

const PastScheduleCard = (props) => {
    return html `
    <div class="card bg-cin-card mx-auto m-2 mt-4 mb-4" style="width: 35em";>
        <div class="card-body">
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h3>${props.candidate}</h3></div>
            <div class="card-text">with <b>${props.interviewer}</b></div>
            <div class="card-text">${props.youtubeStatus}</div>
        </div>
    </div>
    `;
}


function prettifyCandidates(candidates) {
    var names = [];
    candidates.forEach(candidate => {
        names.push(candidate.full_name);
    });
    return names.join(", ");
}

function escapeRegExp(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

const ScheduleArea = () => {

        const [slots, setSlots] = useState([]);
        const searchTerm = useRef("");
        const [searched, setSearched] = useState(true);

        const handleSearch = (event) => {
            searchTerm.current = event.target.value;
            console.log("Searching: " + event.target.value);
            setSearched(false);
        }

        const updateSchedule = (auto) => {
                console.log("CALLED: ", auto, searchTerm.current)
                if (!auto || searchTerm.current == "") {

                    console.log("Update Schedule");
                    var tmp = [html `<input type="search" id="search" class="form-control mx-auto" placeholder="Search" aria-label="Search" onKeyUp=${handleSearch} style="width: 25em;"/>`];
                    interviews.forEach(event => {
                                // Spaces seem to break the search, so just yeet the space characters
                                if (searchTerm.current == "" ||
                                    event.interview.position.full_name.toLowerCase().replace(/\s/g, '').search(escapeRegExp(searchTerm.current).toLowerCase().replace(/\s/g, '')) != -1 ||
                                    prettifyCandidates(event.interview.candidates).toLowerCase().replace(/\s/g, '').search(escapeRegExp(searchTerm.current).toLowerCase().replace(/\s/g, '')) != -1) {
                                    if (
                                        new Date(event.end_time).getTime() < Date.now()
                                    ) {
                                        var youtube = event.interview.youtube_id;
                                        tmp.push(html `<${PastScheduleCard}
                                position=${event.interview.position.full_name} 
                                candidate=${prettifyCandidates(event.interview.candidates)} 
                                interviewer="Interviewer Name"
                                youtubeStatus=${youtube != null ? html `<a class="cin-text-2" href="https://www.youtube.com/watch?v=${youtube}">Watch on YouTube</a>` : "Available on YouTube Soon"}
                                />`)
                            } else if (
                                new Date(event.start_time).getTime() > Date.now()
                            ) {
                                let time = new Date(event.start_time).toLocaleTimeString().slice(0, -3) + " - " + new Date(event.end_time).toLocaleTimeString().slice(0, -3)
                                tmp.push(html `<${FutureScheduleCard}
                                time=${time}
                                position=${event.interview.position.full_name} 
                                candidate=${prettifyCandidates(event.interview.candidates)} 
                                interviewer="Interviewer Name"
                                />`)
                            } else {
                                let time = "Now - " + new Date(event.end_time).toLocaleTimeString().slice(0, -3)
                                tmp.push(html `<${LiveScheduleCard} 
                                position=${event.interview.position.full_name} 
                                candidate=${prettifyCandidates(event.interview.candidates)} 
                                interviewer="Interviewer Name"
                                time=${time}
                                />`)
                            }
                        }
        })
        if (tmp.length == 1) { // No Interviews, Only Search
            if (searchTerm.current == "") {
                setSlots([html `<h2 class="text-center">Coming Soon...</h2>`]);
            } else {
                setSlots(tmp.concat([html `<br /><h2 class="text-center">No Results</h2>`]));
            }
        } else {
            setSlots(tmp);
        }
    }
}


    useEffect(() => {
        if (!searched) {
            updateSchedule(false);
        } else if (searchTerm.current == "") { // Saves generating loads of refreshes by searching
            setTimeout(() => { updateSchedule(true) }, refreshTime);
        }
        setSearched(true);
    })

    return html `
    <div>
        <h1 class="display-3 cin-text text-center">All Interviews</h1>
        
        ${slots}
    </div>
    `
}

const LiveArea = () => {

    const [keepAlive, setKeepAlive] = useState(false);
    const [positions, setPositions] = useState(["Live Position", "Next Position"])
    const [candidates, setCandidates] = useState(["Live Candidate", "Next Candidate"])
    const [interviewers, setInterviewers] = useState(["Live Interviewer", "Next Interviewer"])
    const [times, setTimes] = useState(["Live Time", "Next Time"])
    const [showLive, setShowLive] = useState(false);
    const [showNext, setShowNext] = useState(false);

    const updateLives = async() => {
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
                        interviews[i + 1].interview.position.full_name
                    ])

                    setCandidates([
                        prettifyCandidates(interviews[i].interview.candidates),
                        prettifyCandidates(interviews[i + 1].interview.candidates)
                    ])

                    setInterviewers(["some interviewer", "some other interviewer"]);

                    setTimes([
                        "Now - " + new Date(interviews[i].end_time).toLocaleTimeString().slice(0, -3),
                        new Date(interviews[i + 1].start_time).toLocaleTimeString().slice(0, -3) + " - " + new Date(interviews[i + 1].end_time).toLocaleTimeString().slice(0, -3),
                    ])
                    break;

                } else {
                    setShowNext(false);
                    setPositions([interviews[i].interview.position.full_name, ""])
                    setCandidates([prettifyCandidates(interviews[i].interview.candidates), ""])
                    setInterviewers(["some interviewer", "some other interviewer"]);
                    setTimes(["Now - " + new Date(interviews[i].end_time).toLocaleTimeString().slice(0, -3), ""])
                    break;
                }
            } else {
                setShowLive(false);
                var intDate = new Date(interviews[i].start_time);
                if (intDate.getTime() > Date.now() && intDate.getTime() < Date.now() + 900000) {
                    setShowNext(true);
                    setPositions(["", interviews[i].interview.position.full_name]);
                    setCandidates(["", prettifyCandidates(interviews[i].interview.candidates)]);
                    setInterviewers(["", "Interviewr"]);
                    setTimes(["", intDate.toLocaleTimeString().slice(0, -3) + " - " + new Date(interviews[i].end_time).toLocaleTimeString().slice(0, -3)]);
                    break;
                } else {
                    setShowNext(false);
                }
            }
        }

    }

    useEffect(() => {
        setTimeout(() => { updateLives() }, refreshTime);
    })

    return html `
    <div class="row">
    <${LiveCard} class="col" live="Live Now" 
        show=${showLive}
        position=${positions[0]} 
        candidate=${candidates[0]} 
        interviewer=${interviewers[0]} 
        time=${times[0]}
        />

    <${LiveCard} class="col" live="Next Up" 
        show=${showNext}
        position=${positions[1]} 
        candidate=${candidates[1]} 
        interviewer=${interviewers[1]} 
        time=${times[1]}
        />
    </div>
    `
}

const quickSortTime = (data) => {
    if (data.length === 0) {
        return [];
    } else {
        var before = [];
        var after = [];
        var pivot = data[0];
        for (let i = 1; i < data.length; i++) {
            if (new Date(data[i].start_time).getTime() < new Date(pivot.start_time).getTime()) {
                before.push(data[i])
            } else {
                after.push(data[i])
            }
        }
        return quickSortTime(before).concat([pivot]).concat(quickSortTime(after))
    }
}

const getData = async() => {
    console.log("Updating API Data");
    fetch(api + "/events/")
        .then(r => r.json())
        .then(data => {
            interviews = quickSortTime(data.data);
        })

}

const App = () => {

    useEffect(() => {
        setTimeout(() => { getData() }, refreshTime);
    })

    return html `
    <${LiveArea}/>
    <${ScheduleArea}/>
    `
}

const getApi = async() => {
    const apiCall = await fetch("/cinapi");
    const res = await apiCall.text();
    return res;
}

getApi().then((x) => {
    api = x;
    getData().then(() => {
        render(html `<${App} />`, interactive)

        // This stuff is fun. It (mostly) works though.
        setInterval(() => { getData() }, longTermRefreshTime);
        setTimeout(() => refreshTime = longTermRefreshTime, refreshTime * 5);
        
    })
})