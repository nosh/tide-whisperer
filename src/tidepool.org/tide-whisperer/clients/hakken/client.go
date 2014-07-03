package hakken

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tidepool.org/common/errors"
	"tidepool.org/tide-whisperer/clients/disc"
)

type coordinatorClient struct {
	coordinator Coordinator
}

func (client *coordinatorClient) getCoordinators() ([]Coordinator, error) {
	url := fmt.Sprintf("%s/v1/coordinator", client.coordinator.URL.String())
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Problem when looking up coordinators[%s].", url)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Newf("Unknown response code[%s] from [%s]", res.StatusCode, url)
	}

	var retVal []Coordinator
	if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
		return nil, errors.Wrapf(err, "Error parsing JSON results from [%s]", url)
	}
	return retVal, nil
}

func (client *coordinatorClient) getListings(service string) ([]disc.ServiceListing, error) {
	url := fmt.Sprintf("%s/v1/listings/%s", client.coordinator.URL.String(), service)
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Problem when looking up listings at url[%s].", url)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Newf("Unknown response code[%s] from url[%s]", res.StatusCode, url)
	}

	var retVal []disc.ServiceListing
	if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
		return nil, errors.Wrapf(err, "Error parsing JSON results from url[%s]", url)
	}
	return retVal, nil
}

func (client *coordinatorClient) listingHearbeat(sl disc.ServiceListing) error {
	url := fmt.Sprintf("%s/v1/listings?heartbeat=true", client.coordinator.URL.String())

	out, in := io.Pipe()
	json.NewEncoder(in).Encode(sl)
	res, err := http.Post(url, "application/json", out)
	if err != nil {
		return errors.Wrapf(err, "Problem when updating heartbeat for service[%s] at [%s].", sl.Service, url)
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return errors.Newf("Unknown response code[%s] from url[%s]", res.StatusCode, url)
	}
	return nil
}
