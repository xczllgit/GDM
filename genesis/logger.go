package genesis

import (
	"log"
	"xcz/gdm/logs"
)

var (
	Logger *log.Logger
)

func init() {
	Logger = logs.GetLogger()
}
