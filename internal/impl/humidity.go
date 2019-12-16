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

// HumidityServiceImpl structure to access database.
type HumidityServiceImpl struct {
	DB *sql.DB
}

// NewHumidityServiceImpl initialize humidity service and access to database.
func NewHumidityServiceImpl(db *sql.DB) *HumidityServiceImpl {
	return &HumidityServiceImpl{
		DB: db,
	}
}

// Add create a humidity metric.
func (s *HumidityServiceImpl) Add(ctx context.Context, req *message.Humidity) (*message.Humidity, error) {
	err := models.AddHumidity(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return req, nil
}

// GetHumidity get data points for humidity metrics.
func (s *HumidityServiceImpl) GetHumidity(req *message.HumidityRequest, stream service.HumidityService_GetHumidityServer) error {
	row, err := models.GetHumidity(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Humidity{}
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

// GetAvgHumidity get data points for average humidity metrics.
func (s *HumidityServiceImpl) GetAvgHumidity(req *message.HumidityRequest, stream service.HumidityService_GetAvgHumidityServer) error {
	row, err := models.GetAvgHumidity(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Humidity{}
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
