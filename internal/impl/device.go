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

// DeviceServiceImpl structure to access database.
type DeviceServiceImpl struct {
	DB *sql.DB
}

// NewDeviceServiceImpl initialize device service and access to database.
func NewDeviceServiceImpl(db *sql.DB) *DeviceServiceImpl {
	return &DeviceServiceImpl{
		DB: db,
	}
}

// Add creates new device in database.
func (s *DeviceServiceImpl) Add(ctx context.Context, req *message.Device) (*message.Device, error) {
	err := models.AddDevice(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return req, nil
}

// Get retrieve devices from database.
func (s *DeviceServiceImpl) Get(req *message.DevicesRequest, stream service.DeviceService_GetServer) error {
	if req.ID != 0 {
		device, err := models.GetDevice(s.DB, &message.Device{ID: req.GetID()})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		stream.Send(device)
		return err
	}
	row, err := models.GetDevices(s.DB, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for row.Next() {
		resp := message.Device{}
		err = row.Scan(&resp.ID, &resp.Name, &resp.Active)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		stream.Send(&resp)
	}

	return err
}

// Update updates the device name.
func (s *DeviceServiceImpl) Update(ctx context.Context, req *message.Device) (*message.Empty, error) {
	if err := models.UpdateDevice(s.DB, req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &message.Empty{}, nil
}

// Delete deletes the device.
func (s *DeviceServiceImpl) Delete(ctx context.Context, req *message.Device) (*message.Empty, error) {
	if err := models.DeleteDevice(s.DB, req); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &message.Empty{}, nil
}

// Count counts the number of devices.
func (s *DeviceServiceImpl) Count(ctx context.Context, req *message.DeviceCountRequest) (*message.DeviceCount, error) {
	count, err := models.CountDevice(s.DB, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return count, nil
}
