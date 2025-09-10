(function () {
  const trainingContinueButton = document.getElementById('training-continue');

  function onTrainingSessionSelected() {
    console.log('onTrainingSessionSelected')
    trainingContinueButton.classList.remove('disabled');
    trainingContinueButton.ariaDisabled = false;
    if (trainingContinueButton.href === "") {
      trainingContinueButton.href = trainingContinueButton.dataset.href;
      trainingContinueButton.dataset.href = undefined;
    }
  }

  document.querySelectorAll('.training-slot-select').forEach(ele => ele.addEventListener('change', onTrainingSessionSelected));

  trainingContinueButton.classList.add('disabled');
  trainingContinueButton.ariaDisabled = false;
  trainingContinueButton.dataset.href = trainingContinueButton.href;
  trainingContinueButton.removeAttribute('href');
})();
