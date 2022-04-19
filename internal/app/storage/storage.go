package storage

//data base

type URL struct {
	Short map[string]string
}

func NewURL() *URL {
	return &URL{
		Short: make(map[string]string),
	}
}
