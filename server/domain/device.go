package domain

import "context"

type Device struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	DeviceID     string `json:"device_id" bson:"device_id"`
	DeviceName   string `json:"device_name" bson:"device_name"`
	DeviceStatus string `json:"device_status" bson:"device_status"`
}

type DeviceRepository interface {
	GetAllDevices(ctx context.Context) ([]*Device, error)
	GetByID(ctx context.Context, id string) (*Device, error)
	Update(ctx context.Context, id string, device *Device) error
	Save(ctx context.Context, device *Device) error
}
