module github.com/ipfs-force-community/go-fvm-sdk

go 1.17

require github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799

require (
	github.com/ipfs/go-cid v0.1.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.1.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	lukechampine.com/blake3 v1.1.6 // indirect
)

replace hellocontract => ./hellocontract

replace github.com/ipfs-force-community/go-fvm-sdk => ./

replace github.com/ipfs-force-community/go-fvm-sdk/sdk => ./sdk
