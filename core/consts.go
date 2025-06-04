package core

import "time"

const (
	DefaultDateFormat     = "2006-01-02"
	DefaultTimeFormat     = "15:04:05"
	DefaultDateTimeFormat = "2006-01-02 15:04:05"
	DefaultLogTimeFormat  = "2006-01-02 15:04:05.000"
	UserNameKey           = "username"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)
