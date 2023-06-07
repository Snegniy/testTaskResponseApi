package model

type SiteResponseInfo struct {
	SiteName     string  `json:"site_name"`
	ResponseTime float64 `json:"response_time"`
	Code         int     `json:"code"`
}

type SiteMinMaxInfo struct {
	MinName string
	MaxName string
}

type SiteMinMaxStat struct {
	MinCount *uint64
	MaxCount *uint64
}
