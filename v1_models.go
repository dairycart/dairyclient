package dairyclient

import (
	"time"
)

////////////////////////////////////////////////////////
//                                                    //
//                      Helpers                       //
//                                                    //
////////////////////////////////////////////////////////

// DBRow is meant to represent the base columns that every database table should have
type DBRow struct {
	ID         uint64    `json:"id"`
	CreatedOn  time.Time `json:"created_on"`
	UpdatedOn  time.Time `json:"updated_on,omitempty"`
	ArchivedOn time.Time `json:"archived_on,omitempty"`
}

// ListResponse is a generic list response struct containing values that represent
// pagination, meant to be embedded into other object response structs
type ListResponse struct {
	Count uint64      `json:"count"`
	Limit uint8       `json:"limit"`
	Page  uint64      `json:"page"`
	Data  interface{} `json:"data"`
}

// ErrorResponse is a handy struct we can respond with in the event we have an error to report
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

////////////////////////////////////////////////////////
//                                                    //
//                       Users                        //
//                                                    //
////////////////////////////////////////////////////////

// User represents a Dairycart user
type User struct {
	DBRow
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
}

// UserCreationInput represents the payload used to create a Dairycart user
type UserCreationInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"is_admin"`
}

// UserLoginInput represents the payload used to log in a Dairycart user
type UserLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserUpdateInput represents the payload used to update a Dairycart user
type UserUpdateInput struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

////////////////////////////////////////////////////////
//                                                    //
//                      Products                      //
//                                                    //
////////////////////////////////////////////////////////

// ProductRoot represents the object that products inherit from
type ProductRoot struct {
	DBRow

	// Basic Info
	Name               string    `json:"name"`
	Subtitle           string    `json:"subtitle"`
	Description        string    `json:"description"`
	SKUPrefix          string    `json:"sku_prefix"`
	Manufacturer       string    `json:"manufacturer"`
	Brand              string    `json:"brand"`
	AvailableOn        time.Time `json:"available_on"`
	QuantityPerPackage uint32    `json:"quantity_per_package"`

	// Pricing Fields
	Taxable bool    `json:"taxable"`
	Cost    float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight float32 `json:"package_weight"`
	PackageHeight float32 `json:"package_height"`
	PackageWidth  float32 `json:"package_width"`
	PackageLength float32 `json:"package_length"`

	Options  []ProductOption `json:"options"`
	Products []Product       `json:"products"`
}

// Product describes something a user can buy
type Product struct {
	DBRow
	// Basic Info
	ProductRootID      uint64 `json:"product_root_id"`
	Name               string `json:"name"`
	Subtitle           string `json:"subtitle"`
	Description        string `json:"description"`
	OptionSummary      string `json:"option_summary"`
	SKU                string `json:"sku"`
	UPC                string `json:"upc"`
	Manufacturer       string `json:"manufacturer"`
	Brand              string `json:"brand"`
	Quantity           uint32 `json:"quantity"`
	QuantityPerPackage uint32 `json:"quantity_per_package"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight float32 `json:"package_weight"`
	PackageHeight float32 `json:"package_height"`
	PackageWidth  float32 `json:"package_width"`
	PackageLength float32 `json:"package_length"`

	ApplicableOptionValues []ProductOptionValue `json:"applicable_options,omitempty"`

	AvailableOn time.Time `json:"available_on"`
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	// Core Product stuff
	Name         string `json:"name"`
	Subtitle     string `json:"subtitle"`
	Description  string `json:"description"`
	SKU          string `json:"sku"`
	UPC          string `json:"upc"`
	Manufacturer string `json:"manufacturer"`
	Brand        string `json:"brand"`
	Quantity     uint32 `json:"quantity"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight      float32 `json:"package_weight"`
	PackageHeight      float32 `json:"package_height"`
	PackageWidth       float32 `json:"package_width"`
	PackageLength      float32 `json:"package_length"`
	QuantityPerPackage uint32  `json:"quantity_per_package"`

	AvailableOn time.Time `json:"available_on"`

	// Other things
	Options []ProductOptionCreationInput `json:"options"`
}

// ProductUpdateInput is a struct that represents a product update body
type ProductUpdateInput struct {
	// Core Product stuff
	Name         string `json:"name"`
	Subtitle     string `json:"subtitle"`
	Description  string `json:"description"`
	SKU          string `json:"sku"`
	UPC          string `json:"upc"`
	Manufacturer string `json:"manufacturer"`
	Brand        string `json:"brand"`
	Quantity     uint32 `json:"quantity"`

	// Pricing Fields
	Taxable   bool    `json:"taxable"`
	Price     float32 `json:"price"`
	OnSale    bool    `json:"on_sale"`
	SalePrice float32 `json:"sale_price"`
	Cost      float32 `json:"cost"`

	// Product Dimensions
	ProductWeight float32 `json:"product_weight"`
	ProductHeight float32 `json:"product_height"`
	ProductWidth  float32 `json:"product_width"`
	ProductLength float32 `json:"product_length"`

	// Package dimensions
	PackageWeight      float32 `json:"package_weight"`
	PackageHeight      float32 `json:"package_height"`
	PackageWidth       float32 `json:"package_width"`
	PackageLength      float32 `json:"package_length"`
	QuantityPerPackage uint32  `json:"quantity_per_package"`

	AvailableOn time.Time `json:"available_on"`
}

////////////////////////////////////////////////////////
//                                                    //
//                  Product Options                   //
//                                                    //
////////////////////////////////////////////////////////

// ProductOption represents a products variant options. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size.
type ProductOption struct {
	DBRow
	ProductRootID uint64               `json:"product_root_id"`
	Name          string               `json:"name"`
	Values        []ProductOptionValue `json:"values"`
}

// ProductOptionUpdateInput is a struct to use for updating product options
type ProductOptionUpdateInput struct {
	Name string `json:"name"`
}

// ProductOptionCreationInput is a struct to use for creating product options
type ProductOptionCreationInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

////////////////////////////////////////////////////////
//                                                    //
//               Product Option Values                //
//                                                    //
////////////////////////////////////////////////////////

// ProductOptionValue represents a product's option values. If you have a t-shirt that comes in three colors
// and three sizes, then there are two ProductOptions for that base_product, color and size, and six ProductOptionValues,
// One for each color and one for each size.
type ProductOptionValue struct {
	DBRow
	ProductOptionID uint64 `json:"product_option_id"`
	Value           string `json:"value"`
}

// ProductOptionValueCreationInput is a struct to use for creating product option values
type ProductOptionValueCreationInput struct {
	ProductOptionID uint64
	Value           string `json:"value"`
}

// ProductOptionValueUpdateInput is a struct to use for updating product option values
type ProductOptionValueUpdateInput struct {
	Value string `json:"value"`
}

////////////////////////////////////////////////////////
//                                                    //
//                     Discounts                      //
//                                                    //
////////////////////////////////////////////////////////

// Discount represents pricing changes that apply temporarily to products
type Discount struct {
	DBRow
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Amount        float32   `json:"amount"`
	StartsOn      time.Time `json:"starts_on"`
	ExpiresOn     time.Time `json:"expires_on"`
	RequiresCode  bool      `json:"requires_code"`
	Code          string    `json:"code,omitempty"`
	LimitedUse    bool      `json:"limited_use"`
	NumberOfUses  int64     `json:"number_of_uses,omitempty"`
	LoginRequired bool      `json:"login_required"`
}

// DiscountCreationInput represents user input for creating new discounts
type DiscountCreationInput struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Amount        float32   `json:"amount"`
	StartsOn      time.Time `json:"starts_on"`
	ExpiresOn     time.Time `json:"expires_on"`
	RequiresCode  bool      `json:"requires_code"`
	Code          string    `json:"code"`
	LimitedUse    bool      `json:"limited_use"`
	NumberOfUses  int64     `json:"number_of_uses"`
	LoginRequired bool      `json:"login_required"`
}
