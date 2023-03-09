package gcs

type InMemoryClient struct{}

func NewInMemoryClient() Client {

	return &InMemoryClient{}
}

func (c *InMemoryClient) Update(name string, data []byte) error {

	return nil
}

func (c *InMemoryClient) Close() error { return nil }
