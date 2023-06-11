package cronjob

import (
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"
)

type Updater interface {
	UpdateData(site map[string]model.SiteResponseInfo, minmax model.SiteMinMaxInfo)
	GetSiteNames() map[string]int
}

func SiteCheckResponse(db Updater, timeout, refresh int) {
	logger.Debug("Starting CheckResponse...")

	for {
		var wg sync.WaitGroup
		names := db.GetSiteNames()
		ch := make(chan model.SiteResponseInfo, len(names))
		ticker := time.NewTicker(time.Duration(refresh) * time.Second)
		client := http.Client{Timeout: time.Duration(timeout) * time.Second}

		for site := range names {
			wg.Add(1)

			go func(site string) {
				defer wg.Done()
				start := time.Now()
				code := http.StatusForbidden
				resp, err := client.Head("https://" + site)
				if err == nil {
					code = resp.StatusCode
					defer resp.Body.Close()
					_, err = io.Copy(io.Discard, resp.Body)
					if err != nil {
						logger.Error("Error copying body", zap.Error(err))
					}
				}

				ch <- model.SiteResponseInfo{
					SiteName:     site,
					ResponseTime: time.Since(start).Milliseconds(),
					Code:         code,
				}
			}(site)
		}
		wg.Wait()
		close(ch)
		cacheSite := make(map[string]model.SiteResponseInfo, len(names))
		for v := range ch {
			cacheSite[v.SiteName] = v
		}
		db.UpdateData(cacheSite, checkMinMax(cacheSite))
		<-ticker.C
	}
}

func checkMinMax(cacheSite map[string]model.SiteResponseInfo) model.SiteMinMaxInfo {
	var min, max int64
	var cacheMinMaxInfo model.SiteMinMaxInfo
	for key := range cacheSite {
		if (cacheSite[key].ResponseTime < min || min == 0) && cacheSite[key].Code == http.StatusOK {
			min = cacheSite[key].ResponseTime
			cacheMinMaxInfo.MinName = key
		}
		if (cacheSite[key].ResponseTime > max) && cacheSite[key].Code == http.StatusOK {
			max = cacheSite[key].ResponseTime
			cacheMinMaxInfo.MaxName = key
		}
	}
	logger.Debug("MinMax site info updated")
	return cacheMinMaxInfo
}
