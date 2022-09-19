package usecase

import (
	"github.com/labstack/gommon/log"

	"github.com/cyril-jump/shortener/internal/app/dto"
)

type Usecase struct {
	repo Repo
	cfg  Config
}

type Repo interface {
	GetAllURLsByUserIDDB(userID string) ([]dto.ModelURL, error)
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

func (u *Usecase) ProcessGetAllURLsByUserID(userID string) ([]dto.ModelURL, error) {
	log.Info("processing GetAllURLsByUserID")

	urls, err := u.repo.GetAllURLsByUserIDDB(userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}
