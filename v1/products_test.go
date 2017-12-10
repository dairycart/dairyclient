package dairyclient_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
)

func buildNotFoundProductResponse(sku string) string {
	return fmt.Sprintf(`
		{
			"status": 404,
			"message": "The product you were looking for (sku '%s') does not exist"
		}
	`, sku)
}

func TestProductExists(t *testing.T) {
	t.Parallel()

	existentSKU := "existent_sku"
	nonexistentSKU := "nonexistent_sku"

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", existentSKU):    generateHeadHandler(t, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU): generateHeadHandler(t, http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(_t *testing.T) {
		exists, err := c.ProductExists(existentSKU)
		assert.Nil(t, err)
		assert.True(t, exists)

	})

	t.Run("with nonexistent product", func(_t *testing.T) {
		exists, err := c.ProductExists(nonexistentSKU)
		assert.Nil(t, err)
		assert.False(t, exists)
	})
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	goodResponseSKU := "good"
	nonexistentSKU := "nonexistent"
	badResponseSKU := "bad"

	exampleResponse := fmt.Sprintf(`
		{
			"name": "Your Favorite Band's T-Shirt",
			"subtitle": "A t-shirt you can wear",
			"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			"option_summary": "Size: Small, Color: Red",
			"sku": "%s",
			"upc": "",
			"manufacturer": "Record Company",
			"brand": "Your Favorite Band",
			"quantity": 666,
			"quantity_per_package": 1,
			"taxable": true,
			"price": 20,
			"on_sale": false,
			"sale_price": 0,
			"cost": 10,
			"product_weight": 1,
			"product_height": 5,
			"product_width": 5,
			"product_length": 5,
			"package_weight": 1,
			"package_height": 5,
			"package_width": 5,
			"package_length": 5
		}
	`, goodResponseSKU)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", goodResponseSKU): generateGetHandler(t, exampleResponse, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", badResponseSKU):  generateGetHandler(t, exampleBadJSON, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU):  generateGetHandler(t, buildNotFoundProductResponse(nonexistentSKU), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(_t *testing.T) {
		expected := &models.Product{
			Name:               "Your Favorite Band's T-Shirt",
			Subtitle:           "A t-shirt you can wear",
			Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			OptionSummary:      "Size: Small, Color: Red",
			SKU:                goodResponseSKU,
			UPC:                "",
			Manufacturer:       "Record Company",
			Brand:              "Your Favorite Band",
			Quantity:           666,
			QuantityPerPackage: 1,
			Taxable:            true,
			Price:              20,
			OnSale:             false,
			Cost:               10,
			ProductWeight:      1,
			ProductHeight:      5,
			ProductWidth:       5,
			ProductLength:      5,
			PackageWeight:      1,
			PackageHeight:      5,
			PackageWidth:       5,
			PackageLength:      5,
		}
		actual, err := c.GetProduct(goodResponseSKU)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected product doesn't match actual product")
	})

	t.Run("nonexistent product", func(_t *testing.T) {
		_, err := c.GetProduct(nonexistentSKU)
		assert.NotNil(t, err)
	})

	t.Run("bad response from server", func(_t *testing.T) {
		_, err := c.GetProduct(badResponseSKU)
		assert.NotNil(t, err)
	})

	t.Run("with request error", func(_t *testing.T) {
		ts.Close()
		_, err := c.GetProduct(exampleSKU)
		assert.NotNil(t, err)
	})
}

func TestGetProducts(t *testing.T) {
	t.Parallel()

	exampleGoodResponse := `
		{
			"count": 5,
			"limit": 25,
			"page": 1,
			"products": [{
					"id": 1,
					"product_root_id": 1,
					"name": "Your Favorite Band's T-Shirt",
					"subtitle": "A t-shirt you can wear",
					"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					"option_summary": "Size: Small, Color: Red",
					"sku": "t-shirt-small-red",
					"manufacturer": "Record Company",
					"brand": "Your Favorite Band",
					"quantity": 666,
					"quantity_per_package": 1,
					"price": 20,
					"cost": 10
				},
				{
					"id": 2,
					"product_root_id": 1,
					"name": "Your Favorite Band's T-Shirt",
					"subtitle": "A t-shirt you can wear",
					"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					"option_summary": "Size: Medium, Color: Red",
					"sku": "t-shirt-medium-red",
					"manufacturer": "Record Company",
					"brand": "Your Favorite Band",
					"quantity": 666,
					"quantity_per_package": 1,
					"price": 20,
					"cost": 10
				},
				{
					"id": 3,
					"product_root_id": 1,
					"name": "Your Favorite Band's T-Shirt",
					"subtitle": "A t-shirt you can wear",
					"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					"option_summary": "Size: Large, Color: Red",
					"sku": "t-shirt-large-red",
					"manufacturer": "Record Company",
					"brand": "Your Favorite Band",
					"quantity": 666,
					"quantity_per_package": 1,
					"price": 20,
					"cost": 10
				},
				{
					"id": 4,
					"product_root_id": 1,
					"name": "Your Favorite Band's T-Shirt",
					"subtitle": "A t-shirt you can wear",
					"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					"option_summary": "Size: Small, Color: Blue",
					"sku": "t-shirt-small-blue",
					"manufacturer": "Record Company",
					"brand": "Your Favorite Band",
					"quantity": 666,
					"quantity_per_package": 1,
					"price": 20,
					"cost": 10
				},
				{
					"id": 5,
					"product_root_id": 1,
					"name": "Your Favorite Band's T-Shirt",
					"subtitle": "A t-shirt you can wear",
					"description": "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					"option_summary": "Size: Medium, Color: Blue",
					"sku": "t-shirt-medium-blue",
					"manufacturer": "Record Company",
					"brand": "Your Favorite Band",
					"quantity": 666,
					"quantity_per_package": 1,
					"price": 20,
					"cost": 10
				}
			]
		}
	`

	t.Run("normal usage", func(_t *testing.T) {
		expected := []models.Product{
			{
				ID:                 1,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Small, Color: Red",
				SKU:                "t-shirt-small-red",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 2,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Medium, Color: Red",
				SKU:                "t-shirt-medium-red",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 3,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Large, Color: Red",
				SKU:                "t-shirt-large-red",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 4,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Small, Color: Blue",
				SKU:                "t-shirt-small-blue",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				ID:                 5,
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Medium, Color: Blue",
				SKU:                "t-shirt-medium-blue",
				UPC:                "",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
		}

		handlers := map[string]http.HandlerFunc{
			"/v1/products": generateGetHandler(t, exampleGoodResponse, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		actual, err := c.GetProducts(nil)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected product doesn't match actual product")
	})

	t.Run("with bad server response", func(_t *testing.T) {
		handlers := map[string]http.HandlerFunc{
			"/v1/products": generateGetHandler(t, exampleBadJSON, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		_, err := c.GetProducts(nil)
		assert.NotNil(t, err, "GetProducts should return an error when it receives nonsense")
	})
}

func TestCreateProduct(t *testing.T) {
	exampleProductCreationInput := models.ProductCreationInput{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKU:                "sku",
		UPC:                "upc",
		Manufacturer:       "manufacturer",
		Brand:              "brand",
		Quantity:           666,
		Price:              20,
		SalePrice:          10,
		Cost:               1.23,
		ProductWeight:      9,
		ProductHeight:      9,
		ProductWidth:       9,
		ProductLength:      9,
		PackageWeight:      9,
		PackageHeight:      9,
		PackageWidth:       9,
		PackageLength:      9,
		QuantityPerPackage: 1,
	}

	t.Run("normal response", func(_t *testing.T) {
		var normalEndpointCalled bool

		handlers := map[string]http.HandlerFunc{
			"/v1/product": func(res http.ResponseWriter, req *http.Request) {
				normalEndpointCalled = true
				assert.Equal(t, req.Method, http.MethodPost, "CreateProduct should only be making POST requests")

				bodyBytes, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)

				expected := `
					{
						"product_width": 9,
						"package_length": 9,
						"sale_price": 10,
						"description": "description",
						"package_weight": 9,
						"price": 20,
						"product_weight": 9,
						"quantity": 666,
						"product_height": 9,
						"taxable": false,
						"brand": "brand",
						"product_length": 9,
						"available_on": "0001-01-01T00:00:00Z",
						"quantity_per_package": 1,
						"on_sale": false,
						"name": "name",
						"sku": "sku",
						"manufacturer": "manufacturer",
						"subtitle": "subtitle",
						"package_width": 9,
						"cost": 1.23,
						"package_height": 9,
						"option_summary": "",
						"updated_on": null,
						"upc": "upc",
						"options": null
					}
				`
				actual := string(bodyBytes)
				assert.Equal(t, minifyJSON(t, expected), actual, "CreateProduct should attach the correct JSON to the request body")

				exampleResponse := `
					{
						"name": "name",
						"subtitle": "subtitle",
						"description": "description",
						"option_summary": "option_summary",
						"sku": "sku",
						"upc": "upc",
						"manufacturer": "manufacturer",
						"brand": "brand",
						"quantity": 666,
						"quantity_per_package": 1,
						"taxable": false,
						"price": 20,
						"on_sale": false,
						"sale_price": 10,
						"cost": 1.23,
						"product_weight": 9,
						"product_height": 9,
						"product_width": 9,
						"product_length": 9,
						"package_weight": 9,
						"package_height": 9,
						"package_width": 9,
						"package_length": 9
					}
				`
				fmt.Fprintf(res, exampleResponse)
			},
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		expected := &models.Product{
			Name:               "name",
			Subtitle:           "subtitle",
			Description:        "description",
			OptionSummary:      "option_summary",
			SKU:                "sku",
			UPC:                "upc",
			Manufacturer:       "manufacturer",
			Brand:              "brand",
			Quantity:           666,
			Price:              20,
			SalePrice:          10,
			Cost:               1.23,
			ProductWeight:      9,
			ProductHeight:      9,
			ProductWidth:       9,
			ProductLength:      9,
			PackageWeight:      9,
			PackageHeight:      9,
			PackageWidth:       9,
			PackageLength:      9,
			QuantityPerPackage: 1,
		}
		actual, err := c.CreateProduct(exampleProductCreationInput)

		assert.Nil(t, err, "CreateProduct with valid input and response should never produce an error")
		assert.Equal(t, expected, actual, "expected and actual products should match")
		assert.True(t, normalEndpointCalled, "the normal endpoint should be called")
	})

	t.Run("with bad server response", func(_t *testing.T) {
		var badEndpointCalled bool
		handlers := map[string]http.HandlerFunc{
			"/v1/product": func(res http.ResponseWriter, req *http.Request) {
				badEndpointCalled = true
				fmt.Fprintf(res, exampleBadJSON)
			},
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		_, err := c.CreateProduct(exampleProductCreationInput)
		assert.NotNil(t, err, "CreateProduct should return an error when it fails to load a response")
		assert.True(t, badEndpointCalled, "the bad response endpoint should be called")
	})

	t.Run("with request error", func(_t *testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.CreateProduct(models.ProductCreationInput{})
		assert.NotNil(t, err, "CreateProduct should return an error when faililng to execute a request")
	})
}

// Note: this test is basically the same as TestCreateProduct, because those functions are incredibly similar, but with different purposes.
// I could probably sleep well at night with no tests for this, if only it wouldn't lower my precious coverage number.
func TestUpdateProduct(t *testing.T) {
	exampleProductUpdateInput := models.Product{
		Name:               "name",
		Subtitle:           "subtitle",
		Description:        "description",
		SKU:                "sku",
		UPC:                "upc",
		Manufacturer:       "manufacturer",
		Brand:              "brand",
		Quantity:           666,
		Price:              20,
		SalePrice:          10,
		Cost:               1.23,
		ProductWeight:      9,
		ProductHeight:      9,
		ProductWidth:       9,
		ProductLength:      9,
		PackageWeight:      9,
		PackageHeight:      9,
		PackageWidth:       9,
		PackageLength:      9,
		QuantityPerPackage: 1,
	}

	t.Run("normal response", func(_t *testing.T) {
		var normalEndpointCalled bool

		handlers := map[string]http.HandlerFunc{
			"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
				normalEndpointCalled = true
				assert.Equal(t, req.Method, http.MethodPatch, "UpdateProduct should only be making PATCH requests")

				bodyBytes, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)

				expected := `
					{
						"product_width": 9,
						"package_length": 9,
						"sale_price": 10,
						"description": "description",
						"package_weight": 9,
						"price": 20,
						"product_weight": 9,
						"quantity": 666,
						"product_root_id": 0,
						"product_height": 9,
						"taxable": false,
						"brand": "brand",
						"product_length": 9,
						"created_on": "0001-01-01T00:00:00Z",
						"available_on": "0001-01-01T00:00:00Z",
						"quantity_per_package": 1,
						"on_sale": false,
						"name": "name",
						"sku": "sku",
						"manufacturer": "manufacturer",
						"subtitle": "subtitle",
						"package_width": 9,
						"cost": 1.23,
						"id": 0,
						"package_height": 9,
						"archived_on": null,
						"option_summary": "",
						"updated_on": null,
						"upc": "upc"
					}
				`
				actual := string(bodyBytes)
				assert.Equal(t, minifyJSON(t, expected), actual, "UpdateProduct should attach the correct JSON to the request body")

				exampleResponse := `
					{
						"name": "name",
						"subtitle": "subtitle",
						"description": "description",
						"option_summary": "option_summary",
						"sku": "sku",
						"upc": "upc",
						"manufacturer": "manufacturer",
						"brand": "brand",
						"quantity": 666,
						"quantity_per_package": 1,
						"taxable": false,
						"price": 20,
						"on_sale": false,
						"sale_price": 10,
						"cost": 1.23,
						"product_weight": 9,
						"product_height": 9,
						"product_width": 9,
						"product_length": 9,
						"package_weight": 9,
						"package_height": 9,
						"package_width": 9,
						"package_length": 9
					}
				`
				fmt.Fprintf(res, exampleResponse)
			},
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		expected := &models.Product{
			Name:               "name",
			Subtitle:           "subtitle",
			Description:        "description",
			OptionSummary:      "option_summary",
			SKU:                "sku",
			UPC:                "upc",
			Manufacturer:       "manufacturer",
			Brand:              "brand",
			Quantity:           666,
			Price:              20,
			SalePrice:          10,
			Cost:               1.23,
			ProductWeight:      9,
			ProductHeight:      9,
			ProductWidth:       9,
			ProductLength:      9,
			PackageWeight:      9,
			PackageHeight:      9,
			PackageWidth:       9,
			PackageLength:      9,
			QuantityPerPackage: 1,
		}
		actual, err := c.UpdateProduct(exampleSKU, exampleProductUpdateInput)
		assert.Nil(t, err, "UpdateProduct with valid input and response should never produce an error")
		assert.Equal(t, expected, actual, "expected and actual products should match")
		assert.True(t, normalEndpointCalled, "the normal endpoint should be called")
	})

	t.Run("bad response", func(_t *testing.T) {
		var badEndpointCalled bool
		handlers := map[string]http.HandlerFunc{
			"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
				badEndpointCalled = true
				fmt.Fprintf(res, exampleBadJSON)
			},
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		_, err := c.UpdateProduct(exampleSKU, exampleProductUpdateInput)
		assert.NotNil(t, err, "UpdateProduct should return an error when it fails to load a response")
		assert.True(t, badEndpointCalled, "the bad response endpoint should be called")
	})

	t.Run("with request error", func(_t *testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.UpdateProduct(exampleSKU, models.Product{})
		assert.NotNil(t, err, "UpdateProduct should return an error when faililng to execute a request")
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Parallel()

	exampleResponseJSON := `
		{
			"id": 1,
			"product_root_id": 1,
			"name": "New Product",
			"subtitle": "this is a product",
			"description": "this product is neat or maybe its not who really knows for sure?",
			"option_summary": "",
			"sku": "test-product-updating",
			"upc": "",
			"manufacturer": "Manufacturer",
			"brand": "Brand",
			"quantity": 123,
			"taxable": false,
			"price": 12.34,
			"on_sale": true,
			"sale_price": 10,
			"cost": 5,
			"product_weight": 9,
			"product_height": 9,
			"product_width": 9,
			"product_length": 9,
			"package_weight": 9,
			"package_height": 9,
			"package_width": 9,
			"package_length": 9,
			"quantity_per_package": 3,
			"available_on": "0001-01-01T00:00:00Z",
			"created_on": "2017-12-10T06:03:54.394692Z",
			"updated_on": "",
			"archived_on": "2017-12-10T06:04:09.779255Z"
		}
	`

	existentSKU := "existent_sku"
	nonexistentSKU := "nonexistent_sku"

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", existentSKU):    generateDeleteHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU): generateDeleteHandler(t, buildNotFoundProductResponse(nonexistentSKU), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(_t *testing.T) {
		err := c.DeleteProduct(existentSKU)
		assert.Nil(t, err)
	})

	t.Run("with nonexistent product", func(_t *testing.T) {
		err := c.DeleteProduct(nonexistentSKU)
		assert.NotNil(t, err)
	})
}

func TestGetProductRoot(t *testing.T) {
	t.Skip()
}

func TestGetProductRoots(t *testing.T) {
	t.Skip()
}

func TestDeleteProductRoot(t *testing.T) {
	t.Skip()
}

func TestGetProductOptions(t *testing.T) {
	t.Skip()
}

func TestCreateProductOptionForProduct(t *testing.T) {
	t.Skip()
}

func TestUpdateProductOption(t *testing.T) {
	t.Skip()
}

func TestDeleteProductOption(t *testing.T) {
	t.Skip()
}

func TestCreateProductOptionValueForOption(t *testing.T) {
	t.Skip()
}

func TestUpdateProductOptionValueForOption(t *testing.T) {
	t.Skip()
}

func TestDeleteProductOptionValueForOption(t *testing.T) {
	t.Skip()
}
