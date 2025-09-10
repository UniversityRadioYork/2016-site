import { c, g, h, s, x } from './xeact.js';

(function () {
  const trainingContinueButton = g('training-continue');

  async function updateAvailableSlots() {
    const resp = await fetch('https://ury.org.uk/api/v2/demo/listdemos/?api_key=' + MyRadioAPIKey);
    const data = await resp.json();
    const trainingsList = g('available-trainings');

    document.querySelectorAll('.is-actual-training-slot').forEach(ele => ele.remove());

    console.log(data.payload);
    trainingsList.prepend(...data.payload.filter(demo => demo.presenterstatusid === 'Studio Trained').map(demo => 
      h('label', {
        className: 'card h4',
        for: `training-${demo.demo_id}`,
      }, [
        h('div', {
          className: 'card-body',
        }, [
          h('div', {
            className: 'form-check',
          }, [
            h('input', {
              className: 'form-check-input training-slot-select',
              type: 'radio',
              name: 'sessionid',
              value: demo.demo_id,
              id: `training-${demo.demo_id}`,
            }),
            h('div', {
              className: 'form-check-label',
              innerText: demo.demo_time,
            }),
          ]),
        ])
      ])
    ));

    disableTrainingContinue();
  }

  function onTrainingSessionSelected() {
    console.log('onTrainingSessionSelected')
    trainingContinueButton.classList.remove('disabled');
    trainingContinueButton.ariaDisabled = false;
    if (trainingContinueButton.href === "") {
      trainingContinueButton.href = trainingContinueButton.dataset.href;
      trainingContinueButton.dataset.href = undefined;
    }
  }

  function disableTrainingContinue() {
    trainingContinueButton.classList.add('disabled');
    trainingContinueButton.ariaDisabled = false;
    trainingContinueButton.dataset.href = trainingContinueButton.href;
    trainingContinueButton.removeAttribute('href');
  }

  disableTrainingContinue();

  s('.training-slot-select').forEach(ele => ele.addEventListener('change', onTrainingSessionSelected));

  g('profile-continue').addEventListener('click', updateAvailableSlots);
})();
