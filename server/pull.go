package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cmwaters/skymap/types"
)

// Pull requests all registry information from a github repo and updates the
// handlers local registry. It expects a directory structure as follows:
// - [chain_name]
//   - chain.json
//   - assetlist.json
// It works on a best effort basis. All chain names should be unique. chain.json and
// assetlist.json should comply with the respective schemas
// TODO: Add support for relayer paths
func (h *Handler) Pull(ctx context.Context) error {
	// If there have been no recent commits we can return immediately
	recent, err := h.recentCommits()
	if err != nil {
		return err
	}
	if !recent {
		h.log.Printf("no new recent commits since %s", h.lastUpdated.String())
		h.lastUpdated = time.Now()
		return nil
	}

	// update chains
	if err := h.getChains(); err != nil {
		return err
	}

	// for each chain update the chain info and asset list
	// TODO: If we wanted to be more creative we could first check
	// to see if the file had actually changed since the last time
	// it was pulled
	for _, chain := range h.chains {
		if err := h.getChain(chain); err != nil {
			return err
		}
		if err := h.getAssetList(chain); err != nil {
			return err
		}
	}

	// Index assets by display
	assets := make([]string, 0)
	for _, assetList := range h.assetList {
		name := h.chainById[assetList.ChainID]
		for _, asset := range assetList.Assets {
			assets = append(assets, asset.Display)
			h.chainByAsset[asset.Display] = name
		}
	}

	// update timestamp
	h.lastUpdated = time.Now()
	h.log.Printf("successfully updated registry (%d chains)", len(h.chains))

	return nil
}

func (h *Handler) getChains() error {
	query := fmt.Sprintf("https://api.github.com/repos/%s/contents", h.registryUrl)
	resp, err := http.Get(query)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from query %s: %d", query, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var repo []map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &repo); err != nil {
		return fmt.Errorf("unmarshalling repo: %w", err)
	}

	chains := make([]string, 0)
	for _, entry := range repo {
		// only accept directories
		entryType := entry["type"].(string)
		if entryType != "dir" {
			continue
		}

		name := entry["name"].(string)
		if strings.Contains(name, "testnets") {
			continue
		}
		if strings.Contains(name, ".") {
			continue
		}

		chains = append(chains, name)

		fmt.Printf("%s, ", name)
	}
	fmt.Printf("\n")
	h.chains = chains
	fmt.Print(len(h.chains))
	return nil
}

func (h *Handler) getChain(name string) error {
	query := fmt.Sprintf("https://raw.githubusercontent.com/%s/master/%s/chain.json", h.registryUrl, name)
	resp, err := http.Get(query)
	if err != nil {
		return err
	}

	// If the chain.json file doesn't exist we simply ignore it
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from query %s: %d", query, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var chain types.Chain
	err = json.Unmarshal(bodyBytes, &chain)
	if err != nil {
		return err
	}

	h.chainList[name] = chain
	h.chainById[chain.ChainID] = name
	return nil
}

func (h *Handler) getAssetList(name string) error {
	query := fmt.Sprintf("https://raw.githubusercontent.com/%s/master/%s/assetlist.json", h.registryUrl, name)
	resp, err := http.Get(query)
	if err != nil {
		return err
	}

	// If the chain.json file doesn't exist we simply ignore it
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code for query %s: %d", query, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var assetList types.AssetList
	err = json.Unmarshal(bodyBytes, &assetList)
	if err != nil {
		return err
	}

	h.assetList[name] = assetList
	return nil
}

// recentCommits returns true if there has been a commit more recent than the time the handler
// last updated
func (h Handler) recentCommits() (bool, error) {
	lastUpdated := h.lastUpdated.Format(time.RFC3339)
	query := fmt.Sprintf("https://api.github.com/repos/%s/commits?since=%s", h.registryUrl, lastUpdated)
	resp, err := http.Get(query)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code for query %s: %d", query, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	var body []interface{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return false, err
	}

	return len(body) > 0, nil
}
