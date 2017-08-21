package dairyclient

import (
	"net/http"
	"strings"
)

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func (dc *V1Client) GetDiscountByID(discountID uint64) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) GetListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	u := dc.buildURL(queryFilter, "discounts")
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return dc.executeRequest(req)
}

func (dc *V1Client) CreateDiscount(JSONBody string) (*http.Response, error) {
	u := dc.buildURL(nil, "discount")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) UpdateDiscount(discountID uint64, JSONBody string) (*http.Response, error) {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return dc.executeRequest(req)
}

func (dc *V1Client) DeleteDiscount(discountID uint64) error {
	discountIDString := convertIDToString(discountID)
	u := dc.buildURL(nil, "discount", discountIDString)
	return dc.Delete(u)
}
