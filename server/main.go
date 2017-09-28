package main

import (
	"flag"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/trtstm/gosubspace/log"
)

type ServerSettings struct {
	DataPath   string
	ServerPath string
	ZonePath   string
}

var serverSettings ServerSettings

func init() {
	flag.StringVar(&serverSettings.DataPath, "data", "data", "The path to the data folder.")
}

type ArenaConfig struct {
	Misc struct {
		LevelFile string
	}
}

type Arena struct {
	Config ArenaConfig
}

func NewArenaFromFile(file string) (*Arena, error) {
	arena := &Arena{
		Config: ArenaConfig{},
	}
	_, err := toml.DecodeFile(file, &arena.Config)
	if err != nil {
		return nil, err
	}

	return arena, nil
}

type ZoneConfig struct {
	Misc struct {
		DefaultArena string
	}
}

type Zone struct {
	Config ZoneConfig
	Arena  *Arena
}

func NewZoneFromFile(file string) (*Zone, error) {
	zone := &Zone{
		Config: ZoneConfig{},
	}
	_, err := toml.DecodeFile(file, &zone.Config)
	if err != nil {
		return nil, err
	}

	arena, err := NewArenaFromFile(path.Join(serverSettings.DataPath, "server", "zone", zone.Config.Misc.DefaultArena))
	if err != nil {
		return nil, err
	}

	zone.Arena = arena

	return zone, nil
}

var zone *Zone

func main() {
	flag.Parse()

	serverSettings.ServerPath = path.Join(serverSettings.DataPath, "server")
	serverSettings.ZonePath = path.Join(serverSettings.ServerPath, "zone")

	log.Infof("Starting GoSubSpace")
	log.Infof("Data: %s", serverSettings.DataPath)

	zonePath := path.Join(serverSettings.DataPath, "server", "zone", "zone.toml")
	var err error
	zone, err = NewZoneFromFile(zonePath)
	if err != nil {
		log.Errorf("Could not load '%s'. Reason: %v", zonePath, err)
	}

	log.Info("Starting http server.")
	handleHttp()

	_ = zone
}
