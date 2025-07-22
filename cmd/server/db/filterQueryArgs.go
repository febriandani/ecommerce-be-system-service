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

func BuildQueryStatementGetFilterRegencies(baseQuery string, filter *systemPb.Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Search != "" {
		conditions = append(conditions, "r.name ILIKE '%' || ? || '%'")
		args = append(args, filter.Search)
	}

	if filter.Id != 0 {
		conditions = append(conditions, "r.province_id = ?")
		args = append(args, filter.Id)
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

func BuildQueryStatementGetFilterDistricts(baseQuery string, filter *systemPb.Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Search != "" {
		conditions = append(conditions, "d.name ILIKE '%' || ? || '%'")
		args = append(args, filter.Search)
	}

	if filter.Id != 0 {
		conditions = append(conditions, "d.regency_id = ?")
		args = append(args, filter.Id)
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

func BuildQueryStatementGetFilterSubDistricts(baseQuery string, filter *systemPb.Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Search != "" {
		conditions = append(conditions, "sd.name ILIKE '%' || ? || '%'")
		args = append(args, filter.Search)
	}

	if filter.Id != 0 {
		conditions = append(conditions, "sd.district_id = ?")
		args = append(args, filter.Id)
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
