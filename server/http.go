package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/trtstm/gosubspace/helpers"
	"github.com/trtstm/gosubspace/log"
	"github.com/trtstm/gosubspace/protocol"
)

func handleHttp() {
	hash, err := helpers.FileHash(path.Join(serverSettings.ZonePath, zone.Arena.Config.Misc.LevelFile))
	if err != nil {
		log.Errorf("Could not calculate hash. Reason: %v", err)
		os.Exit(1)
	}

	levelFiles := map[string]struct{}{}
	levelFiles[zone.Arena.Config.Misc.LevelFile] = struct{}{}

	zoneInfoBytes, err := json.Marshal(&protocol.ZoneInfoJson{
		Name:             "Unknown name",
		DefaultLevel:     zone.Arena.Config.Misc.LevelFile,
		DefaultLevelHash: hex.EncodeToString(hash[:]),
	})
	if err != nil {
		log.Errorf("Failed to marshal zone info. Reason: %v", err)
		os.Exit(1)
	}

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Write(zoneInfoBytes)
	})

	http.Handle("/level/", http.StripPrefix("/level", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Path[1:]
		if _, ok := levelFiles[fileName]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, path.Join(serverSettings.ZonePath, fileName))
	})))

	if err = http.ListenAndServe(":"+strconv.Itoa(protocol.ServerHTTPPort), nil); err != nil {
		log.Errorf("Could not start http server. Reason: %v", err)
		os.Exit(1)
	}
}
