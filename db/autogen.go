package db

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

func AutoGen(gorn *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(gorn)

	g.GenerateAllTable()

	g.Execute()
}
