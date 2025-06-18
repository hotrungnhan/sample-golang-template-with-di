package main

import (
	"github.com/hotrungnhan/surl/models"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./generated/queries",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(models.ShortenUrl{})
	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, models.DnsRecord{})
	// Generate the code
	g.Execute()
}
