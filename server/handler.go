package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cmwaters/skychart/types"
)

// Handler is the core object in the server package. It keeps an in-memory state
// of the chain-registry which can be updated using `Pull`. It handles requests
// for this data through the router.
type Handler struct {
	registryUrl  string
	lastUpdated  time.Time
	chains       []string
	assets       []string
	chainByAsset map[string]string // asset name -> chain name
	chainById    map[string]string // chain id -> chain name
	chainList    map[string]types.Chain
	assetList    map[string]types.AssetList
	log          *log.Logger
}

func NewHandler(registryUrl string, log *log.Logger) *Handler {
	return &Handler{
		registryUrl:  registryUrl,
		lastUpdated:  time.Unix(0, 0),
		chains:       make([]string, 0),
		assets:       make([]string, 0),
		chainByAsset: make(map[string]string),
		chainById:    make(map[string]string),
		chainList:    make(map[string]types.Chain),
		assetList:    make(map[string]types.AssetList),
		log:          log,
	}
}

func (h Handler) Chains(res http.ResponseWriter, req *http.Request) {
	respondWithJSON(res, h.chains)
}

// Chain searches for a chain by either name or ID and
// returns it if it exists
func (h Handler) Chain(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	chainName, ok := vars["chain"]
	if !ok {
		badRequest(res)
		return
	}

	exists, chain := h.findChain(chainName)
	if !exists {
		resourceNotFound(res)
		return
	}
	respondWithJSON(res, chain)
}

func (h Handler) Endpoints(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	chainName, ok := vars["chain"]
	if !ok {
		badRequest(res)
		return
	}
	endpointType, ok := vars["type"]
	if !ok {
		badRequest(res)
		return
	}
	exists, chain := h.findChain(chainName)
	if !exists {
		resourceNotFound(res)
		return
	}

	switch endpointType {
	case "rpc":
		respondWithJSON(res, chain.Apis.RPC)
	case "grpc":
		respondWithJSON(res, chain.Apis.Grpc)
	case "rest":
		respondWithJSON(res, chain.Apis.REST)
	case "peers":
		respondWithJSON(res, chain.Peers.PersistentPeers)
	case "seeds":
		respondWithJSON(res, chain.Peers.Seeds)
	default:
		badRequest(res)
	}
}

func (h Handler) ChainAsset(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	chainName, ok := vars["chain"]
	if !ok {
		badRequest(res)
		return
	}
	assets, ok := h.assetList[chainName]
	if !ok {
		chainName, ok = h.chainById[chainName]
		if !ok {
			badRequest(res)
		}
		assets = h.assetList[chainName]
	}
	respondWithJSON(res, assets)
}

func (h Handler) Assets(res http.ResponseWriter, req *http.Request) {
	respondWithJSON(res, h.assets)
}

func (h Handler) Asset(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	assetName, ok := vars["asset"]
	if !ok {
		badRequest(res)
		return
	}
	chainName, ok := h.chainByAsset[assetName]
	if !ok {
		resourceNotFound(res)
		return
	}

	assetList := h.assetList[chainName]
	for _, asset := range assetList.Assets {
		if asset.Display == assetName {
			respondWithJSON(res, asset)
			return
		}
	}

	resourceNotFound(res)
}

func (h Handler) findChain(name string) (bool, types.Chain) {
	chain, ok := h.chainList[name]
	if ok {
		return true, chain
	}

	name, ok = h.chainById[name]
	if !ok {
		return false, types.Chain{}
	}

	return true, h.chainList[name]
}

func respondWithJSON(w http.ResponseWriter, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func resourceNotFound(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.WriteHeader(http.StatusNotFound)
}

func badRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.WriteHeader(http.StatusBadRequest)
}
