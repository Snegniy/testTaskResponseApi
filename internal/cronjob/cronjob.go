package cronjob

import (
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	tick   *time.Ticker // убрать
	client http.Client  // убрать
	// список полных сайтов urls
}

func Response(db *repository.UrlRepository, timeout, refresh int) {
	c := &Cache{
		tick:   time.NewTicker(time.Duration(refresh) * time.Second),
		client: http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
	logger.Debug("Starting CheckResponse...")
	loopCheckSite(db, c)
}

func loopCheckSite(r *repository.UrlRepository, c *Cache) {
	for {
		ch := make(chan model.SiteResponseInfo, len(r.RepoSiteName))
		var wg sync.WaitGroup
		for site := range r.RepoSiteName {
			wg.Add(1)
			go func(c *Cache, site string) {
				defer wg.Done()
				start := time.Now()
				code := http.StatusForbidden
				resp, err := c.client.Head("https://" + site)
				if err == nil {
					code = resp.StatusCode
					defer resp.Body.Close()
					_, err = io.Copy(io.Discard, resp.Body)
					if err != nil {
						logger.Warn("Error copying body", zap.Error(err))
					}
				}
				ch <- model.SiteResponseInfo{
					SiteName:     site,
					ResponseTime: time.Since(start).Milliseconds(),
					Code:         code,
				}
			}(c, site)
		}
		wg.Wait()
		close(ch)
		cacheSite := make(map[string]model.SiteResponseInfo, len(r.RepoSiteInfo))
		for v := range ch {
			cacheSite[v.SiteName] = v
		}
		r.UpdateData(cacheSite, c.checkMinMax(cacheSite))
		logger.Info("data updated")
		<-c.tick.C
	}
}

func (c *Cache) checkMinMax(cacheSite map[string]model.SiteResponseInfo) model.SiteMinMaxInfo {
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
	return cacheMinMaxInfo
}
