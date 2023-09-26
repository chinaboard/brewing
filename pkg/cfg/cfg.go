package cfg

import (
	"os"
	"strconv"
)

var (
	OpenAiBaseURL         = getEnvDefault("OPENAI_BASE_URL", "https://api.openai.com/")
	OpenAiToken           = getEnvDefault("OPENAI_TOKEN", "")
	OpenAiProxy           = getEnvDefault("OPENAI_PROXY", "")
	WhisperEndpoint       = getEnvDefault("WHISPER_ENDPOINT", "")
	WhisperEndpointSchema = getEnvDefault("WHISPER_ENDPOINT_SCHEMA", "http")
)

var (
	HttpPort    = getEnvDefault("HTTP_PORT", "6636")
	ShareDomain = getEnvDefault("SHARE_DOMAIN", "http://localhost")
)

var (
	MongoDBConn = getEnvDefault("MONGODB_CONN", "mongodb://localhost:27017")
)

var (
	BarkNotifyDomain = os.Getenv("BARK_NOTIFY_DOMAIN")
	BarkNotifyToken  = os.Getenv("BARK_NOTIFY_TOKEN")
)

func getEnvDefault(key string, value string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return value
}

func getEnvDefaultInt(key string, value int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}
	return value
}
