package api

import "time"

type APIConfig struct {
	TodoPassword string
	TokenExpiry  time.Duration
}
