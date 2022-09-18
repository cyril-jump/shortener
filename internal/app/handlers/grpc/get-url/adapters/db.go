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
	GetBaseURL(shortURL string) (string, error)
}

func (r *Repo) GetBaseURLDB(shortURL string) (string, error) {

	url, err := r.provider.GetBaseURL(shortURL)
	if err != nil {
		return "", err
	}

	return url, nil
}
