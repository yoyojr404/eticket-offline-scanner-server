package scandevice

import (
	"context"

	"github.com/yoyojr404/internal/background"
)

type DeviceScan interface {
	background.Worker
	ScanDevice(ctx context.Context) error
}
