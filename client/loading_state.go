package main

import (
	"encoding/hex"
	"os"
	"path"

	"github.com/trtstm/gosubspace/helpers"
	"github.com/trtstm/gosubspace/log"
	"github.com/trtstm/gosubspace/protocol"
)

// GameState Is a state the game can be in. E.g Loading,Playing,...
type GameState interface {
	// Run Runs the current state. Return true if this state is done and NextState can be called.
	Run() bool
	// NextState returns the next state to run.
	NextState() GameState

	// Name gets the name of the state. Used mainly for debugging.
	Name() string
}

type LoadingState struct {
	infoLoading   bool
	infoLoadingCh chan *protocol.ZoneInfoJson
	infoLoaded    bool
}

func (s *LoadingState) Run() bool {
	if !s.infoLoaded && !s.infoLoading {
		s.infoLoading = true
		s.infoLoadingCh = make(chan *protocol.ZoneInfoJson)
		log.Infof("Fetching info from server %s.", clientSettings.Server)
		go func() {
			zoneInfo := &protocol.ZoneInfoJson{}
			err := helpers.GetJSON(clientSettings.ServerHTTP+protocol.HTTPInfoUrl, zoneInfo)
			if err != nil {
				log.Errorf("Fetching info failed. Reason: %v", err)
				s.infoLoadingCh <- nil
				return
			}

			err = os.MkdirAll(path.Join(clientSettings.ZonesPath, zoneInfo.Name), 0755)
			if err != nil {
				log.Errorf("Could not create zone directory '%s'. Reason: %v", zoneInfo.Name, err)
				s.infoLoadingCh <- nil
				return
			}

			levelFilePath := path.Join(clientSettings.ZonesPath, zoneInfo.Name, zoneInfo.DefaultLevel)
			hash, err := helpers.FileHash(levelFilePath)
			if err != nil || hex.EncodeToString(hash[:]) != zoneInfo.DefaultLevelHash {
				log.Infof("Validation of '%s' failed. Redownloading...", zoneInfo.DefaultLevel)
				err = helpers.DownloadFile(
					clientSettings.ServerHTTP+protocol.HTTPLevelUrl+"/"+zoneInfo.DefaultLevel,
					levelFilePath,
				)

				if err != nil {
					log.Errorf("Downloading '%s' failed. Reason: %v", zoneInfo.DefaultLevel, err)
					s.infoLoadingCh <- nil
					return
				}
			} else {
				log.Infof("No need to download '%s'.", zoneInfo.DefaultLevel)
			}

			s.infoLoadingCh <- zoneInfo
		}()
	}

	if s.infoLoading {
		select {
		case zoneInfo := <-s.infoLoadingCh:
			if zoneInfo == nil {
				log.Error("Failed to load zone.")
				os.Exit(1)
			}

			log.Infof("Received zone: %v", zoneInfo)
			s.infoLoading = false
			s.infoLoaded = true
		default:

		}
	}

	return false
}

func (s *LoadingState) NextState() GameState {
	return nil
}

func (s *LoadingState) Name() string {
	return "LoadingState"
}
