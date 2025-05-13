package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/model"
)

type Product interface {
	Creat(ctx context.Context, product *model.Product) error
}
