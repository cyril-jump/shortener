package usecase

import (
	"github.com/labstack/gommon/log"

	"github.com/cyril-jump/shortener/internal/app/dto"
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

func (u *Usecase) ProcessPostURLBatch(userID string, urls []dto.ModelURLBatchRequest) ([]dto.ModelURLBatchResponse, error) {
	log.Info("processing PostURLBatch")

	resURL := dto.ModelURLBatchResponse{}
	var resURLs []dto.ModelURLBatchResponse

	hostName, err := u.cfg.Get("base_url_str")
	utils.CheckErr(err, "base_url_str")

	for _, val := range urls {
		shortURL := utils.Hash([]byte(val.BaseURL), hostName)
		err = u.repo.SetShortURLDB(userID, shortURL, val.BaseURL)
		if err != nil {
			return nil, err
		}
		resURL.ShortURL = shortURL
		resURL.CorID = val.CorID

		resURLs = append(resURLs, resURL)
	}

	return resURLs, err
}
