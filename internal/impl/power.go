package impl

import (
	"context"
	"database/sql"
	"time"

	"github.com/afairon/smeter/message"
	"github.com/afairon/smeter/models"
	"github.com/afairon/smeter/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PowerServiceImpl structure to access database.
type PowerServiceImpl struct {
	DB *sql.DB
}

// NewPowerServiceImpl initialize power service and access to database.
func NewPowerServiceImpl(db *sql.DB) *PowerServiceImpl {
	return &PowerServiceImpl{
		DB: db,
	}
}

// Add create a power metric.
func (s *PowerServiceImpl) Add(ctx context.Context, req *message.Power) (*message.Power, error) {
	err := models.AddPower(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return req, nil
}

// GetPower get data points for power metrics.
func (s *PowerServiceImpl) GetPower(req *message.PowerRequest, stream service.PowerService_GetPowerServer) error {
	row, err := models.GetPower(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Power{}
		timestamp := time.Time{}
		err = row.Scan(&timestamp, &resp.SensorID, &resp.Value)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		resp.Time = timestamp.Unix()
		stream.Send(&resp)
	}

	return err
}

// GetAvgConsumption get data points for average consumption metrics.
func (s *PowerServiceImpl) GetAvgConsumption(req *message.ConsumptionRequest, stream service.PowerService_GetAvgConsumptionServer) error {
	row, err := models.GetAvgConsumption(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Energy{}
		timestamp := time.Time{}
		err = row.Scan(&timestamp, &resp.SensorID, &resp.Value)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		resp.Time = timestamp.Unix()
		stream.Send(&resp)
	}

	return err
}
