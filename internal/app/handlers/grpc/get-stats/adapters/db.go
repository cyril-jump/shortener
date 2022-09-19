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
	GetStats() (dto.Stat, error)
}

func (r *Repo) GetStatsDB() (dto.Stat, error) {

	stats, err := r.provider.GetStats()
	if err != nil {
		return dto.Stat{}, err
	}

	return stats, nil
}
