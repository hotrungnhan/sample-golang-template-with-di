package models

import (
	"github.com/bluele/factory-go/factory"
	"github.com/go-faker/faker/v4"
)

// 'Location: "Tokyo"' is default value.
var ShortenUrlFactory = factory.NewFactory(
	&ShortenUrl{},
).Attr("OriginalUrl", func(args factory.Args) (interface{}, error) {
	return faker.URL(), nil
}).OnCreate(func(args factory.Args) error {
	ins := args.Instance().(*ShortenUrl)

	ins.BeforeCreate(nil)

	return nil
})
