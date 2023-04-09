package controllers

import (
	"github.com/ewa-go/ewa"
	"github.com/ewa-go/ewa/consts"
)

type Api struct{}

func (Api) Get(route *ewa.Route) {
	route.SetDescription("Swagger")
	route.Handler = func(c *ewa.Context) error {
		b, err := c.Swagger.JSON()
		if err != nil {
			return c.SendString(422, err.Error())
		}
		return c.Send(200, consts.MIMEApplicationJSON, b)
	}
}
