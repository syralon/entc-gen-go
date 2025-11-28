package data

import (
	"entgo.io/ent/dialect/sql"

	"{{.module}}/ent"
	"{{.module}}/internal/conf"
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
