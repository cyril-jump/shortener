package usecase

import "github.com/labstack/gommon/log"

type Usecase struct {
	repo Repo
}

type Repo interface {
	GetBaseURLDB(shortURL string) (string, error)
}

func New(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) ProcessGetBaseURL(shortURL string) (string, error) {
	log.Info("processing GetBaseURL")

	url, err := u.repo.GetBaseURLDB(shortURL)

	if err != nil {
		return "", err
	}

	return url, nil
}
