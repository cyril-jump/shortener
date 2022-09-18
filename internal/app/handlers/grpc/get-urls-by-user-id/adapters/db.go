package adapters

import "github.com/cyril-jump/shortener/internal/app/dto"

type Repo struct {
	provider Provider
}

func New(p Provider) *Repo {
	return &Repo{
		provider: p,
	}
}

type Provider interface {
	GetAllURLsByUserID(userID string) ([]dto.ModelURL, error)
}

func (r *Repo) GetAllURLsByUserIDDB(userID string) ([]dto.ModelURL, error) {

	urls, err := r.provider.GetAllURLsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}
