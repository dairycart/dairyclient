package dairyclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func NewV1ClientFromCookie(apiURL string, cookie *http.Cookie, client *http.Client) (*V1Client, error) {
	var dc *V1Client
	if client != nil {
		dc = &V1Client{Client: client}
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, errors.Wrap(err, "API URL is not valid")
	}
	dc.URL = u

	dc.AuthCookie = cookie
	dc.Client.Timeout = 5 * time.Second

	return dc, nil
}

func (dc *V1Client) executeRequest(req *http.Request) (*http.Response, error) {
	req.AddCookie(dc.AuthCookie)
	return dc.Do(req)
}

func (dc *V1Client) buildURL(queryParams map[string]string, parts ...string) string {
	parts = append([]string{"v1"}, parts...)
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
	parts = append([]string{"v1"}, parts...)

	u, err := url.Parse(strings.Join(parts, "/"))
	if err != nil {
		return "", err
	}

	queryString := mapToQueryValues(queryParams)
	u.RawQuery = queryString.Encode()
	return dc.URL.ResolveReference(u).String(), nil
}

func (dc *V1Client) exists(uri string) (bool, error) {
	req, _ := http.NewRequest(http.MethodHead, uri, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return false, err
	}

	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (dc *V1Client) get(uri string, obj interface{}) error {
	if err := interfaceArgIsNotPointerOrNil(obj); err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return err
	}

	err = unmarshalBody(res, &obj)
	if err != nil {
		return err
	}
	return nil
}

func (dc *V1Client) delete(uri string) error {
	req, _ := http.NewRequest(http.MethodDelete, uri, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("user couldn't be deleted, status returned: %d", res.StatusCode)
	}
	return nil
}

func (dc *V1Client) makeDataRequest(method string, uri string, in interface{}, out interface{}) error {
	if err := interfaceArgIsNotPointerOrNil(out); err != nil {
		return err
	}

	body, err := createBodyFromStruct(in)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(method, uri, body)
	res, err := dc.executeRequest(req)
	if err != nil {
		return err
	}

	err = unmarshalBody(res, &out)
	if err != nil {
		return err
	}
	return nil
}

func (dc *V1Client) post(uri string, in interface{}, out interface{}) error {
	return dc.makeDataRequest(http.MethodPost, uri, in, out)
}

func (dc *V1Client) patch(uri string, in interface{}, out interface{}) error {
	return dc.makeDataRequest(http.MethodPatch, uri, in, out)
}
