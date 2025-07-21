package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/febriandani/ecommerce-be-system-service/cmd/server/infra"
	systemPb "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"

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
	insertNotificationLog = `INSERT INTO public.systems
	(recipient, channel, subject, message, metadata, status, retry_count, max_retry, reason, sent_at, created_at, updated_at, trace_id, otp_id)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) returning id;
	`
)

type provincesDB struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
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
