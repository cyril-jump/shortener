package usecase

import (
	"github.com/labstack/gommon/log"

	"github.com/cyril-jump/shortener/internal/app/utils"
)

type Usecase struct {
	repo Repo
	cfg  Config
}

type Repo interface {
	SetShortURLDB(userID, shortURL, baseURL string) error
}

type Config interface {
	Get(key string) (string, error)
}

func New(repo Repo, cfg Config) *Usecase {
	return &Usecase{
		repo: repo,
		cfg:  cfg,
	}
}

func (u *Usecase) ProcessPostURL(userID string, baseURL string) (string, error) {
	log.Info("processing PostURL")

	hostName, err := u.cfg.Get("base_url_str")
	utils.CheckErr(err, "base_url_str")

	shortURL := utils.Hash([]byte(baseURL), hostName)

	err = u.repo.SetShortURLDB(userID, shortURL, baseURL)
	if err != nil {
		return shortURL, err
	}

	return shortURL, err
}
