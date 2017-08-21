package dairyclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	currentAPIVersion = `v1`
)

type V1Client struct {
	*http.Client
	URL        *url.URL
	AuthCookie *http.Cookie
}

func NewV1Client(storeURL string, username string, password string, client *http.Client) (*V1Client, error) {
	var dc *V1Client
	if client != nil {
		dc = &V1Client{Client: client}
	}

	u, err := url.Parse(storeURL)
	if err != nil {
		return nil, errors.Wrap(err, "Store URL is not valid")
	}
	dc.URL = u

	p := fmt.Sprintf("%s://%s/login", u.Scheme, u.Host)
	body := strings.NewReader(fmt.Sprintf(`
		{
			"username": "%s",
			"password": "%s"
		}
	`, username, password))
	req, _ := http.NewRequest(http.MethodPost, p, body)
	res, err := dc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered logging into store")
	}
	cookies := res.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("No cookies returned with login response")
	}

	for _, c := range cookies {
		if c.Name == "dairycart" {
			dc.AuthCookie = c
		}
	}
	dc.Client.Timeout = 5 * time.Second

	return dc, nil
}

func (dc *V1Client) executeRequest(req *http.Request) (*http.Response, error) {
	req.AddCookie(dc.AuthCookie)
	return dc.Do(req)
}

func (dc *V1Client) buildURL(queryParams map[string]string, parts ...string) string {
	parts = append([]string{currentAPIVersion}, parts...)
	u, _ := url.Parse(strings.Join(parts, "/"))
	queryString := mapToQueryValues(queryParams)
	u.RawQuery = queryString.Encode()
	return dc.URL.ResolveReference(u).String()
}

// BuildURL is the same as the unexported build URL, except I trust myself to never call the
// unexported function with variables that could lead to an error being returned. This function
// returns the error in the event a user needs to build an API url, but tries to do so with an
// invalid value.
func (dc *V1Client) BuildURL(queryParams map[string]string, parts ...string) (string, error) {
	parts = append([]string{currentAPIVersion}, parts...)

	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		return "", err
	}

	queryString := mapToQueryValues(queryParams)
	u.RawQuery = queryString.Encode()
	return dc.URL.ResolveReference(u).String(), nil
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) ProductExists(sku string) (bool, error) {
	u := dc.buildURL(nil, "product", sku)
	req, _ := http.NewRequest(http.MethodHead, u, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return false, err
	}

	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (dc *V1Client) GetProduct(sku string) (*Product, error) {
	u := dc.buildURL(nil, "product", sku)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	p := Product{}
	err = unmarshalBody(res, p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) GetProducts(queryFilter map[string]string) (*http.Response, error) {
	u := dc.buildURL(queryFilter, "products")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) CreateProduct(np ProductCreationInput) (*http.Response, error) {
	body, err := createBodyFromStruct(np)
	if err != nil {
		return nil, err
	}

	u := dc.buildURL(nil, "product")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) UpdateProduct(sku string, up ProductUpdateInput) (*http.Response, error) {
	body, err := createBodyFromStruct(up)
	if err != nil {
		return nil, err
	}

	u := dc.buildURL(nil, "product", sku)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProduct(sku string) (*http.Response, error) {
	u := dc.buildURL(nil, "product", sku)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//              Product Root Functions                //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetProductRoot(rootID uint64) (*http.Response, error) {
	rootIDString := convertIDToString(rootID)
	u := dc.buildURL(nil, "product_root", rootIDString)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) GetProductRoots(queryFilter map[string]string) (*http.Response, error) {
	u := dc.buildURL(queryFilter, "product_roots")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProductRoot(rootID uint64) (*http.Response, error) {
	rootIDString := convertIDToString(rootID)
	u := dc.buildURL(nil, "product_root", rootIDString)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//             Product Option Functions               //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetProductOptions(productID uint64, queryFilter map[string]string) (*http.Response, error) {
	productIDString := convertIDToString(productID)
	u := dc.buildURL(queryFilter, "product", productIDString, "options")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) CreateProductOptionForProduct(productID uint64, JSONBody string) (*http.Response, error) {
	productIDString := convertIDToString(productID)
	body := strings.NewReader(JSONBody)
	u := dc.buildURL(nil, "product", productIDString, "options")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) UpdateProductOption(optionID uint64, JSONBody string) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	body := strings.NewReader(JSONBody)
	u := dc.buildURL(nil, "product_options", optionIDString)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProductOption(optionID uint64) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//          Product Option Value Functions            //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) createProductOptionValueForOption(optionID uint64, JSONBody string) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	body := strings.NewReader(JSONBody)
	u := dc.buildURL(nil, "product_options", optionIDString, "value")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) updateProductOptionValueForOption(valueID uint64, JSONBody string) (*http.Response, error) {
	valueIDString := convertIDToString(valueID)
	body := strings.NewReader(JSONBody)
	u := dc.buildURL(nil, "product_option_values", valueIDString)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) deleteProductOptionValueForOption(optionID uint64) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_option_values", optionIDString)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) getDiscountByID(discountID uint64) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) getListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	u := dc.buildURL(queryFilter, "discounts")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) createDiscount(JSONBody string) (*http.Response, error) {
	u := dc.buildURL(nil, "discount")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) updateDiscount(discountID uint64, JSONBody string) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) deleteDiscount(discountID uint64) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}
