package model

import "sync/atomic"

type SiteResponseInfo struct {
	SiteName     string  `json:"site_name"`
	ResponseTime float64 `json:"response_time"`
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
