package main

import (
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	p "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
)

func main() {
	// @TODO: reuse cfg vars, and use as a separate go pkg to isolate deps
	err := postgres.Generate("../../db/gen", postgres.DBConnection{
		Host:       "localhost",
		Port:       5432,
		User:       "test",
		Password:   "test",
		SslMode:    "disable",
		DBName:     "test",
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
