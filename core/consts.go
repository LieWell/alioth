package core

import "time"

const (
	DefaultTimeFormat    = "2006-01-02 15:04:05"
	DefaultLogTimeFormat = "2006-01-02 15:04:05.000"
	DefaultDayFormat     = "2006-01-02"
	UserNameKey          = "username"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)
