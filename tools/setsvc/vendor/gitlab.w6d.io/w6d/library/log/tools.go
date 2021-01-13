package log

import (
	"os"

	"github.com/google/uuid"
)

//SetLogFields take key ([]string) and values ([]string) and them on map[string]interface
func SetLogFields(key []string, values []string) map[string]interface{} {
	ret := map[string]interface{}{}
	for i := range key {
		ret[key[i]] = values[i]
	}
	return ret
}

// AddLogFields add keys/values to an existing map[string]interface{}
func AddLogFields(logField map[string]interface{}, key []string, values []string) map[string]interface{} {
	for i := range key {
		logField[key[i]] = values[i]
	}
	return logField
}

// SetLogUUID ...
func SetLogUUID() map[string]interface{} {
	return map[string]interface{}{
		"uuid": uuid.New().String(),
	}
}

// SetConfig ...
func SetConfig() {

	SetOutput(os.Getenv("OUTPUT_FORMAT"))

	switch os.Getenv("SET_LEVEL_LOG") {
	case "error":
		SetLevel(ERROR)
	case "debug":
		SetLevel(DEBUG)
	case "trace":
		SetLevel(TRACE)
	default:
		SetLevel(INFO)
	}
}
