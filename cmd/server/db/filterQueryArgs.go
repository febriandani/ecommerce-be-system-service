package db

import (
	"strings"

	systemPb "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
)

func BuildQueryStatementGetFilterProvinces(baseQuery string, filter *systemPb.Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Search != "" {
		conditions = append(conditions, "p.name ILIKE '%' || ? || '%'")
		args = append(args, filter.Search)
	}

	if len(conditions) > 0 {
		whereClause := "WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
	}

	if filter.Page != 0 && filter.Limit != 0 {
		baseQuery += " OFFSET ((? - 1) * ?) ROWS FETCH NEXT ? ROWS ONLY"
		args = append(args, filter.Page, filter.Limit, filter.Limit)
	}

	return baseQuery, args
}
