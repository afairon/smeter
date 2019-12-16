package impl

import (
	"context"
	"database/sql"

	"github.com/afairon/smeter/message"
	"github.com/afairon/smeter/models"
	"github.com/afairon/smeter/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SensorServiceImpl structure to access database.
type SensorServiceImpl struct {
	DB *sql.DB
}

// NewSensorServiceImpl initialize sensor service and access to database.
func NewSensorServiceImpl(db *sql.DB) *SensorServiceImpl {
	return &SensorServiceImpl{
		DB: db,
	}
}

// Add creates sensor in database.
func (s *SensorServiceImpl) Add(ctx context.Context, req *message.Sensor) (*message.Sensor, error) {
	err := models.AddSensor(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return req, nil
}

// Get retrieves sensors from database.
func (s *SensorServiceImpl) Get(req *message.SensorsRequest, stream service.SensorService_GetServer) error {
	if req.ID != 0 {
		sensor, err := models.GetSensor(s.DB, &message.Sensor{ID: req.GetID()})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		stream.Send(sensor)
		return err
	}
	row, err := models.GetSensors(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Sensor{}
		err = row.Scan(&resp.ID, &resp.DeviceID, &resp.Type, &resp.Name, &resp.Active)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		stream.Send(&resp)
	}

	return err
}

// Update updates the device ID and sensor name.
func (s *SensorServiceImpl) Update(ctx context.Context, req *message.Sensor) (*message.Empty, error) {
	if err := models.UpdateSensor(s.DB, req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &message.Empty{}, nil
}

// Delete deletes the sensor.
func (s *SensorServiceImpl) Delete(ctx context.Context, req *message.Sensor) (*message.Empty, error) {
	if err := models.DeleteSensor(s.DB, req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &message.Empty{}, nil
}

// Count counts the number of sensors.
func (s *SensorServiceImpl) Count(ctx context.Context, req *message.SensorCountRequest) (*message.SensorCount, error) {
	count, err := models.CountSensor(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return count, nil
}
