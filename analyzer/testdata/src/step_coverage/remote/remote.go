package remote

type Client struct{}

func (*Client) Send() {}
func Send[T any](T)   {}
func Helper()         { new(Client).Send() }
