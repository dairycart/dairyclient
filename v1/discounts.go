package dairyclient

import (
	"github.com/dairycart/dairymodels/v1"
)

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetDiscountByID(discountID uint64) (*models.Discount, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	d := models.Discount{}

	err := dc.get(u, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

//func (dc *V1Client) GetListOfDiscounts(queryFilter map[string]string) ([]models.Discount, error) {
	func (dc *V1Client) GetListOfDiscounts(queryFilter map[string]string) (*models.ListResponse, error) {
	u := dc.buildURL(nil, "discount")
	d := &models.ListResponse{}

	err := dc.get(u, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (dc *V1Client) CreateDiscount(nd models.Discount) (*models.Discount, error) {
	d := models.Discount{}
	u := dc.buildURL(nil, "discount")

	err := dc.post(u, nd, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (dc *V1Client) UpdateDiscount(discountID uint64, ud models.Discount) (*models.Discount, error) {
	d := models.Discount{}
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
