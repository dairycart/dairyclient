package dairyclient

type DairyclientV1 interface {
	ProductExists(sku string) (bool, error)
	GetProduct(sku string) (*Product, error)
	GetProducts(queryFilter map[string]string) ([]Product, error)
	CreateProduct(np ProductInput) (*Product, error)
	UpdateProduct(sku string, up ProductInput) (*Product, error)
	DeleteProduct(sku string) error
	GetProductRoot(rootID uint64) (*ProductRoot, error)
	GetProductRoots(queryFilter map[string]string) ([]ProductRoot, error)
	DeleteProductRoot(rootID uint64) error
	GetProductOptions(productID uint64, queryFilter map[string]string) ([]ProductOption, error)
	CreateProductOptionForProduct(productID uint64, no ProductOption) (*ProductOption, error)
	UpdateProductOption(optionID uint64, uo ProductOption) (*ProductOption, error)
	DeleteProductOption(optionID uint64) error
	CreateProductOptionValueForOption(optionID uint64, nv ProductOptionValue) (*ProductOptionValue, error)
	UpdateProductOptionValueForOption(valueID uint64, uv ProductOptionValue) (*ProductOptionValue, error)
	DeleteProductOptionValueForOption(optionID uint64) error
}
