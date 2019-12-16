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

// TemperatureServiceImpl structure to access database.
type TemperatureServiceImpl struct {
	DB *sql.DB
}

// NewTemperatureServiceImpl initialize temperature service and access to database.
func NewTemperatureServiceImpl(db *sql.DB) *TemperatureServiceImpl {
	return &TemperatureServiceImpl{
		DB: db,
	}
}

// Add create a temperature metric.
func (s *TemperatureServiceImpl) Add(ctx context.Context, req *message.Temperature) (*message.Temperature, error) {
	err := models.AddTemperature(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return req, nil
}

// GetTemperature get data points for temperature metrics.
func (s *TemperatureServiceImpl) GetTemperature(req *message.TemperatureRequest, stream service.TemperatureService_GetTemperatureServer) error {
	row, err := models.GetTemperature(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Temperature{}
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

// GetAvgTemperature get data points for average temperature metrics.
func (s *TemperatureServiceImpl) GetAvgTemperature(req *message.TemperatureRequest, stream service.TemperatureService_GetAvgTemperatureServer) error {
	row, err := models.GetAvgTemperature(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Temperature{}
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
