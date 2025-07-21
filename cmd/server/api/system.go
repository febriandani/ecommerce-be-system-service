package api

import (
	"context"
	"fmt"

	udb "github.com/febriandani/ecommerce-be-system-service/cmd/server/db"
	"github.com/febriandani/ecommerce-be-system-service/cmd/server/utils"
	systemPb "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
	"github.com/sirupsen/logrus"
)

// SystemServer implements the Users gRPC service
type SystemServer struct {
	systemPb.UnimplementedSystemsServer // Embeds the default implementation
	db                                  *udb.DatabaseConfig
	log                                 *logrus.Logger
}

// NewSystemServer initializes a new SystemServer
func NewSystemServer(database *udb.DatabaseConfig, logger *logrus.Logger) *SystemServer {
	return &SystemServer{db: database, log: logger}
}

func (s *SystemServer) GetProvinces(ctx context.Context, request *systemPb.Filter) (*systemPb.ProvincesResponse, error) {
	if request.GetTraceId() == "" {
		request.TraceId = utils.NewTraceID()
	}

	logger := utils.NewLoggerWithTrace(s.log, request.GetTraceId(), "")
	logger.Info("Api [GetProvinces] started processing")

	// Ambil data dari DB
	tpl, err := s.db.GetProvinces(ctx, request)
	if err != nil {
		logger.Error(err, "Gagal ambil provinces dari DB")
		return &systemPb.ProvincesResponse{
			Meta: &systemPb.Meta{
				Code:            500,
				Status:          "error",
				Message:         "Gagal mengambil data provinces",
				InternalMessage: err.Error(),
				TraceId:         request.TraceId,
			},
			Data: nil,
		}, nil
	}

	// Konversi slice ke proto
	var protoProvinces []*systemPb.Provinces
	for _, p := range tpl {
		protoProvinces = append(protoProvinces, &systemPb.Provinces{
			Id:   p.Id,
			Name: p.Name,
		})
	}

	logger.Info(fmt.Sprintf("Berhasil dapat %d provinsi", len(protoProvinces)))

	// Return langsung ke ProvincesResponse (tanpa Any)
	return &systemPb.ProvincesResponse{
		Meta: &systemPb.Meta{
			Code:    200,
			Status:  "success",
			Message: "provinces found",
			TraceId: request.TraceId,
		},
		Data: &systemPb.ProvincesList{
			Provinces: protoProvinces,
		},
	}, nil
}
