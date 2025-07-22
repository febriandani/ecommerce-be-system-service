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

func (s *SystemServer) GetRegencies(ctx context.Context, request *systemPb.Filter) (*systemPb.RegenciesResponse, error) {
	if request.GetTraceId() == "" {
		request.TraceId = utils.NewTraceID()
	}

	logger := utils.NewLoggerWithTrace(s.log, request.GetTraceId(), "")
	logger.Info("Api [GetRegencies] started processing")

	// Ambil data dari DB
	tpl, err := s.db.GetRegencies(ctx, request)
	if err != nil {
		logger.Error(err, "Gagal ambil regency dari DB")
		return &systemPb.RegenciesResponse{
			Meta: &systemPb.Meta{
				Code:            500,
				Status:          "error",
				Message:         "Gagal mengambil data regency",
				InternalMessage: err.Error(),
				TraceId:         request.TraceId,
			},
			Data: nil,
		}, nil
	}

	// Konversi slice ke proto
	var protoRegencies []*systemPb.Regencies
	for _, p := range tpl {
		protoRegencies = append(protoRegencies, &systemPb.Regencies{
			Id:         p.Id,
			ProvinceId: p.ProvinceId,
			Name:       p.Name,
		})
	}

	logger.Info(fmt.Sprintf("Berhasil dapat %d regency", len(protoRegencies)))

	// Return langsung ke ProvincesResponse (tanpa Any)
	return &systemPb.RegenciesResponse{
		Meta: &systemPb.Meta{
			Code:    200,
			Status:  "success",
			Message: "regencies found",
			TraceId: request.TraceId,
		},
		Data: &systemPb.RegenciesList{
			Regencies: protoRegencies,
		},
	}, nil
}

func (s *SystemServer) GetDistricts(ctx context.Context, request *systemPb.Filter) (*systemPb.DistrictsResponse, error) {
	if request.GetTraceId() == "" {
		request.TraceId = utils.NewTraceID()
	}

	logger := utils.NewLoggerWithTrace(s.log, request.GetTraceId(), "")
	logger.Info("Api [GetDistricts] started processing")

	// Ambil data dari DB
	tpl, err := s.db.GetDistricts(ctx, request)
	if err != nil {
		logger.Error(err, "Gagal ambil districts dari DB")
		return &systemPb.DistrictsResponse{
			Meta: &systemPb.Meta{
				Code:            500,
				Status:          "error",
				Message:         "Gagal mengambil data districts",
				InternalMessage: err.Error(),
				TraceId:         request.TraceId,
			},
			Data: nil,
		}, nil
	}

	// Konversi slice ke proto
	var protoDistricts []*systemPb.Districts
	for _, p := range tpl {
		protoDistricts = append(protoDistricts, &systemPb.Districts{
			Id:        p.Id,
			RegencyId: p.RegencyId,
			Name:      p.Name,
		})
	}

	logger.Info(fmt.Sprintf("Berhasil dapat %d districts", len(protoDistricts)))

	// Return langsung ke ProvincesResponse (tanpa Any)
	return &systemPb.DistrictsResponse{
		Meta: &systemPb.Meta{
			Code:    200,
			Status:  "success",
			Message: "districts found",
			TraceId: request.TraceId,
		},
		Data: &systemPb.DistrictsList{
			Districts: protoDistricts,
		},
	}, nil
}

func (s *SystemServer) GetSubDistricts(ctx context.Context, request *systemPb.Filter) (*systemPb.SubDistrictsResponse, error) {
	if request.GetTraceId() == "" {
		request.TraceId = utils.NewTraceID()
	}

	logger := utils.NewLoggerWithTrace(s.log, request.GetTraceId(), "")
	logger.Info("Api [GetSubDistricts] started processing")

	// Ambil data dari DB
	tpl, err := s.db.GetSubDistricts(ctx, request)
	if err != nil {
		logger.Error(err, "Gagal ambil sub districts dari DB")
		return &systemPb.SubDistrictsResponse{
			Meta: &systemPb.Meta{
				Code:            500,
				Status:          "error",
				Message:         "Gagal mengambil data sub districts",
				InternalMessage: err.Error(),
				TraceId:         request.TraceId,
			},
			Data: nil,
		}, nil
	}

	// Konversi slice ke proto
	var protoSubDistricts []*systemPb.SubDistricts
	for _, p := range tpl {
		protoSubDistricts = append(protoSubDistricts, &systemPb.SubDistricts{
			Id:         p.Id,
			DistrictId: p.DistrictId,
			Name:       p.Name,
		})
	}

	logger.Info(fmt.Sprintf("Berhasil dapat %d sub districts", len(protoSubDistricts)))

	// Return langsung ke ProvincesResponse (tanpa Any)
	return &systemPb.SubDistrictsResponse{
		Meta: &systemPb.Meta{
			Code:    200,
			Status:  "success",
			Message: "sub districts found",
			TraceId: request.TraceId,
		},
		Data: &systemPb.SubDistrictsList{
			Subdistricts: protoSubDistricts,
		},
	}, nil
}
