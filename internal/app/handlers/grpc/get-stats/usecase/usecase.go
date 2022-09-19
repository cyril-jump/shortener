package usecase

import (
	"github.com/labstack/gommon/log"

	"github.com/cyril-jump/shortener/internal/app/dto"
)

type Usecase struct {
	repo Repo
}

type Repo interface {
	GetStatsDB() (dto.Stat, error)
}

func New(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) ProcessGetStats() (dto.Stat, error) {
	log.Info("processing GetStats")

	stats, err := u.repo.GetStatsDB()
	if err != nil {
		return dto.Stat{}, err
	}

	return stats, nil
}
