package domain

import "context"

type Product struct {
	ProductCode            string  `json:"product_code"             bson:"product_code"`
	ProductName            string  `json:"product_name"             bson:"product_name"`
	Description            string  `json:"description"              bson:"description"`
	StandardCost           float64 `json:"standard_cost"            bson:"standard_cost"`
	ListPrice              float64 `json:"list_price"               bson:"list_price"`
	ReorderLevel           int     `json:"reorder_level"            bson:"reorder_level"`
	TargetLevel            int     `json:"target_level"             bson:"target_level"`
	QuantityPerUnit        string  `json:"quantity_per_unit"        bson:"quantity_per_unit"`
	Discontinued           bool    `json:"discontinued"             bson:"discontinued"`
	MinimumReorderQuantity int     `json:"minimum_reorder_quantity" bson:"minimum_reorder_quantity"`
	Category               string  `json:"category"                 bson:"category"`
}

type ProductRepository interface {
	FetchActive(ctx context.Context) ([]Product, error)
	GetByCode(ctx context.Context, code string) (Product, error)
	Store(ctx context.Context, p *Product) error
	UpdatePrice(ctx context.Context, productCode string, newPrice float64) error
	SoftDelete(ctx context.Context, productCode string) error
}

type ProductUsecase interface {
	GetActiveProducts(ctx context.Context) ([]Product, error)
	GetProduct(ctx context.Context, code string) (Product, error)
	CreateProduct(ctx context.Context, p *Product) error
	UpdateProductPrice(ctx context.Context, code string, price float64) error
	DiscontinueProduct(ctx context.Context, code string) error
}
