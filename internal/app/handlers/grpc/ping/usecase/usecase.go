package usecase

import (
	"github.com/labstack/gommon/log"
)

type Usecase struct {
	repo Repo
}

type Repo interface {
	PingDB() error
}

func New(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) ProcessPing() error {
	log.Info("processing Ping")

	return u.repo.PingDB()
}
