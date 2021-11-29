package pkg

import (
	"context"

	"tianxu.xin/phone/cloud/pkg/photo"
)

type Client interface {
	Login(ctx context.Context) error
	Requires2sa(ctx context.Context) error

	PhotoService(ctx context.Context) (photo.Service, error)
}
