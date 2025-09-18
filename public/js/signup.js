import { g, h, s, t } from './xeact.js';

(function () {
  const MYRADIO_API_BASE = 'https://ury.org.uk/api/v2';

  const profileNextButton = g('profile-continue');

  function colourForAvailability(current, max) {
    const diff = max - current;
    switch (diff) {
      case 0:
        return 'secondary';
      case 1:
        return 'warning';
      default:
        return 'success';
    }
  }

  async function updateAvailableSlots() {
    const resp = await fetch(`${MYRADIO_API_BASE}/demo/listdemosforsignup/?api_key=${MyRadioAPIKey}`);
    const data = await resp.json();
    const trainingsList = g('available-trainings');

    document.querySelectorAll('.is-actual-training-slot').forEach(ele => ele.remove());

    trainingsList.prepend(...data.payload.filter(demo => demo.presenterstatusid === 'Studio Trained').map(demo => {
      const remainingTime = demo.demo_time_ - (Date.now() / 1000);
      const unavailable = (demo.attendee_count === demo.max_participants) || (remainingTime < (3600 * demo.signup_cutoff_hours));
      if (unavailable) {
        return null;
      }
      const col = colourForAvailability(demo.attendee_count, demo.max_participants);
      return h('label', {
        className: 'card h4 is-actual-training-slot',
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
              disabled: unavailable,
            }),
            h('div', {
              className: 'form-check-label',
            }, [
              h('span', {
                className: `badge rounded-pill text-bg-${col}`,
              }, [
                t(`${demo.attendee_count} / ${demo.max_participants}`),
              ]),
              t(' ' + demo.demo_time + ' (with ' + demo.member + ')'),
            ]),
          ]),
        ])
      ]);
    }).filter(r => r !== null));

    g('signup').disabled = true;
    initTrainingSessionEvents();
  }

  function checkPersonalDetails() {
    if (g('first-name').value === '') {
      return 'You need to enter your first name';
    }
    if (g('last-name').value === '') {
      return 'You need to enter your last name';
    }
    let email = g('email').value;
    if (email === '') {
      return 'You need to enter your email address';
    }
    if (email.endsWith('@york.ac.uk')) {
      email = email.replace(/@york\.ac\.uk$/, '');
    }
    if (!email.match(/^([a-z]|[A-Z]){1,6}[0-9]{1,6}$/)) {
      return 'Your email address looks incorrect';
    }
  }

  function onTrainingSessionSelected() {
    console.log('onTrainingSessionSelected');
    g('signup').disabled = false;
  }

  function initTrainingSessionEvents() {
    s('.training-slot-select').forEach(ele => ele.addEventListener('change', onTrainingSessionSelected));
  }

  function enableProfileContinue() {
    profileNextButton.classList.remove('disabled');
    profileNextButton.ariaDisabled = false;
    if (profileNextButton.href === "") {
      profileNextButton.href = profileNextButton.dataset.href;
      profileNextButton.dataset.href = undefined;
    }
  }

  function disableProfileContinue() {
    if (profileNextButton.ariaDisabled === 'true') return;
    profileNextButton.classList.add('disabled');
    profileNextButton.ariaDisabled = true;
    profileNextButton.dataset.href = profileNextButton.href;
    profileNextButton.removeAttribute('href');
  }

  g('signup').disabled = true;
  disableProfileContinue();
  profileNextButton.addEventListener('click', function () {
    if (profileNextButton.ariaDisabled === 'true') {
      alert(checkPersonalDetails());
    }
  });
  initTrainingSessionEvents();

  s('#signup-form input').forEach(input => {
    input.addEventListener('keydown', function (event) {
      if (event.key === 'Enter') {
        event.preventDefault();
      }
    });
    input.addEventListener('change', function () {
      if (!checkPersonalDetails()) {
        enableProfileContinue();
      } else {
        disableProfileContinue();
      }
    });
  });

  g('profile-continue').addEventListener('click', updateAvailableSlots);
})();
