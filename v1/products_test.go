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

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	exists, err := c.ProductExists(existentSKU)
	assert.Nil(t, err)
	assert.True(t, exists)

	exists, err = c.ProductExists(nonexistentSKU)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

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
	`, exampleSKU)

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", exampleSKU): generateGetHandler(t, exampleResponse, http.StatusOK),
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	expected := &dairyclient.Product{
		Name:               "Your Favorite Band's T-Shirt",
		Subtitle:           "A t-shirt you can wear",
		Description:        "Wear this if you'd like. Or don't, I'm not in charge of your actions",
		OptionSummary:      "Size: Small, Color: Red",
		SKU:                exampleSKU,
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
	actual, err := c.GetProduct(exampleSKU)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "expected product doesn't match actual product")
}

func TestGetProductReturnsErrorWhenExecutingRequestFails(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := buildTestClient(t, ts)
	ts.Close()

	_, err := c.GetProduct(exampleSKU)

	assert.NotNil(t, err)
}

func TestGetProductReturnsErrorWhenReceivingBadJSON(t *testing.T) {
	t.Parallel()
	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", exampleSKU): generateGetHandler(t, exampleBadJSON, http.StatusOK),
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	_, err := c.GetProduct(exampleSKU)

	assert.NotNil(t, err)
}
