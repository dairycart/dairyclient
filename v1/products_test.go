package dairyclient_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dairycart/dairyclient/v1"
)

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

	testExistsWithExistentProduct := func(t *testing.T) {
		exists, err := c.ProductExists(existentSKU)
		assert.Nil(t, err)
		assert.True(t, exists)
	}

	testExistsWithNonexistentProduct := func(t *testing.T) {
		exists, err := c.ProductExists(nonexistentSKU)
		assert.Nil(t, err)
		assert.False(t, exists)
	}

	subtests := []subtest{
		{
			Message: "existent product",
			Test:    testExistsWithExistentProduct,
		},
		{
			Message: "nonexistent product",
			Test:    testExistsWithNonexistentProduct,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	goodResponseSKU := "good"
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
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	normalResponse := func(t *testing.T) {
		expected := &dairyclient.Product{
			Name:               "Your Favorite Band's T-Shirt",
			Subtitle:           "A t-shirt you can wear",
			Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
			OptionSummary:      "Size: Small, Color: Red",
			SKU:                goodResponseSKU,
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
	}

	badResponse := func(t *testing.T) {
		_, err := c.GetProduct(badResponseSKU)
		assert.NotNil(t, err)
	}

	requestError := func(t *testing.T) {
		ts.Close()
		_, err := c.GetProduct(exampleSKU)
		assert.NotNil(t, err)
	}

	subtests := []subtest{
		{
			Message: "normal response",
			Test:    normalResponse,
		},
		{
			Message: "bad response",
			Test:    badResponse,
		},
		{
			Message: "error executing request",
			Test:    requestError,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestGetProducts(t *testing.T) {
	t.Parallel()

	exampleGoodResponse := `
		{
			"count": 5,
			"limit": 25,
			"page": 1,
			"data": [{
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

	normalResponse := func(t *testing.T) {
		expected := []dairyclient.Product{
			{
				DBRow: dairyclient.DBRow{
					ID: 1,
				},
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Small, Color: Red",
				SKU:                "t-shirt-small-red",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				DBRow: dairyclient.DBRow{
					ID: 2,
				},
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Medium, Color: Red",
				SKU:                "t-shirt-medium-red",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				DBRow: dairyclient.DBRow{
					ID: 3,
				},
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Large, Color: Red",
				SKU:                "t-shirt-large-red",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				DBRow: dairyclient.DBRow{
					ID: 4,
				},
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Small, Color: Blue",
				SKU:                "t-shirt-small-blue",
				Manufacturer:       "Record Company",
				Brand:              "Your Favorite Band",
				Quantity:           666,
				QuantityPerPackage: 1,
				Price:              20,
				Cost:               10,
			},
			{
				DBRow: dairyclient.DBRow{
					ID: 5,
				},
				ProductRootID:      1,
				Name:               "Your Favorite Band's T-Shirt",
				Subtitle:           "A t-shirt you can wear",
				Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
				OptionSummary:      "Size: Medium, Color: Blue",
				SKU:                "t-shirt-medium-blue",
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
	}

	badResponse := func(t *testing.T) {
		handlers := map[string]http.HandlerFunc{
			"/v1/products": generateGetHandler(t, exampleBadJSON, http.StatusOK),
		}
		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		c := buildTestClient(t, ts)

		_, err := c.GetProducts(nil)
		assert.NotNil(t, err, "GetProducts should return an error when it receives nonsense")
	}

	subtests := []subtest{
		{
			Message: "normal response",
			Test:    normalResponse,
		},
		{
			Message: "bad response",
			Test:    badResponse,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestCreateProduct(t *testing.T) {
	exampleProductCreationInput := dairyclient.ProductInput{
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

	normalResponse := func(t *testing.T) {
		var normalEndpointCalled bool

		handlers := map[string]http.HandlerFunc{
			"/v1/product": func(res http.ResponseWriter, req *http.Request) {
				normalEndpointCalled = true
				assert.Equal(t, req.Method, http.MethodPost, "CreateProduct should only be making POST requests")

				bodyBytes, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)

				expected := `
					{
						"name": "name",
						"subtitle": "subtitle",
						"description": "description",
						"sku": "sku",
						"upc": "upc",
						"manufacturer": "manufacturer",
						"brand": "brand",
						"quantity": 666,
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
						"package_length": 9,
						"quantity_per_package": 1,
						"available_on": "0001-01-01T00:00:00Z",
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

		expected := &dairyclient.Product{
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
	}

	badResponse := func(t *testing.T) {
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
	}

	requestError := func(t *testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.CreateProduct(dairyclient.ProductInput{})
		assert.NotNil(t, err, "CreateProduct should return an error when faililng to execute a request")
	}

	subtests := []subtest{
		{
			Message: "normal response",
			Test:    normalResponse,
		},
		{
			Message: "bad response",
			Test:    badResponse,
		},
		{
			Message: "error executing request",
			Test:    requestError,
		},
	}
	runSubtestSuite(t, subtests)
}

// Note: this test is basically the same as TestCreateProduct, because those functions are incredibly similar, but with different purposes.
// I could probably sleep well at night with no tests for this, if only it wouldn't lower my precious coverage number.
func TestUpdateProduct(t *testing.T) {
	exampleProductUpdateInput := dairyclient.ProductInput{
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

	normalResponse := func(t *testing.T) {
		var normalEndpointCalled bool

		handlers := map[string]http.HandlerFunc{
			"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
				normalEndpointCalled = true
				assert.Equal(t, req.Method, http.MethodPatch, "UpdateProduct should only be making PATCH requests")

				bodyBytes, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)

				expected := `
					{
						"name": "name",
						"subtitle": "subtitle",
						"description": "description",
						"sku": "sku",
						"upc": "upc",
						"manufacturer": "manufacturer",
						"brand": "brand",
						"quantity": 666,
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
						"package_length": 9,
						"quantity_per_package": 1,
						"available_on": "0001-01-01T00:00:00Z",
						"options": null
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

		expected := &dairyclient.Product{
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
	}

	badResponse := func(t *testing.T) {
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
	}

	requestError := func(t *testing.T) {
		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()
		_, err := c.UpdateProduct(exampleSKU, dairyclient.ProductInput{})
		assert.NotNil(t, err, "UpdateProduct should return an error when faililng to execute a request")
	}

	subtests := []subtest{
		{
			Message: "normal response",
			Test:    normalResponse,
		},
		{
			Message: "bad response",
			Test:    badResponse,
		},
		{
			Message: "error executing request",
			Test:    requestError,
		},
	}
	runSubtestSuite(t, subtests)
}
