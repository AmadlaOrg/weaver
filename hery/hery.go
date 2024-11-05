package hery

type IHery interface {
	HeryFunc(query string) (string, error)
}

type SHery struct{}

func (s *SHery) HeryFunc(query string) (string, error) {
	result, err := herypkg.ExecuteQuery(query)
	if err != nil {
		return "", err
	}
	return result, nil
}
