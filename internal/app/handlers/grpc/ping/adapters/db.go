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
	Ping() error
}

func (r *Repo) PingDB() error {

	return r.provider.Ping()
}
