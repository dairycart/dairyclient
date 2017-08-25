package dairyclient_test

import (
	"fmt"
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
