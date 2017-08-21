package dairyclient

import (
	"fmt"
	"net/http"
	"strings"
)

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
	err = unmarshalBody(res, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) GetProducts(queryFilter map[string]string) ([]Product, error) {
	u := dc.buildURL(queryFilter, "products")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	p := ProductList{}
	err = unmarshalBody(res, &p)
	if err != nil {
		return nil, err
	}
	return p.Data, nil
}

func (dc *V1Client) CreateProduct(np ProductCreationInput) (*Product, error) {
	body, err := createBodyFromStruct(np)
	if err != nil {
		return nil, err
	}

	u := dc.buildURL(nil, "product")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	p := Product{}
	err = unmarshalBody(res, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) UpdateProduct(sku string, up ProductUpdateInput) (*Product, error) {
	body, err := createBodyFromStruct(up)
	if err != nil {
		return nil, err
	}

	u := dc.buildURL(nil, "product", sku)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	p := Product{}
	err = unmarshalBody(res, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) DeleteProduct(sku string) error {
	u := dc.buildURL(nil, "product", sku)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return dc.delete(u)
}

////////////////////////////////////////////////////////
//                                                    //
//              Product Root Functions                //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetProductRoot(rootID uint64) (*ProductRoot, error) {
	rootIDString := convertIDToString(rootID)
	u := dc.buildURL(nil, "product_root", rootIDString)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	r := ProductRoot{}
	err = unmarshalBody(res, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (dc *V1Client) GetProductRoots(queryFilter map[string]string) (*http.Response, error) {
	u := dc.buildURL(queryFilter, "product_roots")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProductRoot(rootID uint64) error {
	rootIDString := convertIDToString(rootID)
	u := dc.buildURL(nil, "product_root", rootIDString)
	return dc.delete(u)
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

func (dc *V1Client) DeleteProductOption(optionID uint64) error {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString)
	return dc.delete(u)
}

////////////////////////////////////////////////////////
//                                                    //
//          Product Option Value Functions            //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) CreateProductOptionValueForOption(optionID uint64, JSONBody string) (*http.Response, error) {
	optionIDString := convertIDToString(optionID)
	body := strings.NewReader(JSONBody)
	u := dc.buildURL(nil, "product_options", optionIDString, "value")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) UpdateProductOptionValueForOption(valueID uint64, JSONBody string) (*http.Response, error) {
	valueIDString := convertIDToString(valueID)
	body := strings.NewReader(JSONBody)
	u := dc.buildURL(nil, "product_option_values", valueIDString)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteProductOptionValueForOption(optionID uint64) error {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_option_values", optionIDString)
	return dc.delete(u)
}
