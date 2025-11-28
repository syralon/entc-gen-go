package data

import (
	"entgo.io/ent/dialect/sql"

	"github.com/syralon/entc-gen-go/example/ent"
	"github.com/syralon/entc-gen-go/example/internal/conf"
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
