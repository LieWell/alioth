package utils

import (
	"time"

	"liewell.fun/alioth/core"
)

func StringDate(t time.Time) string {
	return t.Format(core.DefaultDateFormat)
}
