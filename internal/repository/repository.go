package repository

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"go.uber.org/zap"
	"os"
	"sync/atomic"
)

type Repository interface {
	ReadSite(u *UrlRepository) (string, error)
}

type UrlRepository struct {
	RepoSiteInfo       map[string]model.SiteResponseInfo
	RepoSiteName       map[string]int
	RepoSiteCount      []atomic.Uint64
	RepoSiteMinMaxInfo *model.SiteMinMaxInfo
	RepoSiteMinMaxStat *model.SiteMinMaxStat
	log                *zap.Logger
}

func NewRepository(log *zap.Logger, file string) *UrlRepository {
	log.Debug("Register repository...")
	url := initData(log, file)
	return url
}

func initData(log *zap.Logger, file string) *UrlRepository {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}
	defer f.Close()

	sites := map[string]model.SiteResponseInfo{}
	sitesName := map[string]int{}

	log.Debug("Read sites list..")
	scanner := bufio.NewScanner(f)
	i := 0

	for scanner.Scan() {
		site := scanner.Text()
		sites[site] = model.SiteResponseInfo{}
		sitesName[site] = i
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}

	minMaxInfo := &model.SiteMinMaxInfo{}
	minMaxStat := &model.SiteMinMaxStat{}
	siteCount := make([]atomic.Uint64, len(sites))

	return &UrlRepository{sites, sitesName, siteCount, minMaxInfo, minMaxStat, log}
}

func (u *UrlRepository) ReadSiteInfo(s string) (model.SiteResponseInfo, error) {
	u.log.Debug(fmt.Sprintf("Read site info %s from repository", s))
	output, ok := u.RepoSiteInfo[s]
	if !ok {
		return model.SiteResponseInfo{}, errors.New("incorrect site requested")
	}

	time := u.RepoSiteInfo[s].ResponseTime
	if time == 0.0 {
		if u.RepoSiteInfo[s].Code == 0 {
			return model.SiteResponseInfo{}, errors.New("site data not loaded. Please wait")
		}
		return model.SiteResponseInfo{}, errors.New("site is unavailable")
	}
	return output, nil
}

func (u *UrlRepository) ReadMinResponseSite() (model.SiteResponseInfo, error) {
	u.log.Debug("Get Min Response Site From Repository")
	key := u.RepoSiteMinMaxInfo.MinName
	result, err := u.ReadSiteInfo(key)
	return result, err
}

func (u *UrlRepository) ReadMaxResponseSite() (model.SiteResponseInfo, error) {
	u.log.Debug("Get Max Response Site From Repository")
	key := u.RepoSiteMinMaxInfo.MaxName
	result, err := u.ReadSiteInfo(key)
	return result, err
}

func (u *UrlRepository) GetCountSiteRequest(s string) (uint64, error) {
	u.log.Debug(fmt.Sprintf("Read site count requests %s from repository", s))
	key, ok := u.RepoSiteName[s]
	if !ok {
		return 0, errors.New("incorrect site requested")
	}
	return u.RepoSiteCount[key].Load(), nil
}

func (u *UrlRepository) GetCountMaxRequest() uint64 {
	u.log.Debug("Get Max count request to Repository")
	return u.RepoSiteMinMaxStat.MaxCount.Load()
}

func (u *UrlRepository) GetCountMinRequest() uint64 {
	u.log.Debug("Get Max count request to Repository")
	return u.RepoSiteMinMaxStat.MinCount.Load()
}

func (u *UrlRepository) WriteCountSiteRequest(s string) {
	u.log.Debug(fmt.Sprintf("Write count site %s request to Repository", s))
	key := u.RepoSiteName[s]
	u.RepoSiteCount[key].Add(1)
}

func (u *UrlRepository) WriteCountMaxRequest() {
	u.log.Debug("Write Max count request to Repository")
	u.RepoSiteMinMaxStat.MaxCount.Add(1)
}

func (u *UrlRepository) WriteCountMinRequest() {
	u.log.Debug("Write Min count request to Repository")
	u.RepoSiteMinMaxStat.MinCount.Add(1)
}
