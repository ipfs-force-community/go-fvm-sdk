module hellocontract

go 1.16

replace github.com/klauspost/cpuid/v2 => github.com/ipfs-force-community/cpuid/v2 v2.0.13-0.20220421095210-bfbeb72f34dd

replace github.com/whyrusleeping/cbor-gen => github.com/ipfs-force-community/cbor-gen v0.0.0-20220421100448-dc345220256c

replace github.com/minio/sha256-simd => github.com/ipfs-force-community/sha256-simd v1.0.1-0.20220421100150-fcbba4b6ea96

replace golang.org/x/crypto => github.com/ipfs-force-community/crypto v0.0.0-20220421095836-dd8044371872

replace github.com/ipfs-force-community/go-fvm-sdk/sdk => ../sdk

replace hellocontract => ./

require (
	github.com/ipfs-force-community/go-fvm-sdk/sdk v0.0.0-00010101000000-000000000000
	github.com/ipfs/go-cid v0.1.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)

require github.com/multiformats/go-multihash v0.1.0 // indirect
