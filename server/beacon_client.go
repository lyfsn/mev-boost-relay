package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type ValidatorService interface {
	IsValidator(string) bool
	NumValidators() uint64
	FetchValidators() error
}

// type DevValidatorService struct {
// 	mu           sync.RWMutex
// 	validatorSet map[string]validatorResponseEntry
// }

// func (d *DevValidatorService) IsValidator(pubkey string) bool {
// 	d.mu.RLock()
// 	pkLower := strings.ToLower(pubkey)
// 	_, found := d.validatorSet[pkLower]
// 	d.mu.RUnlock()
// 	return found
// }

// func (d *DevValidatorService) NumValidators() uint64 {
// 	d.mu.RLock()
// 	defer d.mu.RUnlock()
// 	return uint64(len(d.validatorSet))
// }

// func (d *DevValidatorService) FetchValidators() error {
// 	return nil
// }

type BeaconClientValidatorService struct {
	beaconEndpoint string
	mu             sync.RWMutex
	validatorSet   map[string]validatorResponseEntry
}

func NewBeaconClientValidatorService(beaconEndpoint string) *BeaconClientValidatorService {
	return &BeaconClientValidatorService{
		beaconEndpoint: beaconEndpoint,
		validatorSet:   make(map[string]validatorResponseEntry),
	}
}

func (b *BeaconClientValidatorService) IsValidator(pubkey string) bool {
	b.mu.RLock()
	pkLower := strings.ToLower(pubkey)
	_, found := b.validatorSet[pkLower]
	b.mu.RUnlock()
	return found
}

func (b *BeaconClientValidatorService) NumValidators() uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return uint64(len(b.validatorSet))
}

func (b *BeaconClientValidatorService) FetchValidators() error {
	vd, err := fetchAllValidators(b.beaconEndpoint)
	if err != nil {
		return err
	}

	newValidatorSet := make(map[string]validatorResponseEntry)
	for _, vs := range vd.Data {
		pkLower := strings.ToLower(vs.Validator.Pubkey)
		newValidatorSet[pkLower] = vs
	}

	b.mu.Lock()
	b.validatorSet = newValidatorSet
	b.mu.Unlock()
	return nil
}

type validatorResponseEntry struct {
	Validator struct {
		Pubkey string `json:"pubkey"`
	} `json:"validator"`
}

type allValidatorsResponse struct {
	Data []validatorResponseEntry
}

func fetchAllValidators(endpoint string) (*allValidatorsResponse, error) {
	uri := endpoint + "/eth/v1/beacon/states/head/validators?status=active,pending"

	// https://ethereum.github.io/beacon-APIs/#/Beacon/getStateValidators
	vd := new(allValidatorsResponse)
	err := fetchBeacon(uri, "GET", vd)
	return vd, err
}

func fetchBeacon(url string, method string, dst any) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("invalid reqest for %s: %w", url, err)
	}
	req.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("client refused for %s: %w", url, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body for %s: %w", url, err)
	}

	if resp.StatusCode >= 300 {
		ec := &struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{}
		if err = json.Unmarshal(bodyBytes, ec); err != nil {
			return fmt.Errorf("could not unmarshal error response from beacon node for %s from %s: %w", url, string(bodyBytes), err)
		}
		return errors.New(ec.Message)
	}

	err = json.Unmarshal(bodyBytes, dst)
	if err != nil {
		return fmt.Errorf("could not unmarshal response for %s from %s: %w", url, string(bodyBytes), err)
	}

	return nil
}
