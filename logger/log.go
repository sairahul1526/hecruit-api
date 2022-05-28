package logger

import (
	"fmt"
	CONFIG "hecruit-backend/config"
)

// Log - log based on test value
func Log(str ...interface{}) {
	if CONFIG.Log {
		fmt.Println(str)
	}
}
