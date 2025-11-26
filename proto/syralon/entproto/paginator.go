package entproto

import "entgo.io/ent/dialect/sql"

func (o *ClassicalPaginator) OrderSelector() func(*sql.Selector) {
	if o.GetOrder() == nil {
		return func(selector *sql.Selector) {
		}
	}

	var opts []sql.OrderTermOption
	if o.GetOrder().Desc {
		opts = append(opts, sql.OrderDesc())
	}
	return sql.OrderByField(o.GetOrder().By, opts...).ToFunc()
}
