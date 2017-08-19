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

	return dc, nil
}

func (dc *V1Client) executeRequest(req *http.Request) (*http.Response, error) {
	req.AddCookie(dc.AuthCookie)
	return dc.Do(req)
}

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
//                  Auth Functions                    //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) CreateUser(JSONBody string, createAsSuperUser bool) (*User, error) {
	u, err := dc.BuildURL(nil, "user")
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)

	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	ru := User{}
	err = unmarshalBody(res.Body, ru)
	if err != nil {
		return nil, err
	}

	return &ru, nil
}

func (dc *V1Client) DeleteUser(userID uint64) error {
	userIDString := convertIDToString(userID)
	u, err := dc.BuildURL(nil, "user", userIDString)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("user couldn't be deleted, status returned: %d", res.StatusCode))
	}
	return nil
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) ProductExists(sku string) (bool, error) {
	u, err := dc.BuildURL(nil, "product", sku)
	if err != nil {
		return false, err
	}

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
	u, err := dc.BuildURL(nil, "product", sku)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)

	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	p := Product{}
	err = unmarshalBody(res.Body, p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) GetProducts(queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.BuildURL(queryFilter, "products")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) CreateProduct(np ProductCreationInput) (*http.Response, error) {
	body, err := createBodyFromStruct(np)
	if err != nil {
		return nil, err
	}

	u, err := dc.BuildURL(nil, "product")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) UpdateProduct(sku string, up ProductUpdateInput) (*http.Response, error) {
	body, err := createBodyFromStruct(up)
	if err != nil {
		return nil, err
	}

	u, err := dc.BuildURL(nil, "product", sku)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProduct(sku string) (*http.Response, error) {
	u, err := dc.BuildURL(nil, "product", sku)
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

func (dc *V1Client) GetProductRoot(rootID uint64) (*http.Response, error) {
	rootIDString := convertIDToString(rootID)
	u, err := dc.BuildURL(nil, "product_root", rootIDString)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) GetProductRoots(queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.BuildURL(queryFilter, "product_roots")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProductRoot(rootID uint64) (*http.Response, error) {
	rootIDString := convertIDToString(rootID)
	u, err := dc.BuildURL(nil, "product_root", rootIDString)
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

func (dc *V1Client) GetProductOptions(productID uint64, queryFilter map[string]string) (*http.Response, error) {
	productIDString := convertIDToString(productID)
	u, err := dc.BuildURL(queryFilter, "product", productIDString, "options")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) CreateProductOptionForProduct(productID uint64, JSONBody string) (*http.Response, error) {
	productIDString := convertIDToString(productID)
	body := strings.NewReader(JSONBody)
	u, err := dc.BuildURL(nil, "product", productIDString, "options")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) UpdateProductOption(optionID uint64, JSONBody string) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	body := strings.NewReader(JSONBody)
	u, err := dc.BuildURL(nil, "product_options", optionIDString)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProductOption(optionID uint64) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	u, err := dc.BuildURL(nil, "product_options", optionIDString)
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

func (dc *V1Client) createProductOptionValueForOption(optionID uint64, JSONBody string) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	body := strings.NewReader(JSONBody)
	u, err := dc.BuildURL(nil, "product_options", optionIDString, "value")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) updateProductOptionValueForOption(valueID uint64, JSONBody string) (*http.Response, error) {
	valueIDString := convertIDToString(valueID)
	body := strings.NewReader(JSONBody)
	u, err := dc.BuildURL(nil, "product_option_values", valueIDString)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) deleteProductOptionValueForOption(optionID uint64) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	u, err := dc.BuildURL(nil, "product_option_values", optionIDString)
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

func (dc *V1Client) getDiscountByID(discountID uint64) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u, err := dc.BuildURL(nil, "discount", discountIDString)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) getListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	u, err := dc.BuildURL(queryFilter, "discounts")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) createDiscount(JSONBody string) (*http.Response, error) {
	u, err := dc.BuildURL(nil, "discount")
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) updateDiscount(discountID uint64, JSONBody string) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u, err := dc.BuildURL(nil, "discount", discountIDString)
	if err != nil {
		return nil, err
	}

	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) deleteDiscount(discountID uint64) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u, err := dc.BuildURL(nil, "discount", discountIDString)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.executeRequest(req)
}
