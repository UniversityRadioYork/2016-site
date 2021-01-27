/**
 * University Radio York Candidate Interview Night
 */

import { h, render } from 'https://unpkg.com/preact@latest?module';
import { useState, useEffect } from 'https://unpkg.com/preact@latest/hooks/dist/hooks.module.js?module';
import htm from 'https://unpkg.com/htm?module';

// Initialize htm with Preact
const html = htm.bind(h);

var api = "";

const getApiEndpoint = async() => {
    const apiCall = await fetch("/cinapi");
    const res = await apiCall.text();
    return res;
}

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

    const [loaded, setLoaded] = useState(false);

    const [positions, setPositions] = useState(["Live Position", "Next Position"])
    const [candidates, setCandidates] = useState(["Live Candidate", "Next Candidate"])
    const [interviewers, setInterviewers] = useState(["Live Interviewer", "Next Interviewer"])

    useEffect(() => {
        const updateLives = async() => {
            fetch(api + "/events/live")
                .then(res => res.json())
                .then((data) => {
                    console.log(data);
                    setPositions([data.current.interview.position.full_name, data.next.interview.position.full_name]);
                    setCandidates([prettifyCandidates(data.current.interview.candidates), prettifyCandidates(data.next.interview.candidates)]);
                    setInterviewers(["some interviewer", "some other interviewer"]);
                })
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

function App(props) {

    return html `
    <${LiveArea}/>
    `
}

getApiEndpoint().then((x) => {
    api = x;
    console.log(api);
    render(html `<${App} />`, interactive)
})