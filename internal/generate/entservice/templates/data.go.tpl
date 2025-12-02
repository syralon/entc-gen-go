package data

import (
	"entgo.io/ent/dialect/sql"

	"{{.module}}/ent"
	"{{.module}}/internal/conf"

    # _ "github.com/mattn/go-sqlite3"
    _ "github.com/go-sql-driver/mysql"
)

func NewData(conf *conf.Data) (*ent.Client, error) {
	drv, err := sql.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	if err != nil {
		return nil, err
	}
	client := ent.NewClient(ent.Driver(drv))
	return client, nil
}
