package protocol

const (
	ServerHTTPPort = 8080
)

const (
	HTTPInfoUrl  = "/info"
	HTTPLevelUrl = "/level"
)

type ZoneInfoJson struct {
	Name             string
	DefaultLevel     string
	DefaultLevelHash string
}
