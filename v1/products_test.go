package dairyclient_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	existentSKU := "existent_sku"
	nonexistentSKU := "nonexistent_sku"

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", existentSKU):    generateHeadHandler(t, http.StatusOK),
		fmt.Sprintf("/v1/product/%s", nonexistentSKU): generateHeadHandler(t, http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("with existent product", func(*testing.T) {
		exists, err := c.ProductExists(existentSKU)
		assert.Nil(t, err)
		assert.True(t, exists)

	})

	t.Run("with nonexistent product", func(*testing.T) {
		exists, err := c.ProductExists(nonexistentSKU)
		assert.Nil(t, err)
		assert.False(t, exists)
	})
}

func TestGetProduct(t *testing.T) {
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

	t.Run("normal usage", func(*testing.T) {
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

	t.Run("nonexistent product", func(*testing.T) {
		_, err := c.GetProduct(nonexistentSKU)
		assert.NotNil(t, err)
	})

	t.Run("bad response from server", func(*testing.T) {
		_, err := c.GetProduct(badResponseSKU)
		assert.NotNil(t, err)
	})

	t.Run("with request error", func(*testing.T) {
		ts.Close()
		_, err := c.GetProduct(exampleSKU)
		assert.NotNil(t, err)
	})
}

func TestGetProducts(t *testing.T) {
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

	t.Run("normal usage", func(*testing.T) {
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

	t.Run("with bad server response", func(*testing.T) {
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

	t.Run("normal response", func(*testing.T) {
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

	t.Run("with bad server response", func(*testing.T) {
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

	t.Run("with request error", func(*testing.T) {
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

	t.Run("normal response", func(*testing.T) {
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

	t.Run("bad response", func(*testing.T) {
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

	t.Run("with request error", func(*testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.UpdateProduct(exampleSKU, models.Product{})
		assert.NotNil(t, err, "UpdateProduct should return an error when faililng to execute a request")
	})
}

func TestDeleteProduct(t *testing.T) {
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

	t.Run("with existent product", func(*testing.T) {
		err := c.DeleteProduct(existentSKU)
		assert.Nil(t, err)
	})

	t.Run("with nonexistent product", func(*testing.T) {
		err := c.DeleteProduct(nonexistentSKU)
		assert.NotNil(t, err)
	})
}

func compareProductOptionValues(t *testing.T, expected, actual models.ProductOptionValue, optionIndex, optionValueIndex int) {
	assert.Equal(t, expected.Value, actual.Value, "expected and actual Value for option value %d (option %d) should match", optionValueIndex, optionIndex)
	assert.Equal(t, expected.CreatedOn, actual.CreatedOn, "expected and actual CreatedOn for option value %d (option %d) should match", optionValueIndex, optionIndex)
	assert.Equal(t, expected.ID, actual.ID, "expected and actual ID for option value %d (option %d) should match", optionValueIndex, optionIndex)
	assert.Equal(t, expected.ArchivedOn, actual.ArchivedOn, "expected and actual ArchivedOn for option value %d (option %d) should match", optionValueIndex, optionIndex)
	assert.Equal(t, expected.UpdatedOn, actual.UpdatedOn, "expected and actual UpdatedOn for option value %d (option %d) should match", optionValueIndex, optionIndex)
	assert.Equal(t, expected.ProductOptionID, actual.ProductOptionID, "expected and actual ProductOptionID for option value %d (option %d) should match", optionValueIndex, optionIndex)
}

func compareProductOptions(t *testing.T, expected, actual models.ProductOption, optionIndex int) {
	assert.Equal(t, expected.ProductRootID, actual.ProductRootID, "expected and actual ProductRootID for option %d should match", optionIndex)
	assert.Equal(t, expected.CreatedOn, actual.CreatedOn, "expected and actual CreatedOn for option %d should match", optionIndex)
	assert.Equal(t, expected.ID, actual.ID, "expected and actual ID for option %d should match", optionIndex)
	assert.Equal(t, expected.ArchivedOn, actual.ArchivedOn, "expected and actual ArchivedOn for option %d should match", optionIndex)
	assert.Equal(t, expected.UpdatedOn, actual.UpdatedOn, "expected and actual UpdatedOn for option %d should match", optionIndex)
	assert.Equal(t, expected.Name, actual.Name, "expected and actual Name for option %d should match", optionIndex)

	for i := range expected.Values {
		if len(actual.Values)-1 < i {
			t.Logf("expected %d option values, got %d instead.", len(expected.Values), len(actual.Values))
			t.Fail()
			break
		}
		compareProductOptionValues(t, expected.Values[i], actual.Values[i], optionIndex, i)
	}
}

func compareProductRoots(t *testing.T, expected, actual *models.ProductRoot) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID, "expected and actual ID should match")
	assert.Equal(t, expected.AvailableOn, actual.AvailableOn, "expected and actual AvailableOn should match")
	assert.Equal(t, expected.ProductLength, actual.ProductLength, "expected and actual ProductLength should match")
	assert.Equal(t, expected.UpdatedOn, actual.UpdatedOn, "expected and actual UpdatedOn should match")
	assert.Equal(t, expected.SKUPrefix, actual.SKUPrefix, "expected and actual SKUPrefix should match")
	assert.Equal(t, expected.PackageHeight, actual.PackageHeight, "expected and actual PackageHeight should match")
	assert.Equal(t, expected.ProductWeight, actual.ProductWeight, "expected and actual ProductWeight should match")
	assert.Equal(t, expected.ProductWidth, actual.ProductWidth, "expected and actual ProductWidth should match")
	assert.Equal(t, expected.QuantityPerPackage, actual.QuantityPerPackage, "expected and actual QuantityPerPackage should match")
	assert.Equal(t, expected.Name, actual.Name, "expected and actual Name should match")
	assert.Equal(t, expected.ProductHeight, actual.ProductHeight, "expected and actual ProductHeight should match")
	assert.Equal(t, expected.PackageLength, actual.PackageLength, "expected and actual PackageLength should match")
	assert.Equal(t, expected.CreatedOn, actual.CreatedOn, "expected and actual CreatedOn should match")
	assert.Equal(t, expected.Cost, actual.Cost, "expected and actual Cost should match")
	assert.Equal(t, expected.Brand, actual.Brand, "expected and actual Brand should match")
	assert.Equal(t, expected.Subtitle, actual.Subtitle, "expected and actual Subtitle should match")
	assert.Equal(t, expected.PackageWeight, actual.PackageWeight, "expected and actual PackageWeight should match")
	assert.Equal(t, expected.ArchivedOn, actual.ArchivedOn, "expected and actual ArchivedOn should match")
	assert.Equal(t, expected.PackageWidth, actual.PackageWidth, "expected and actual PackageWidth should match")
	assert.Equal(t, expected.Description, actual.Description, "expected and actual Description should match")
	assert.Equal(t, expected.Manufacturer, actual.Manufacturer, "expected and actual Manufacturer should match")
	assert.Equal(t, expected.Taxable, actual.Taxable, "expected and actual Taxable should match")

	for i := range expected.Options {
		if len(actual.Options)-1 < i {
			t.Logf("expected %d options, got %d instead.", len(expected.Options), len(actual.Options))
			t.Fail()
			break
		}
		compareProductOptions(t, expected.Options[i], actual.Options[i], i)
	}

	for i := range expected.Products {
		if len(actual.Products)-1 < i {
			t.Logf("expected %d products, got %d instead.", len(expected.Products), len(actual.Products))
			t.Fail()
			break
		}
		compareProducts(t, expected.Products[i], actual.Products[i])
	}
}

// TODO: maybe these functions should just set the values that we don't care about equality for rather than check for the equality of each field
// for instance, we don't really worry about IDs, so make this function set the expected.ID to actual.ID and then use assert to check equality
func compareProducts(t *testing.T, expected models.Product, actual models.Product) {
	t.Helper()
	assert.Equal(t, expected.ProductWidth, actual.ProductWidth, "expected and actual ProductWidth should match")
	assert.Equal(t, expected.PackageLength, actual.PackageLength, "expected and actual PackageLength should match")
	assert.Equal(t, expected.SalePrice, actual.SalePrice, "expected and actual SalePrice should match")
	assert.Equal(t, expected.Description, actual.Description, "expected and actual Description should match")
	assert.Equal(t, expected.PackageWeight, actual.PackageWeight, "expected and actual PackageWeight should match")
	assert.Equal(t, expected.Price, actual.Price, "expected and actual Price should match")
	assert.Equal(t, expected.ProductWeight, actual.ProductWeight, "expected and actual ProductWeight should match")
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected and actual Quantity should match")
	assert.Equal(t, expected.ProductRootID, actual.ProductRootID, "expected and actual ProductRootID should match")
	assert.Equal(t, expected.ProductHeight, actual.ProductHeight, "expected and actual ProductHeight should match")
	assert.Equal(t, expected.Taxable, actual.Taxable, "expected and actual Taxable should match")
	assert.Equal(t, expected.Brand, actual.Brand, "expected and actual Brand should match")
	assert.Equal(t, expected.ProductLength, actual.ProductLength, "expected and actual ProductLength should match")
	assert.Equal(t, expected.CreatedOn, actual.CreatedOn, "expected and actual CreatedOn should match")
	assert.Equal(t, expected.AvailableOn, actual.AvailableOn, "expected and actual AvailableOn should match")
	assert.Equal(t, expected.QuantityPerPackage, actual.QuantityPerPackage, "expected and actual QuantityPerPackage should match")
	assert.Equal(t, expected.OnSale, actual.OnSale, "expected and actual OnSale should match")
	assert.Equal(t, expected.Name, actual.Name, "expected and actual Name should match")
	assert.Equal(t, expected.SKU, actual.SKU, "expected and actual SKU should match")
	assert.Equal(t, expected.Manufacturer, actual.Manufacturer, "expected and actual Manufacturer should match")
	assert.Equal(t, expected.Subtitle, actual.Subtitle, "expected and actual Subtitle should match")
	assert.Equal(t, expected.PackageWidth, actual.PackageWidth, "expected and actual PackageWidth should match")
	assert.Equal(t, expected.Cost, actual.Cost, "expected and actual Cost should match")
	assert.Equal(t, expected.ID, actual.ID, "expected and actual ID should match")
	assert.Equal(t, expected.PackageHeight, actual.PackageHeight, "expected and actual PackageHeight should match")
	assert.Equal(t, expected.ArchivedOn, actual.ArchivedOn, "expected and actual ArchivedOn should match")
	assert.Equal(t, expected.OptionSummary, actual.OptionSummary, "expected and actual OptionSummary should match")
	assert.Equal(t, expected.UpdatedOn, actual.UpdatedOn, "expected and actual UpdatedOn should match")
	assert.Equal(t, expected.UPC, actual.UPC, "expected and actual UPC should match")

	for i := range expected.ApplicableOptionValues {
		if len(actual.ApplicableOptionValues)-1 < i {
			t.Logf("expected %d option values attached to product, got %d instead.", len(expected.ApplicableOptionValues), len(actual.ApplicableOptionValues))
			t.Fail()
			break
		}
		compareProductOptionValues(t, expected.ApplicableOptionValues[i], actual.ApplicableOptionValues[i], 0, i)
	}
}

func buildNotFoundProductRootResponse(id uint64) string {
	return fmt.Sprintf(`
		{
			"status": 404,
			"message": "The product_root you were looking for (identified by '%d') does not exist"
		}
	`, id)
}

func TestGetProductRoot(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "product_root")
	existentID := uint64(1)
	nonexistentID := uint64(2)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_root/%d", existentID):    generateGetHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product_root/%d", nonexistentID): generateGetHandler(t, buildNotFoundProductRootResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		pTime, err := time.Parse(timeLayout, "2017-12-10T15:58:43.136458Z")
		expected := &models.ProductRoot{
			ID:                 1,
			Name:               "Your Favorite Band's T-Shirt",
			Subtitle:           "A t-shirt you can wear",
			Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			SKUPrefix:          "t-shirt",
			Manufacturer:       "Record Company",
			Brand:              "Your Favorite Band",
			Taxable:            true,
			Cost:               20,
			ProductWeight:      1,
			ProductHeight:      5,
			ProductWidth:       5,
			ProductLength:      5,
			PackageWeight:      1,
			PackageHeight:      5,
			PackageWidth:       5,
			PackageLength:      5,
			QuantityPerPackage: 1,
			AvailableOn:        pTime,
			CreatedOn:          pTime,
			Options: []models.ProductOption{
				{
					ID:            1,
					ProductRootID: 1,
					Name:          "color",
					CreatedOn:     pTime,
				},
				{
					ID:            2,
					ProductRootID: 1,
					Name:          "size",
					CreatedOn:     pTime,
				},
			},
			Products: []models.Product{
				{
					ID:                 1,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-small-red",
					OptionSummary:      "Size: Small, Color: Red",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 2,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-medium-red",
					OptionSummary:      "Size: Medium, Color: Red",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 3,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-large-red",
					OptionSummary:      "Size: Large, Color: Red",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 4,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-small-blue",
					OptionSummary:      "Size: Small, Color: Blue",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 5,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-medium-blue",
					OptionSummary:      "Size: Medium, Color: Blue",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 6,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-large-blue",
					OptionSummary:      "Size: Large, Color: Blue",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 7,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-small-green",
					OptionSummary:      "Size: Small, Color: Green",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 8,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-medium-green",
					OptionSummary:      "Size: Medium, Color: Green",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
				{
					ID:                 9,
					ProductRootID:      1,
					Name:               "Your Favorite Band's T-Shirt",
					Subtitle:           "A t-shirt you can wear",
					Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
					SKU:                "t-shirt-large-green",
					OptionSummary:      "Size: Large, Color: Green",
					Manufacturer:       "Record Company",
					Brand:              "Your Favorite Band",
					Taxable:            true,
					Quantity:           666,
					Price:              20,
					Cost:               10,
					ProductWeight:      1,
					ProductHeight:      5,
					ProductWidth:       5,
					ProductLength:      5,
					PackageWeight:      1,
					PackageHeight:      5,
					PackageWidth:       5,
					PackageLength:      5,
					QuantityPerPackage: 1,
					AvailableOn:        pTime,
					CreatedOn:          pTime,
				},
			},
		}

		actual, err := c.GetProductRoot(existentID)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "expected and actual product roots don't match.")
	})

	t.Run("with nonexistent product root", func(*testing.T) {
		_, err := c.GetProductRoot(nonexistentID)
		assert.NotNil(t, err)
	})
}

func TestGetProductRoots(t *testing.T) {
	exampleResponseJSON := loadExampleResponse(t, "product_roots")
	handlers := map[string]http.HandlerFunc{
		"/v1/product_roots": generateGetHandler(t, exampleResponseJSON, http.StatusOK),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	pTime, err := time.Parse(timeLayout, "2017-12-10T15:58:43.136458Z")
	assert.Nil(t, err)

	t.Run("normal usage", func(*testing.T) {
		expected := []models.ProductRoot{
			{
				ID:                 5,
				Name:               "Animals As Leaders - The Joy Of Motion",
				Subtitle:           "A solid prog metal album",
				Description:        "Arbitrary description can go here because real product descriptions are technically copywritten.",
				Brand:              "Animals As Leaders",
				Manufacturer:       "Record Company",
				SKUPrefix:          "the-joy-of-motion",
				QuantityPerPackage: 1,
				ProductLength:      0.5,
				PackageHeight:      12,
				ProductWeight:      1,
				ProductWidth:       12,
				ProductHeight:      12,
				PackageLength:      0.5,
				PackageWeight:      1,
				PackageWidth:       12,
				Cost:               5,
				Taxable:            true,
				AvailableOn:        pTime,
				CreatedOn:          pTime,
			},
			{
				ID:                 6,
				Name:               "Mort Garson - Mother Earth's Plantasia",
				Subtitle:           "A solid synth album",
				Description:        "Arbitrary description can go here because real product descriptions are technically copywritten.",
				Brand:              "Mort Garson",
				Manufacturer:       "Record Company",
				SKUPrefix:          "mother-earths-plantasia",
				QuantityPerPackage: 1,
				ProductLength:      0.5,
				PackageHeight:      12,
				ProductWeight:      1,
				ProductWidth:       12,
				ProductHeight:      12,
				PackageLength:      0.5,
				PackageWeight:      1,
				PackageWidth:       12,
				Cost:               5,
				Taxable:            true,
				AvailableOn:        pTime,
				CreatedOn:          pTime,
			},
		}

		actual, err := c.GetProductRoots(nil)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestDeleteProductRoot(t *testing.T) {
	existentID := uint64(1)
	nonexistentID := uint64(2)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product_root/%d", existentID):    generateDeleteHandler(t, "{}", http.StatusNotFound),
		fmt.Sprintf("/v1/product_root/%d", nonexistentID): generateDeleteHandler(t, buildNotFoundProductRootResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		err := c.DeleteProductRoot(existentID)
		assert.Nil(t, err)
	})

	t.Run("nonexistent product root", func(*testing.T) {
		err := c.DeleteProductRoot(nonexistentID)
		assert.NotNil(t, err)
	})
}

func buildNotFoundProductOptionsResponse(productID uint64) string {
	// FIXME
	return fmt.Sprintf(`
		{
			"count": 2,
			"limit": 25,
			"page": 1,
			"data": null
		}
	`, productID)
}

func TestGetProductOptions(t *testing.T) {
	existentID := uint64(1)
	nonexistentID := uint64(2)
	exampleResponseJSON := loadExampleResponse(t, "product_options")

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%d/options", existentID):    generateGetHandler(t, exampleResponseJSON, http.StatusOK),
		fmt.Sprintf("/v1/product/%d/options", nonexistentID): generateGetHandler(t, buildNotFoundProductOptionsResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	pTime, err := time.Parse(timeLayout, "2017-12-10T15:58:43.136458Z")
	assert.Nil(t, err)

	t.Run("normal operation", func(*testing.T) {
		expected := []models.ProductOption{
			{
				ID:            1,
				Name:          "color",
				ProductRootID: 1,
				CreatedOn:     pTime,
				Values: []models.ProductOptionValue{
					{
						ID:              1,
						ProductOptionID: 1,
						Value:           "red",
						CreatedOn:       pTime,
					},
					{
						ID:              2,
						ProductOptionID: 1,
						Value:           "green",
						CreatedOn:       pTime,
					},
					{
						ID:              3,
						ProductOptionID: 1,
						Value:           "blue",
						CreatedOn:       pTime,
					},
				},
			},
			{
				ID:            2,
				Name:          "size",
				ProductRootID: 1,
				CreatedOn:     pTime,
				Values: []models.ProductOptionValue{
					{
						ID:              4,
						ProductOptionID: 2,
						Value:           "small",
						CreatedOn:       pTime,
					},
					{
						ID:              5,
						ProductOptionID: 2,
						Value:           "medium",
						CreatedOn:       pTime,
					},
					{
						ID:              6,
						ProductOptionID: 2,
						Value:           "large",
						CreatedOn:       pTime,
					},
				},
			},
		}

		actual, err := c.GetProductOptions(existentID, nil)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("for nonexistent product root", func(*testing.T) {
		_, err := c.GetProductOptions(nonexistentID, nil)
		assert.NotNil(t, err)
	})
}

func TestCreateProductOption(t *testing.T) {
	existentID := uint64(1)
	exampleResponseJSON := loadExampleResponse(t, "created_product_options")
	expectedBody := `
		{
			"name": "example_option",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%d/options", existentID): generatePostHandler(t, expectedBody, exampleResponseJSON, http.StatusCreated),
		// fmt.Sprintf("/v1/product/%d/options", nonexistentID): generateGetHandler(t, buildNotFoundProductOptionsResponse(nonexistentID), http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	pTime, err := time.Parse(timeLayout, "2017-12-10T15:58:43.136458Z")
	assert.Nil(t, err)

	t.Run("normal operation", func(*testing.T) {
		exampleInput := models.ProductOptionCreationInput{
			Name:   "example_option",
			Values: []string{"one", "two", "three"},
		}

		expected := &models.ProductOption{
			ID:            3,
			Name:          "example_option",
			ProductRootID: 1,
			CreatedOn:     pTime,
			Values: []models.ProductOptionValue{
				{
					ID:              7,
					ProductOptionID: 3,
					Value:           "one",
					CreatedOn:       pTime,
				},
				{
					ID:              8,
					ProductOptionID: 3,
					Value:           "two",
					CreatedOn:       pTime,
				},
				{
					ID:              9,
					ProductOptionID: 3,
					Value:           "three",
					CreatedOn:       pTime,
				},
			},
		}

		actual, err := c.CreateProductOption(1, exampleInput)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

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
