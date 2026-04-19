package hery

// Service defines the hery query interface.
type Service interface {
	HeryFunc(query string) (string, error)
}

type service struct{}

func (s *service) HeryFunc(query string) (string, error) {
	/*result, err := herypkg.ExecuteQuery(query)
	if err != nil {
		return "", err
	}*/
	return "result", nil
}
