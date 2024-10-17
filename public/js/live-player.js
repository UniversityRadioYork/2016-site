export function makePlayer(config) {
    const { idPrefix, audioUrl, icecastStatusUrl } = config;
    let player = new Audio();
    player.preload = 'none';
    const playPause = document.getElementById(`${idPrefix}-play`);
    const volume = document.getElementById(`${idPrefix}-volume`);
    const currentTrackTitle = document.getElementById(`${idPrefix}-track-title`);
    const currentTrackArtist = document.getElementById(`${idPrefix}-track-artist`);
    const currentTrackArtistContainer = document.getElementById(`${idPrefix}-track-artist-container`);
    const currentTrackContainer = document.getElementById(`${idPrefix}-track-container`);

    function updateButton() {
        if (player.paused) {
            playPause.innerHTML = '<i class="fa fa-play"></i>';
        } else {
            playPause.innerHTML = '<i class="fa fa-stop"></i>';
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
        currentTrackArtistContainer.style.display = 'inline';
        if (artist === null) {
            currentTrackTitle.innerText = "University Radio York"
            currentTrackArtistContainer.style.display = 'none';
        }
    }

    let nowPlayingUpdate = null;

    function fetchNowPlaying() {
        fetch(icecastStatusUrl)
            .then((resp) => {
                if (!resp.ok) {
                    console.error('failed to fetch current track, has the icecastStatus been added to config?', resp.status, resp.statusText);
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
                console.error('failed to fetch now playing, has the icecastStatus been added to config?', e);
            });
    }

    const playbackControls = {
        play() {
            if (playbackError) {
                console.log('playback error, has the audioUrl been set in config?');
                return;
            }
            if (!nowPlayingUpdate) {
                fetchNowPlaying();
            }

            if (!this.playing) {
                player.src = audioUrl;
                player.play();
            }
        },

        pause() {
            if (playbackError) {
                console.error('playback error');
                return;
            }
            if (nowPlayingUpdate) {
                clearTimeout(nowPlayingUpdate);
                nowPlayingUpdate = null;
            }

            player.src = null;
            player.srcObject = null;

            updateButton();
        },

        setVolume(level) {
            player.volume = level;
        },

        get playing() {
            return player.src !== null && !player.paused;
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
        console.log(ev);
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

    window.onbeforeunload = () => {
        console.log('before unload');
        if (playbackControls.playing) {
            return '';
        }
    };

    return playbackControls;
}
