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

function prettifyCandidates(candidates) {
    var names = [];
    candidates.forEach(candidate => {
        names.push(candidate.full_name);
    });
    return names.join(", ");
}

function LiveArea() {

    const [update, setUpdate] = useState(0);

    const [positions, setPositions] = useState(["Live Position", "Next Position"])
    const [candidates, setCandidates] = useState(["Live Candidate", "Next Candidate"])
    const [interviewers, setInterviewers] = useState(["Live Interviewer", "Next Interviewer"])

    useEffect(() => {
        const updateLives = async() => {

            if (update != lastUpdate) {
                for (let i = 0; i < interviews.length; i++) {
                    if (interviews[i].start_time < Date.now() && interviews[i].end_time > Date.now()) {
                        setPositions([
                            interviews[i].position.full_name,
                            interviews[i + 1].position.full_name
                        ])

                        setCandidates([
                            prettifyCandidates(interviews[i].interview.candidates),
                            prettifyCandidates(interviews[i + 1].interview.candidates)
                        ])

                        setInterviewers(["some interviewer", "some other interviewer"]);
                        break;
                    }
                }
                setUpdate(lastUpdate);
            }
        }

        if (loaded) {
            setTimeout(() => { updateLives() }, 20000);
        } else {
            updateLives();
            setLoaded(true);
        }
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
            if (data[i].start_time < pivot.start_time) {
                before.push(data[i])
            } else {
                after.push(data[i])
            }
        }
        return quickSortTime(before).push(pivot).concat(quickSortTime(after))
    }
}

const getData = async() => {
    fetch(api + "/events")
        .then(r => r.json())
        .then(data => {
            interviews = quickSortTime(data);
            lastUpdate = Date.now();
        })
}

const App = () => {

    const [loaded, setLoaded] = useState(false);

    useEffect(() => {

        if (loaded) {
            setTimeout(() => { updateLives() }, 20000);
        } else {
            getData();
            setLoaded(true);
        }
    })

    return html `
    <${LiveArea}/>
    `
}

(() => {
    const apiCall = await fetch("/cinapi");
    const res = await apiCall.text();
    return res;
})().then((x) => {
    api = x;
    console.log(api);
    render(html `<${App} />`, interactive)
})