package cronjob

import (
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/internal/repository"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	cacheName map[string]string
	tick      *time.Ticker
	wg        sync.WaitGroup
	client    http.Client
}

func NewCheckResponse(timeout, refresh int, db *repository.UrlRepository) *Cache {
	name := make(map[string]string, len(db.RepoSiteName))
	for key := range db.RepoSiteName {
		name[key] = fmt.Sprintf("https://%s", key)
	}
	return &Cache{
		cacheName: name,
		tick:      time.NewTicker(time.Duration(refresh) * time.Second),
		client:    http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func Response(db *repository.UrlRepository, timeout, refresh int, log *zap.Logger) {
	c := NewCheckResponse(timeout, refresh, db)
	log.Debug("Starting CheckResponse...")
	loopCheckSite(db, c, log)
}

func loopCheckSite(r *repository.UrlRepository, c *Cache, log *zap.Logger) {
	for {
		ch := make(chan model.SiteResponseInfo)
		go SetSite(ch, c, r)

		for site := range c.cacheName {
			c.wg.Add(1)
			go func(r *repository.UrlRepository, c *Cache, site string) {
				defer c.wg.Done()
				start := time.Now()
				code := http.StatusForbidden
				resp, err := c.client.Head(c.cacheName[site])
				if err == nil {
					code = resp.StatusCode
					defer resp.Body.Close()
					_, err = io.Copy(io.Discard, resp.Body)
					if err != nil {
						log.Warn(fmt.Sprintf("Error copying body: %v", err))
					}
				}
				ch <- model.SiteResponseInfo{
					SiteName:     site,
					ResponseTime: time.Since(start).Milliseconds(),
					Code:         code,
				}
			}(r, c, site)
		}
		c.wg.Wait()
		close(ch)
		log.Info("data updated")
		<-c.tick.C
	}
}

func SetSite(ch chan model.SiteResponseInfo, c *Cache, r *repository.UrlRepository) {
	cacheSite := make(map[string]model.SiteResponseInfo, len(c.cacheName))
	for v := range ch {
		fmt.Println(v)
		cacheSite[v.SiteName] = v
	}
	fmt.Println(cacheSite)
	r.RepoSiteMinMaxInfo = c.checkMinMax(cacheSite)
	r.RepoSiteInfo = cacheSite

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
			min = cacheSite[key].ResponseTime
			cacheMinMaxInfo.MaxName = key
		}
	}
	return cacheMinMaxInfo
}
