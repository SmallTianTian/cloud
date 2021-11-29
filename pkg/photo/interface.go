package photo

import (
	"context"
	"io"
)

type Service interface {
	Albums(ctx context.Context) (map[string]Album, error)
	All(ctx context.Context) ([]Item, error)
}

type Album interface {
	Photos(ctx context.Context) ([]Item, error)
}

type Item interface {
	Name() string
	WriteTo(ctx context.Context, file string) error
	Reader(ctx context.Context) (io.ReadCloser, error)
	Delete(ctx context.Context) error
}
