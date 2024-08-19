package types

import "time"

type Clinet struct {
	Count        int
	LastAccess   time.Time
	BlockedUntil time.Time
}
