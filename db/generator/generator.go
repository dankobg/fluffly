package main

import (
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	p "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
)

func main() {
	cfg, _, e := config.New()
	if e != nil {
		panic(e)
	}

	err := postgres.Generate("../../db/gen", postgres.DBConnection{
		Host:       cfg.Database.Host,
		Port:       cfg.Database.Port,
		User:       cfg.Database.User,
		Password:   cfg.Database.Password,
		SslMode:    cfg.Database.SSLMode,
		DBName:     cfg.Database.DB,
		SchemaName: "public",
	}, template.Default(p.Dialect).
		UseSchema(func(schema metadata.Schema) template.Schema {
			return template.DefaultSchema(schema).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						return template.DefaultTableModel(table).
							UseField(func(column metadata.Column) template.TableModelField {
								defaultTableModelField := template.DefaultTableModelField(column)

								if schema.Name == "public" && table.Name == "animal" && column.Name == "properties" {
									defaultTableModelField.Type = template.NewType(dbcustom.JsonType{})
								}
								return defaultTableModelField
							})
					}),
				)
		}))
	if err != nil {
		panic(err)
	}
}
