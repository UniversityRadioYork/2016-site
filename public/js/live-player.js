export function makePlayer(config) {
    const { idPrefix, audioUrl, icecastStatusUrl } = config;
    const player = document.getElementById(`${idPrefix}-audio`);
    const playPause = document.getElementById(`${idPrefix}-play`);
    const volume = document.getElementById(`${idPrefix}-volume`);
    const currentTrackTitle = document.getElementById(`${idPrefix}-track-title`);
    const currentTrackArtist = document.getElementById(`${idPrefix}-track-artist`);
    const currentTrackContainer = document.getElementById(`${idPrefix}-track-container`);

    function updateButton() {
        if (player.paused) {
            playPause.innerHTML = '<i class="fa fa-play"></i>';
        } else {
            playPause.innerHTML = '<i class="fa fa-pause"></i>';
        }
    }

    let playbackError = false;

    function markLoading() {
        playPause.innerHTML = '<span class="player-load-dots" title="loading"><span>&bull;</span><span>&bull;</span><span>&bull;</span></span>';
    }

    function markError() {
        playbackError = true;
        playPause.disabled = true;
        playPause.innerHTML = '<i class="fa fa-exclamation-triangle"></i>';
    }

    function setNowPlaying(title, artist) {
        currentTrackTitle.innerText = title;
        currentTrackArtist.innerText = artist;
        if (title === 'URY' && artist === null) {
            currentTrackContainer.style.display = 'none';
        } else {
            currentTrackContainer.style.display = 'block';
        }
    }

    let nowPlayingUpdate = null;

    function fetchNowPlaying() {
        fetch(icecastStatusUrl)
            .then((resp) => {
                if (!resp.ok) {
                    console.error('failed to fetch current track', resp.status, resp.statusText);
                    nowPlayingUpdate = setTimeout(fetchNowPlaying, 10_000);
                    return;
                } else {
                    return resp.json();
                }
            })
            .then((resp) => {
                const stream = resp.icestats.source.filter((s) => s.listenurl.indexOf('/live-high-ogg') !== -1)[0];
                const { artist, title } = stream;

                setNowPlaying(title, artist);

                // Update every 10s
                nowPlayingUpdate = setTimeout(fetchNowPlaying, 10_000);
            }).catch((e) => {
                console.error('failed to fetch now playing', e);
            });
    }

    const playbackControls = {
        play() {
            if (playbackError) {
                console.log('playback error');
                return;
            }
            if (!player.src) {
                // Load the live feed
                player.src = audioUrl;
                player.autoplay = false;
            }
            if (!nowPlayingUpdate) {
                fetchNowPlaying();
            }
            player.play();
        },

        pause() {
            if (playbackError) {
                console.error('playback error');
                return;
            }
            if (nowPlayingUpdate) {
                clearTimeout(nowPlayingUpdate);
            }
            player.pause();
        },

        setVolume(level) {
            player.volume = level;
        }
    };

    player.addEventListener('waiting', () => {
        if (playbackError) return;
        markLoading();
    })

    player.addEventListener('pause', () => {
        if (playbackError) return;
        updateButton();
    });

    player.addEventListener('play', () => {
        if (playbackError) return;
        updateButton();
    });

    player.addEventListener('playing', () => {
        if (playbackError) return;
        updateButton();
    });

    player.addEventListener('ended', () => {
        if (playbackError) return;
        console.log('retrying load');
        player.load();
    });

    player.addEventListener('error', (ev) => {
        console.log(ev.error);
        markError();
    });

    playPause.addEventListener('click', () => {
        if (player.paused) {
            playbackControls.play();
        } else {
            playbackControls.pause();
        }
    });

    playbackControls.setVolume(parseInt(volume.value) / 11.0);
    volume.addEventListener('input', () => {
        playbackControls.setVolume(parseInt(volume.value) / 11.0);
    });

    return playbackControls;
}
