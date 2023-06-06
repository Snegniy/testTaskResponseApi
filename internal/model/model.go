package model

import (
	"sync/atomic"
)

type SiteResponse struct {
	Response float64
	Code     int
}

type SiteCount struct {
	Count atomic.Uint64
}

type MinMaxStat struct {
	MinCount atomic.Uint64
	MaxCount atomic.Uint64
	MinName  string
	MaxName  string
}
