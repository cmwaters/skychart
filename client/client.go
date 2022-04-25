package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cmwaters/skychart/types"
)

// Client is a simple wrapper for performing http requests to the registry and
// parsing the corresponding response
type Client struct {
	registryUrl string
}

func New(registryUrl string) (*Client, error) {
	_, err := url.Parse(registryUrl)
	if err != nil {
		return nil, err
	}
	return &Client{registryUrl: registryUrl}, nil
}

func (c Client) Chains() ([]string, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chains", c.registryUrl))
	if err != nil {
		return nil, err
	}
	var chains []string
	err = json.Unmarshal(bz, &chains)
	if err != nil {
		return nil, err
	}

	return chains, nil
}

func (c Client) Assets() ([]string, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/assets", c.registryUrl))
	if err != nil {
		return nil, err
	}
	var assets []string
	err = json.Unmarshal(bz, &assets)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (c Client) Chain(chain string) (types.Chain, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chain/%s", c.registryUrl, chain))
	if err != nil {
		return types.Chain{}, err
	}
	var resp types.Chain
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return types.Chain{}, err
	}
	return resp, nil
}

func (c Client) Asset(name string) (types.AssetElement, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/asset/%s", c.registryUrl, name))
	if err != nil {
		return types.AssetElement{}, err
	}
	var resp types.AssetElement
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return types.AssetElement{}, err
	}
	return resp, nil
}

func (c Client) RPC(chain string) ([]types.GrpcElement, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chain/%s/endpoints/rpc", c.registryUrl, chain))
	if err != nil {
		return []types.GrpcElement{}, err
	}
	var resp []types.GrpcElement
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return []types.GrpcElement{}, err
	}
	return resp, nil
}

func (c Client) GRPC(chain string) ([]types.GrpcElement, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chain/%s/endpoints/grpc", c.registryUrl, chain))
	if err != nil {
		return []types.GrpcElement{}, err
	}
	var resp []types.GrpcElement
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return []types.GrpcElement{}, err
	}
	return resp, nil
}

func (c Client) REST(chain string) ([]types.GrpcElement, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chain/%s/endpoints/rest", c.registryUrl, chain))
	if err != nil {
		return []types.GrpcElement{}, err
	}
	var resp []types.GrpcElement
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return []types.GrpcElement{}, err
	}
	return resp, nil
}

func (c Client) Peers(chain string) ([]types.PersistentPeerElement, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chain/%s/endpoints/peers", c.registryUrl, chain))
	if err != nil {
		return []types.PersistentPeerElement{}, err
	}
	var resp []types.PersistentPeerElement
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return []types.PersistentPeerElement{}, err
	}
	return resp, nil
}

func (c Client) Seeds(chain string) ([]types.PersistentPeerElement, error) {
	bz, err := c.get(fmt.Sprintf("%s/v1/chain/%s/endpoints/seeds", c.registryUrl, chain))
	if err != nil {
		return []types.PersistentPeerElement{}, err
	}
	var resp []types.PersistentPeerElement
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return []types.PersistentPeerElement{}, err
	}
	return resp, nil
}

func (c Client) get(query string) ([]byte, error) {
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("resource not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}
