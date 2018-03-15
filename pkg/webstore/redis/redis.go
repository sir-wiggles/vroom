package redis

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (r *Service) Get(key string) error {
	return nil
}

func (r *Service) Set(key, value string) error {
	return nil
}

func (r *Service) Delete(key string) error {
	return nil
}
