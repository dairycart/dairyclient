package dairyclient

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetDiscountByID(discountID uint64) (*Discount, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	d := Discount{}

	err := dc.get(u, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (dc *V1Client) GetListOfDiscounts(queryFilter map[string]string) ([]Discount, error) {
	u := dc.buildURL(nil, "discount")
	d := DiscountList{}

	err := dc.get(u, &d)
	if err != nil {
		return nil, err
	}
	return d.Data, nil
}

func (dc *V1Client) CreateDiscount(nd Discount) (*Discount, error) {
	d := Discount{}
	u := dc.buildURL(nil, "discount")

	err := dc.post(u, nd, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (dc *V1Client) UpdateDiscount(discountID uint64, ud Discount) (*Discount, error) {
	d := Discount{}
	u := dc.buildURL(nil, "discount")

	err := dc.patch(u, ud, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (dc *V1Client) DeleteDiscount(discountID uint64) error {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	return dc.delete(u)
}
