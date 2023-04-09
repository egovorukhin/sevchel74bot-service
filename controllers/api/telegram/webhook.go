package telegram

import "github.com/ewa-go/ewa"

type Webhook struct{}

func (Webhook) Get(route *ewa.Route) {
	route.Handler = func(c *ewa.Context) error {
		return nil
	}
}
