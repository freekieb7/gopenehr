package database

import (
	"time"
)

type Migration struct {
	Version   uint64
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt time.Time
}
