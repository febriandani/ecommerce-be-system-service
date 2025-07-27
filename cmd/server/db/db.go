package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/febriandani/ecommerce-be-system-service/cmd/server/infra"
	systemPb "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
	"github.com/guregu/null"

	"github.com/sirupsen/logrus"
)

type DatabaseConfig struct {
	DB  *infra.DatabaseList
	log *logrus.Logger
}

func NewDatabaseConfig(db *infra.DatabaseList, logger *logrus.Logger) *DatabaseConfig {
	return &DatabaseConfig{
		DB:  db,
		log: logger,
	}
}

type DatabaseUser interface {
	GetProvinces(ctx context.Context, tx *sql.Tx, request *systemPb.Filter) (*systemPb.Provinces, error)
}

const (
// insertNotificationLog = `INSERT INTO public.systems
// (recipient, channel, subject, message, metadata, status, retry_count, max_retry, reason, sent_at, created_at, updated_at, trace_id, otp_id)
// VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) returning id;
// `
)

type provincesDB struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type regencyDB struct {
	ID         int64  `db:"id"`
	ProvinceID int64  `db:"province_id"`
	Name       string `db:"name"`
}

type districtDB struct {
	ID        int64  `db:"id"`
	RegencyID int64  `db:"regency_id"`
	Name      string `db:"name"`
}

type subDistrictDB struct {
	ID         int64       `db:"id"`
	DistrictID int64       `db:"district_id"`
	Name       string      `db:"name"`
	PostalCode null.String `db:"postal_code"`
}

func (dc *DatabaseConfig) GetProvinces(ctx context.Context, request *systemPb.Filter) ([]systemPb.Provinces, error) {
	var tplDB []provincesDB

	q := `SELECT id, name FROM master.provinces p `

	queryStatement, args2 := BuildQueryStatementGetFilterProvinces(q, request)
	query, args, err := dc.DB.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return nil, err
	}

	query = dc.DB.Backend.Read.Rebind(query)
	err = dc.DB.Backend.Read.Select(&tplDB, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	log.Println("GetProvinces query:", query, "args:", args)

	var result []systemPb.Provinces
	for _, p := range tplDB {
		result = append(result, systemPb.Provinces{
			Id:   p.ID,
			Name: p.Name,
		})
	}
	return result, nil
}

func (dc *DatabaseConfig) GetRegencies(ctx context.Context, request *systemPb.Filter) ([]systemPb.Regencies, error) {
	var tplDB []regencyDB

	q := `SELECT id, province_id, name FROM master.regencies r `

	queryStatement, args2 := BuildQueryStatementGetFilterRegencies(q, request)
	query, args, err := dc.DB.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return nil, err
	}

	query = dc.DB.Backend.Read.Rebind(query)
	err = dc.DB.Backend.Read.Select(&tplDB, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	log.Println("GetRegencies query:", query, "args:", args)

	var result []systemPb.Regencies
	for _, p := range tplDB {
		result = append(result, systemPb.Regencies{
			Id:         p.ID,
			ProvinceId: p.ProvinceID,
			Name:       p.Name,
		})
	}
	return result, nil
}

func (dc *DatabaseConfig) GetDistricts(ctx context.Context, request *systemPb.Filter) ([]systemPb.Districts, error) {
	var tplDB []districtDB

	q := `SELECT id, regency_id, name FROM master.districts d `

	queryStatement, args2 := BuildQueryStatementGetFilterDistricts(q, request)
	query, args, err := dc.DB.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return nil, err
	}

	query = dc.DB.Backend.Read.Rebind(query)
	err = dc.DB.Backend.Read.Select(&tplDB, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	log.Println("GetDistricts query:", query, "args:", args)

	var result []systemPb.Districts
	for _, p := range tplDB {
		result = append(result, systemPb.Districts{
			Id:        p.ID,
			RegencyId: p.RegencyID,
			Name:      p.Name,
		})
	}
	return result, nil
}

func (dc *DatabaseConfig) GetSubDistricts(ctx context.Context, request *systemPb.Filter) ([]systemPb.SubDistricts, error) {
	var tplDB []subDistrictDB

	q := `SELECT id, district_id, name, postal_code FROM master.sub_districts sd `

	queryStatement, args2 := BuildQueryStatementGetFilterSubDistricts(q, request)
	query, args, err := dc.DB.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return nil, err
	}

	query = dc.DB.Backend.Read.Rebind(query)
	err = dc.DB.Backend.Read.Select(&tplDB, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	log.Println("GetProvinces query:", query, "args:", args)

	dc.DB.Backend.Read.Close()

	var result []systemPb.SubDistricts
	for _, p := range tplDB {
		result = append(result, systemPb.SubDistricts{
			Id:         p.ID,
			DistrictId: p.DistrictID,
			Name:       p.Name,
			PostalCode: p.PostalCode.ValueOrZero(),
		})
	}
	return result, nil
}
