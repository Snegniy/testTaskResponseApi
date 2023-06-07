package repository

import (
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"go.uber.org/zap"
	"os"
	"strings"
	"sync/atomic"
)

type UrlRepository struct {
	RepoSiteInfo       map[string]model.SiteResponseInfo
	RepoSiteName       map[string]int
	RepoSiteCount      []atomic.Uint64
	RepoSiteMinMaxInfo model.SiteMinMaxInfo
	RepoSiteMinMaxStat model.SiteMinMaxStat
	log                *zap.Logger
}

func NewRepository(log *zap.Logger, file string) *UrlRepository {
	log.Debug("Register repository...")
	log.Debug("Read sites list..")
	sitesInfo, sitesName, err := initData(file)
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}

	return &UrlRepository{
		RepoSiteInfo:  sitesInfo,
		RepoSiteName:  sitesName,
		RepoSiteCount: make([]atomic.Uint64, len(sitesName)),
		log:           log,
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
		mName[v] = len(mName)
		mInfo[v] = model.SiteResponseInfo{}
	}
	return mInfo, mName, nil
}

func (u *UrlRepository) ReadSiteInfo(s string) (model.SiteResponseInfo, error) {
	u.log.Debug(fmt.Sprintf("Read site info %s from repository", s))
	result, ok := u.RepoSiteInfo[s]
	if !ok {
		return model.SiteResponseInfo{}, errors.New("incorrect site requested")
	}
	return result, nil
}

func (u *UrlRepository) ReadMinResponseSite() (model.SiteResponseInfo, error) {
	u.log.Debug("Read Min Response Site From Repository")
	key := u.RepoSiteMinMaxInfo.MinName
	result, err := u.ReadSiteInfo(key)
	return result, err
}

func (u *UrlRepository) ReadMaxResponseSite() (model.SiteResponseInfo, error) {
	u.log.Debug("Read Max Response Site From Repository")
	key := u.RepoSiteMinMaxInfo.MaxName
	result, err := u.ReadSiteInfo(key)
	return result, err
}

func (u *UrlRepository) ReadCountSiteRequest(s string) (uint64, error) {
	u.log.Debug(fmt.Sprintf("Read site count requests %s from repository", s))
	key, ok := u.RepoSiteName[s]
	if !ok {
		return 0, errors.New("incorrect site requested")
	}
	return u.RepoSiteCount[key].Load(), nil
}

func (u *UrlRepository) ReadCountMinRequest() uint64 {
	u.log.Debug("Read Max count request to Repository")
	return *u.RepoSiteMinMaxStat.MinCount
}

func (u *UrlRepository) ReadCountMaxRequest() uint64 {
	u.log.Debug("Read Max count request to Repository")
	return *u.RepoSiteMinMaxStat.MaxCount
}

func (u *UrlRepository) WriteCountSiteRequest(s string) {
	u.log.Debug(fmt.Sprintf("Write count site %s request to Repository", s))
	key := u.RepoSiteName[s]
	u.RepoSiteCount[key].Add(1)
}

func (u *UrlRepository) WriteCountMinRequest() {
	u.log.Debug("Write Min count request to Repository")
	atomic.AddUint64(u.RepoSiteMinMaxStat.MinCount, 1)
}

func (u *UrlRepository) WriteCountMaxRequest() {
	u.log.Debug("Write Max count request to Repository")
	atomic.AddUint64(u.RepoSiteMinMaxStat.MaxCount, 1)
}
