package dairyclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	currentAPIVersion = `v1`
)

type DairyclientV1 struct {
	*http.Client
	URL        *url.URL
	AuthCookie *http.Cookie
}

func NewV1Client(storeURL string, username string, password string, client *http.Client) (*DairyclientV1, error) {
	var dc *DairyclientV1
	if client != nil {
		dc = &DairyclientV1{Client: client}
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

	return dc, nil
}

func (dc *DairyclientV1) executeRequest(req *http.Request) (*http.Response, error) {
	req.AddCookie(dc.AuthCookie)
	return dc.Do(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                 Helper Functions                   //
//                                                    //
////////////////////////////////////////////////////////

func mapToQueryValues(in map[string]string) url.Values {
	out := url.Values{}
	for k, v := range in {
		out.Set(k, v)
	}
	return out
}

func (dc *DairyclientV1) buildURL(queryParams map[string]string, versioned bool, parts ...string) (string, error) {
	if versioned {
		parts = append([]string{currentAPIVersion}, parts...)
	}

	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		return "", err
	}

	queryString := mapToQueryValues(queryParams)
	u.RawQuery = queryString.Encode()
	return dc.URL.ResolveReference(u).String(), nil
}

func (dc *DairyclientV1) buildVersionlessURL(queryParams map[string]string, parts ...string) (string, error) {
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
//                  Auth Functions                    //
//                                                    //
////////////////////////////////////////////////////////

func (dc *DairyclientV1) createNewUser(JSONBody string, createAsSuperUser bool) (*http.Response, error) {
	u, err := dc.buildURL(nil, false, "user")
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	if createAsSuperUser {
		return dc.executeRequest(req)
	}
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) deleteUser(userID string, deleteAsSuperUser bool) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "user", userID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	if deleteAsSuperUser {
		return dc.executeRequest(req)
	}
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *DairyclientV1) CheckProductExistence(sku string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product", sku)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodHead, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) retrieveProduct(sku string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product", sku)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) retrieveListOfProducts(queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.buildURL(queryFilter, true, "products")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) createProduct(JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u, err := dc.buildURL(nil, true, "product")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u, err := dc.buildURL(nil, true, "product", sku)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) deleteProduct(sku string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product", sku)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//              Product Root Functions                //
//                                                    //
////////////////////////////////////////////////////////

func (dc *DairyclientV1) retrieveProductRoot(rootID string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product_root", rootID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) retrieveProductRoots(queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.buildURL(queryFilter, true, "product_roots")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) deleteProductRoot(rootID string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product_root", rootID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//             Product Option Functions               //
//                                                    //
////////////////////////////////////////////////////////

func (dc *DairyclientV1) retrieveProductOptions(productID string, queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.buildURL(queryFilter, true, "product", productID, "options")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) createProductOptionForProduct(productID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u, err := dc.buildURL(nil, true, "product", productID, "options")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) updateProductOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u, err := dc.buildURL(nil, true, "product_options", optionID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) deleteProductOption(optionID string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product_options", optionID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//          Product Option Value Functions            //
//                                                    //
////////////////////////////////////////////////////////

func (dc *DairyclientV1) createProductOptionValueForOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u, err := dc.buildURL(nil, true, "product_options", optionID, "value")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) updateProductOptionValueForOption(valueID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u, err := dc.buildURL(nil, true, "product_option_values", valueID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) deleteProductOptionValueForOption(optionID string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "product_option_values", optionID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *DairyclientV1) getDiscountByID(discountID string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "discount", discountID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) getListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.buildURL(queryFilter, true, "discounts")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) createDiscount(JSONBody string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "discount")
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) updateDiscount(discountID string, JSONBody string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "discount", discountID)
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *DairyclientV1) deleteDiscount(discountID string) (*http.Response, error) {
	u, err := dc.buildURL(nil, true, "discount", discountID)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}
