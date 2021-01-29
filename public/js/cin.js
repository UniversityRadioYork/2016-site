/**
 * University Radio York Candidate Interview Night
 */

import { h, render } from 'https://unpkg.com/preact@latest?module';
import { useState, useEffect } from 'https://unpkg.com/preact@latest/hooks/dist/hooks.module.js?module';
import htm from 'https://unpkg.com/htm?module';

// Initialize htm with Preact
const html = htm.bind(h);

var api = "";
var lastUpdate = 0;
var interviews = [];
var refreshTime = 5000;


function LiveCard(props) {
    return html `
    <div class="card mx-auto m-2" style="width: 35em;">
        <div class="card-body">
            <div class="card-title"><h1>${props.live}</h1></div>
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h2>${props.candidate}</h2></div>
            <div class="card-text">${props.interviewer}</div>
            <div class="card-text">Time</div>
        </div>
    </div>
    `;
}

const FutureScheduleCard = (props) => {
    return html `
    <div class="card mx-auto m-2" style="width: 35em";>
        <div class="card-body">
            <div class="card-title"><h1>Future</h1></div>
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h2>${props.candidate}</h2></div>
            <div class="card-text">${props.interviewer}</div>
            <div class="card-text">Time</div>
        </div>
    </div>
    `;
}

const LiveScheduleCard = (props) => {
    return html `
    <div class="card mx-auto m-2" style="width: 35em";>
        <div class="card-body">
            <div class="card-title"><h1>Live</h1></div>
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h2>${props.candidate}</h2></div>
            <div class="card-text">${props.interviewer}</div>
            <div class="card-text">Time</div>
        </div>
    </div>
    `;
}

const PastScheduleCard = (props) => {
    return html `
    <div class="card mx-auto m-2" style="width: 35em";>
        <div class="card-body">
            <div class="card-title"><h1>Past</h1></div>
            <div class="card-text"><h2>${props.position}</h2></div>
            <div class="card-text"><h2>${props.candidate}</h2></div>
            <div class="card-text">${props.interviewer}</div>
            <div class="card-text">Time</div>
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

const ScheduleArea = () => {

    const [slots, setSlots] = useState([]);


    const updateSchedule = () => {
        console.log("Update Schedule");
        var tmp = [];

        interviews.forEach(event => {
            if (
                new Date(event.end_time).getTime() < Date.now()
            ) {
                tmp.push(html `<${PastScheduleCard}
                position=${event.interview.position.full_name} 
                candidate=${prettifyCandidates(event.interview.candidates)} 
                interviewer="Interviewer Name"
                />`)
            } else if (
                new Date(event.start_time).getTime() > Date.now()
            ) {
                tmp.push(html `<${FutureScheduleCard}
                position=${event.interview.position.full_name} 
                candidate=${prettifyCandidates(event.interview.candidates)} 
                interviewer="Interviewer Name"
                />`)
            } else {
                tmp.push(html `<${LiveScheduleCard} 
                position=${event.interview.position.full_name} 
                candidate=${prettifyCandidates(event.interview.candidates)} 
                interviewer="Interviewer Name"
                />`)
            }
        })
        setSlots(tmp);
    }

    useEffect(() => {
        setTimeout(() => { updateSchedule() }, refreshTime);
    })

    return html `
    <div>
        <h1 class="display-3 cin-text text-center">Schedule</h1>
        ${slots}
    </div>
    `
}

const LiveArea = () => {

    const [positions, setPositions] = useState(["Live Position", "Next Position"])
    const [candidates, setCandidates] = useState(["Live Candidate", "Next Candidate"])
    const [interviewers, setInterviewers] = useState(["Live Interviewer", "Next Interviewer"])

    const updateLives = async() => {
        console.log("Updating Live Tiles");
        for (let i = 0; i < interviews.length; i++) {
            if (
                new Date(interviews[i].start_time).getTime() < Date.now() &&
                new Date(interviews[i].end_time).getTime() > Date.now()
            ) {
                setPositions([
                    interviews[i].interview.position.full_name,
                    interviews[i + 1].interview.position.full_name
                ])

                setCandidates([
                    prettifyCandidates(interviews[i].interview.candidates),
                    prettifyCandidates(interviews[i + 1].interview.candidates)
                ])

                setInterviewers(["some interviewer", "some other interviewer"]);
                break;
            }
        }

    }

    useEffect(() => {
        setTimeout(() => { updateLives() }, refreshTime);
    })

    return html `
    <${LiveCard} live="Live Now" position=${positions[0]} candidate=${candidates[0]} interviewer=${interviewers[0]}/>
    <${LiveCard} live="Next Up" position=${positions[1]} candidate=${candidates[1]} interviewer=${interviewers[1]}/>
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
            lastUpdate = Date.now();
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
        setInterval(() => { getData() }, refreshTime);
    })
})