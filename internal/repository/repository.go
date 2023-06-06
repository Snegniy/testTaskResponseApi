package repository

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Snegniy/testTaskResponseApi/internal/config"
	"github.com/Snegniy/testTaskResponseApi/internal/model"
	"go.uber.org/zap"
	"os"
)

type Repository interface {
	ReadSite(u *UrlRepository) (string, error)
}

type UrlRepository struct {
	RepoSite   map[string]model.SiteResponse
	RepoCount  map[string]model.SiteCount
	RepoMinMax model.MinMaxStat
}

func NewRepository(log *zap.Logger, cfg *config.Config) *UrlRepository {
	log.Debug("Register repository...")
	url := initData(log, cfg.UrlRepo.SitesFile)
	return url
}

func initData(log *zap.Logger, f string) *UrlRepository {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}
	defer file.Close()

	sites := map[string]model.SiteResponse{}
	sitesCount := map[string]model.SiteCount{}

	log.Debug("Read sites list..")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		site := scanner.Text()
		sites[site] = model.SiteResponse{}
		sitesCount[site] = model.SiteCount{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}
	return &UrlRepository{RepoSite: sites, RepoCount: sitesCount, RepoMinMax: model.MinMaxStat{}}
}

func (u *UrlRepository) ReadSite(site string) (float64, error) {
	_, ok := u.RepoSite[site]
	if !ok {
		return 0.0, errors.New("incorrect site requested")
	}

	time := u.RepoSite[site].Response
	if time == 0.0 {
		if u.RepoSite[site].Code == 0 {
			return 0.0, errors.New("site data not loaded. Please wait")
		}
		return 0.0, errors.New("site is unavailable")
	}
	return time, nil
}

func (u *UrlRepository) AddCount(r Repository) (*UrlRepository, error) {
	return u, nil
}

func (u *UrlRepository) GetMinMaxSite(m string) (float64, string, error) {
	return 0.0, "", nil
}
