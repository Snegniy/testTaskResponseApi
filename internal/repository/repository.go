package repository

import (
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"github.com/Snegniy/testTaskResponseApi/pkg/logger"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type UrlRepository struct {
	RepoSiteInfo       map[string]model.SiteResponseInfo
	RepoSiteName       map[string]int
	RepoSiteCount      []atomic.Uint64
	RepoSiteMinMaxInfo model.SiteMinMaxInfo
	RepoSiteMinMaxStat model.SiteMinMaxStat
	mu                 sync.RWMutex
}

func NewRepository(file string) *UrlRepository {
	logger.Debug("Register repository...")
	logger.Debug("Read sites list..")
	sitesInfo, sitesName, err := initData(file)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%v", err))
	}

	return &UrlRepository{
		RepoSiteInfo:  sitesInfo,
		RepoSiteName:  sitesName,
		RepoSiteCount: make([]atomic.Uint64, len(sitesName)),
	}
}

func initData(file string) (map[string]model.SiteResponseInfo, map[string]int, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	list := strings.Split(string(b), "\r\n") // win - \r\n ; unix - \n ; mac - \r
	mInfo := make(map[string]model.SiteResponseInfo, len(list))
	mName := make(map[string]int, len(list))

	for _, v := range list {
		if v == "" {
			continue
		}
		mInfo[v] = model.SiteResponseInfo{SiteName: v}
		mName[v] = len(mName)
	}
	return mInfo, mName, nil
}

func (u *UrlRepository) ReadSiteInfo(s string) model.SiteResponseInfo {
	logger.Debug(fmt.Sprintf("Read site info %s from repository", s))
	u.mu.RLock()
	result, ok := u.RepoSiteInfo[s]
	u.mu.RUnlock()
	if !ok {
		return model.SiteResponseInfo{
			SiteName: s,
			Code:     -1,
		}
	}
	u.WriteCountSiteRequest(s)
	return result
}

func (u *UrlRepository) ReadMinResponseSite() model.SiteResponseInfo {
	logger.Debug("Read Min response site from repository")
	key := u.RepoSiteMinMaxInfo.MinName
	result := u.ReadSiteInfo(key)
	u.WriteCountMinRequest()
	return result
}

func (u *UrlRepository) ReadMaxResponseSite() model.SiteResponseInfo {
	logger.Debug("Read Max response site from repository")
	key := u.RepoSiteMinMaxInfo.MaxName
	result := u.ReadSiteInfo(key)
	u.WriteCountMaxRequest()
	return result
}

func (u *UrlRepository) ReadCountSiteRequest(s string) (uint64, error) {
	logger.Debug(fmt.Sprintf("Read site count requests %s from repository", s))
	key, ok := u.RepoSiteName[s]
	if !ok {
		return 0, errors.New("incorrect site requested")
	}
	return u.RepoSiteCount[key].Load(), nil
}

func (u *UrlRepository) ReadCountMinRequest() uint64 {
	logger.Debug("Read Max count request to repository")
	return u.RepoSiteMinMaxStat.MinCount.Load()
}

func (u *UrlRepository) ReadCountMaxRequest() uint64 {
	logger.Debug("Read Max count request to repository")
	return u.RepoSiteMinMaxStat.MaxCount.Load()
}

func (u *UrlRepository) WriteCountSiteRequest(s string) {
	logger.Debug(fmt.Sprintf("Write count site %s request to repository", s))
	if key, ok := u.RepoSiteName[s]; ok {
		u.RepoSiteCount[key].Add(1)
	}
}

func (u *UrlRepository) WriteCountMinRequest() {
	logger.Debug("Write Min count request to repository")
	u.RepoSiteMinMaxStat.MinCount.Add(1)
}

func (u *UrlRepository) WriteCountMaxRequest() {
	logger.Debug("Write Max count request to repository")
	u.RepoSiteMinMaxStat.MaxCount.Add(1)
}

func (u *UrlRepository) UpdateData(siteinfo map[string]model.SiteResponseInfo, minmax model.SiteMinMaxInfo) {
	u.mu.Lock()
	u.RepoSiteInfo = siteinfo
	u.RepoSiteMinMaxInfo = minmax
	u.mu.Unlock()
}
