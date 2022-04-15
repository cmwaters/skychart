![](skymap.jpg)
# Skymap

Skymap is a simple golang server and client library for the cosmos chain-registry. It provides a convenient
API, automatically updating itself to any changes in the github repo once a day. In the types package you will
find go generated types from the JSON schemas.

## Usage

Install the binary

```cli
go install github.com/cmwaters/skymap@latest
```

and run the server (this is for port :8080)

```cli
skymap cosmos/chain-registry :8080
```

## API Reference


| Query | Description | Response Type |
|-------|-------------|---------------|
| `/v1/chains` | Returns an array of registered chains by name  | `[]string` |
| `/v1/chain/{chain}` | Returns a registered chain if it exists | `Chain` |
| `/v1/chain/{chain}/endpoints/rpc` | Returns a list of active public RPC endpoints | `[]GrpcElement` |
| `/v1/chain/{chain}/endpoints/rest` | Returns a list of active public REST endpoints | `[]GrpcElement` |
| `/v1/chain/{chain}/endpoints/grpc` | Returns a list of active public gRPC endpoints | `[]GrpcElement` |
| `/v1/chain/{chain}/endpoints/peers` | Returns a list of chain peers | `[]PersistentPeerElement` |
| `/v1/chain/{chain}/endpoints/seeds` | Returns a list of chain seeds | `[]PersistentPeerElement` |
| `/v1/chain/{chain}/assets` | Returns all the native assets of the chain | `AssetList` |
| `/v1/assets` | Returns an array of registered assets by display name | `[]string` |
| `/v1/asset/{asset}` | Returns an asset by display name if it exists | `AssetElement` | 

