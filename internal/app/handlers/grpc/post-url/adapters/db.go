package adapters

type Repo struct {
	provider Provider
}

func New(p Provider) *Repo {
	return &Repo{
		provider: p,
	}
}

type Provider interface {
	SetShortURL(userID, shortURL, baseURL string) error
}

func (r *Repo) SetShortURLDB(userID, shortURL, baseURL string) error {

	return r.provider.SetShortURL(userID, shortURL, baseURL)
}
