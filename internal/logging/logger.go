package logging

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func NewLogger(name string) *zap.Logger {
	rawJSON := []byte(fmt.Sprintf(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "./logs/%s"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`, name))

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())

	return logger
}
