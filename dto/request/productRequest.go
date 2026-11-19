package request

type CreateProductRequest struct {
	SKU         string `form:"sku" binding:"required"`
	Name        string `form:"name" binding:"required,min=3"`
	Description string `form:"description"`
	Stock       *int   `form:"stock" binding:"required,gte=0"`
	Price       *int64 `form:"price" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	SKU         string `form:"sku"`
	Name        string `form:"name"`
	Description string `form:"description"`
	Stock       *int   `form:"stock"`
	Price       *int64 `form:"price"`
}
