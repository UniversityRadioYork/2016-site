function authenticateUser(thisEl) {
  return new Promise((resolve, reject) => {
    const state = (Math.random() * 1e6).toString(10);
    const authWindow = window.open(
      "https://accounts.spotify.com/authorize?" +
      "client_id=59c103779bf341fb80f8ac527ea05808" +
      "&response_type=token" +
      "&redirect_uri=" + encodeURIComponent("http://localhost:3000/_redirect") +
      "&state=" + state +
      "&scope=playlist-modify-public",
      "_blank",
      "width=450,height=600"
    );
    // This is stupid.
    const checkState = () => {
      let done = false;
      try {
        done = authWindow.location.host === window.location.host && authWindow.location.pathname === "/_redirect";
      } catch (e) {
        // Probably cross-origin error
        if (!(e instanceof DOMException)) {
          // It's not.
          throw e;
        }
      }
      if (done) {
        const hash = authWindow.location.hash;
        const query = authWindow.location.search;
        authWindow.close();
        if (query.length > 0 && "error" in query) {
          const params = new URLSearchParams(query);
          $(thisEl).replaceWith(`<span class="error">${params.get("error")}`);
          reject(params.get("error"));
        } else {
          const params = new URLSearchParams(hash.replace(/^#/, ""));
          resolve(params.get("access_token"));
        }
      } else {
        window.setTimeout(checkState, 100);
      }
    };
    window.setTimeout(checkState, 100);
  })
}

async function savePlaylistToSpotify(thisEl, tracklistContainerEl, showName) {
  const accessToken = await authenticateUser(thisEl);

  const api = new SpotifyWebApi();
  api.setAccessToken(accessToken);
  const statusEl = $(`<span class="status">Creating your playlist, please wait...</span>`);
  $(thisEl).replaceWith(statusEl);

  const searches = [];
  $(tracklistContainerEl).find("tr").each(function() {
    const title = $(this).find(".title").text();
    const artist = $(this).find(".artist").text();

    // Exclude Sessions
    if (title.indexOf("URY Sessions") > -1 || artist.indexOf("URY Sessions") > -1) {
      return;
    }

    searches.push(
      api.searchTracks(`track:"${title}" artist:"${artist}"`)
    )
  });

  const searchResults = await Promise.all(searches);
  console.log(searchResults);

  const tracks = searchResults
    .map(x => x.tracks.items).filter(its => its.length > 0).map(x => x[0]);

  const me = await api.getMe();

  const playlist = await api.createPlaylist(me.id, {
    name: showName
  });

  await api.addTracksToPlaylist(playlist.id, tracks.map(x => x.uri));

  let resultStr = "Playlist created!";
  const failed = searchResults.map(x => x.tracks.items).filter(its => its.length === 0).length;
  if (failed > 0) {
    resultStr += ` Sadly, we couldn't find ${failed} ${failed === 1 ? "track" : "tracks"}.`;
  }
  statusEl.text(resultStr);
}
