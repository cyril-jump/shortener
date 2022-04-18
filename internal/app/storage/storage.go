package storage

type Url struct {
	Short map[string]string
}

func NewUrl() *Url {
	return &Url{
		Short: make(map[string]string),
	}
}
