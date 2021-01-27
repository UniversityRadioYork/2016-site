/**
 * University Radio York Candidate Interview Night
 */

import { h, Component, render } from 'https://unpkg.com/preact@latest?module';
import { useState, useEffect } from 'https://unpkg.com/preact@latest/hooks/dist/hooks.module.js?module';
import htm from 'https://unpkg.com/htm?module';

// Initialize htm with Preact
const html = htm.bind(h);

function updateData() {

}

function LiveCard(props) {
    const title = (props.live == "live") ? "Live Now" : "Next Up";

    return html `
    <div class="card mx-auto m-2" style="width: 35em;">
        <div class="card-body">
            <div class="card-title"><h1>${title}</h1></div>
            <div class="card-text"><h2>Position</h2></div>
            <div class="card-text"><h2>Name</h2></div>
            <div class="card-text">Interviewer</div>
            <div class="card-text">Time</div>
        </div>
    </div>
    `;
}

function App(props) {

    return html `
    <${LiveCard} live="live" />
    <${LiveCard} live="next" />
    `
}

render(html `<${App} />`, interactive);