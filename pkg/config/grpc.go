package config

import "time"

// TODO: move into service conf TOML and template
const (
	BackoffDelay   = 5 * time.Second
	MaxMsgSize     = 1 << 11 // 2048
	RequestTimeout = 5000 * time.Millisecond
)
