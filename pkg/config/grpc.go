package config

import "time"

// List of gRCP configuration values.
const (
	BackoffDelay         = 5 * time.Second
	MaxMsgSize           = 1 << 11 // 2048
	ClientRequestTimeout = 15 * time.Second
	ServerRequestTimeout = 5 * time.Second
	RateLimit            = 150
	Burst                = 5
)
