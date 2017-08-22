package dairyclient

import (
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
	return dc.exists(u)
}

func (dc *V1Client) GetProduct(sku string) (*Product, error) {
	u := dc.buildURL(nil, "product", sku)
	p := Product{}

	err := dc.get(u, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (dc *V1Client) GetProducts(queryFilter map[string]string) ([]Product, error) {
	u := dc.buildURL(queryFilter, "products")
	pl := ProductList{}

	err := dc.get(u, &pl)
	if err != nil {
		return nil, err
	}

	return pl.Data, nil
}

func (dc *V1Client) CreateProduct(np ProductCreationInput) (*Product, error) {
	p := Product{}
	u := dc.buildURL(nil, "product")

	err := dc.post(u, np, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (dc *V1Client) UpdateProduct(sku string, up ProductUpdateInput) (*Product, error) {
	p := Product{}
	u := dc.buildURL(nil, "product", sku)

	err := dc.put(u, up, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (dc *V1Client) DeleteProduct(sku string) error {
	u := dc.buildURL(nil, "product", sku)
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

	r := ProductRoot{}
	err := dc.get(u, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (dc *V1Client) GetProductRoots(queryFilter map[string]string) ([]ProductRoot, error) {
	u := dc.buildURL(queryFilter, "product_roots")
	rl := ProductRootList{}
	err := dc.get(u, &rl)
	if err != nil {
		return nil, err
	}

	return rl.Data, nil
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
