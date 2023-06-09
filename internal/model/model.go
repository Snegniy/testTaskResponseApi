package model

import (
	"sync/atomic"
)

type SiteResponseInfo struct {
	SiteName     string
	ResponseTime int64
	Code         int
}

type SiteMinMaxInfo struct {
	MinName string
	MaxName string
}

type SiteMinMaxStat struct {
	MinCount atomic.Uint64
	MaxCount atomic.Uint64
}
