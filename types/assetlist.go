package types

// Asset lists are a similar mechanism to allow frontends and other UIs to fetch metadata
// associated with Cosmos SDK denoms, especially for assets sent over IBC.
type AssetList struct {
	Assets    []AssetElement `json:"assets"`
	ChainName string         `json:"chain_name"`
}

type AssetElement struct {
	Address     *string            `json:"address,omitempty"`
	Base        string             `json:"base"`                   // The base unit of the asset. Must be in denom_units.
	CoingeckoID *string            `json:"coingecko_id,omitempty"` // The coingecko id to fetch asset data from coingecko v3 api. See; https://api.coingecko.com/api/v3/coins/list
	DenomUnits  []DenomUnitElement `json:"denom_units"`
	Description *string            `json:"description,omitempty"` // A short description of the asset
	Display     string             `json:"display"`               // The human friendly unit of the asset. Must be in denom_units.
	Ibc         *Ibc               `json:"ibc,omitempty"`
	Kind        *Kind              `json:"kind,omitempty"` // The potential options for type of asset. By default, assumes sdk.coin
	LogoURIs    *LogoURIs          `json:"logo_URIs,omitempty"`
	Name        *string            `json:"name,omitempty"`   // The project name of the asset. For example Bitcoin.
	Symbol      *string            `json:"symbol,omitempty"` // The symbol of an asset. For example BTC.
}

type DenomUnitElement struct {
	Aliases  []string `json:"aliases,omitempty"`
	Denom    string   `json:"denom"`
	Exponent int64    `json:"exponent"`
}

type Ibc struct {
	DstChannel    string `json:"dst_channel"`
	SourceChannel string `json:"source_channel"`
	SourceDenom   string `json:"source_denom"`
}

type LogoURIs struct {
	PNG *string `json:"png,omitempty"`
	SVG *string `json:"svg,omitempty"`
}

// The potential options for type of asset. By default, assumes sdk.coin
type Kind string

const (
	Cw20    Kind = "cw20"
	Erc20   Kind = "erc20"
	SDKCoin Kind = "sdk.coin"
	Snip20  Kind = "snip20"
)
