package mock

type CacheService struct {
	SetFn      func(string, string) error
	SetInvoked bool

	GetFn      func(string) (string, error)
	GetInvoked bool

	DeleteFn      func(string) error
	DeleteInvoked bool
}

func (c *CacheService) Set(key, value string) error {
	c.SetInvoked = true
	if c.SetFn == nil {
		return nil
	}
	return c.SetFn(key, value)
}

func (c *CacheService) Get(key string) (string, error) {
	c.GetInvoked = true
	if c.GetFn == nil {
		return "", nil
	}
	return c.GetFn(key)
}

func (c *CacheService) Delete(key string) error {
	c.DeleteInvoked = true
	if c.DeleteFn == nil {
		return nil
	}
	return c.DeleteFn(key)
}
