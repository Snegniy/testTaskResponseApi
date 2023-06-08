package response

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
	cacheSite   map[string]model.SiteResponseInfo
	cacheMinMax model.SiteMinMaxInfo
	tick        *time.Ticker
	wg          sync.WaitGroup
	client      http.Client
	mu          sync.Mutex
}

func NewCheckResponse(timeout, refresh int) *Cache {
	return &Cache{
		cacheSite: map[string]model.SiteResponseInfo{},
		tick:      time.NewTicker(time.Duration(refresh) * time.Second),
		client:    http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func loopCheckSite(r *repository.UrlRepository, c *Cache, log *zap.Logger) {
	for {
		_ = <-c.tick.C
		for site := range r.RepoSiteInfo {
			c.wg.Add(1)
			site := site
			go func(r *repository.UrlRepository, c *Cache) {
				defer c.wg.Done()
				start := time.Now()
				resp, err := c.client.Get(fmt.Sprintf("https://%s", site))
				if err != nil {
					c.setSite(model.SiteResponseInfo{
						SiteName:     site,
						ResponseTime: time.Since(start).Milliseconds(),
						Code:         http.StatusForbidden,
					})
					return
				}
				c.setSite(model.SiteResponseInfo{
					SiteName:     site,
					ResponseTime: time.Since(start).Milliseconds(),
					Code:         resp.StatusCode,
				})

				defer resp.Body.Close()
				_, err = io.Copy(io.Discard, resp.Body)
				if err != nil {
					return
				}
			}(r, c)
		}
		c.wg.Wait()
		c.checkMinMax()
		c.swapData(r)
		log.Info("data updated")

	}
}

func Response(db *repository.UrlRepository, timeout, refresh int, log *zap.Logger) {
	c := NewCheckResponse(timeout, refresh)
	log.Debug("Starting CheckResponse...")
	loopCheckSite(db, c, log)
}

func (c *Cache) setSite(m model.SiteResponseInfo) {
	c.mu.Lock()
	c.cacheSite[m.SiteName] = m
	c.mu.Unlock()
}

func (c *Cache) swapData(r *repository.UrlRepository) {
	r.RepoSiteInfo = c.cacheSite
	r.RepoSiteMinMaxInfo = c.cacheMinMax
}

func (c *Cache) checkMinMax() {
	var min, max int64
	for key := range c.cacheSite {
		if (c.cacheSite[key].ResponseTime < min || min == 0) && c.cacheSite[key].Code == http.StatusOK {
			min = c.cacheSite[key].ResponseTime
			c.cacheMinMax.MinName = key
		}
		if (c.cacheSite[key].ResponseTime > max) && c.cacheSite[key].Code == http.StatusOK {
			min = c.cacheSite[key].ResponseTime
			c.cacheMinMax.MaxName = key
		}
	}
}
