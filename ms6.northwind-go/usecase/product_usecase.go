package usecase

import (
	"context"
	"northwind-go/domain"
	"time"
)

type productUsecase struct {
	repo    domain.ProductRepository
	timeout time.Duration
}

func NewProductUsecase(repo domain.ProductRepository, timeout time.Duration) domain.ProductUsecase {
	return &productUsecase{repo: repo, timeout: timeout}
}

func (u *productUsecase) ctx(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, u.timeout)
}

func (u *productUsecase) GetActiveProducts(c context.Context) ([]domain.Product, error) {
	ctx, cancel := u.ctx(c)
	defer cancel()
	return u.repo.FetchActive(ctx)
}

func (u *productUsecase) GetProduct(c context.Context, code string) (domain.Product, error) {
	ctx, cancel := u.ctx(c)
	defer cancel()
	return u.repo.GetByCode(ctx, code)
}

func (u *productUsecase) CreateProduct(c context.Context, p *domain.Product) error {
	ctx, cancel := u.ctx(c)
	defer cancel()
	return u.repo.Store(ctx, p)
}

func (u *productUsecase) UpdateProductPrice(c context.Context, code string, price float64) error {
	ctx, cancel := u.ctx(c)
	defer cancel()
	return u.repo.UpdatePrice(ctx, code, price)
}

func (u *productUsecase) DiscontinueProduct(c context.Context, code string) error {
	ctx, cancel := u.ctx(c)
	defer cancel()
	return u.repo.SoftDelete(ctx, code)
}
