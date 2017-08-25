package dairyclient

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

func (dc *V1Client) CreateProduct(np ProductInput) (*Product, error) {
	p := Product{}
	u := dc.buildURL(nil, "product")

	err := dc.post(u, np, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (dc *V1Client) UpdateProduct(sku string, up ProductInput) (*Product, error) {
	p := Product{}
	u := dc.buildURL(nil, "product", sku)

	err := dc.patch(u, up, &p)
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

func (dc *V1Client) GetProductOptions(productID uint64, queryFilter map[string]string) ([]ProductOption, error) {
	productIDString := convertIDToString(productID)
	u := dc.buildURL(queryFilter, "product", productIDString, "options")
	ol := ProductOptionList{}

	err := dc.get(u, &ol)
	if err != nil {
		return nil, err
	}

	return ol.Data, nil
}

func (dc *V1Client) CreateProductOptionForProduct(productID uint64, no ProductOption) (*ProductOption, error) {
	productIDString := convertIDToString(productID)
	o := ProductOption{}
	u := dc.buildURL(nil, "product", productIDString, "options")

	err := dc.post(u, no, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (dc *V1Client) UpdateProductOption(optionID uint64, uo ProductOption) (*ProductOption, error) {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString)
	o := ProductOption{}

	err := dc.patch(u, uo, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
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

func (dc *V1Client) CreateProductOptionValueForOption(optionID uint64, nv ProductOptionValue) (*ProductOptionValue, error) {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_options", optionIDString, "value")
	v := ProductOptionValue{}

	err := dc.post(u, nv, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (dc *V1Client) UpdateProductOptionValueForOption(valueID uint64, uv ProductOptionValue) (*ProductOptionValue, error) {
	valueIDString := convertIDToString(valueID)
	u := dc.buildURL(nil, "product_option_values", valueIDString)
	v := ProductOptionValue{}

	err := dc.patch(u, uv, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (dc *V1Client) DeleteProductOptionValueForOption(optionID uint64) error {
	optionIDString := convertIDToString(optionID)
	u := dc.buildURL(nil, "product_option_values", optionIDString)
	return dc.delete(u)
}
