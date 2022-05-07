module github.com/ipfs-force-community/go-fvm-sdk/sdk

go 1.17

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-state-types v0.1.3
	github.com/filecoin-project/specs-actors/v2 v2.3.5-0.20210114162132-5b58b773f4fb
	github.com/filecoin-project/specs-actors/v7 v7.0.0
	github.com/ipfs/go-cid v0.1.0
	github.com/stretchr/testify v1.7.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/filecoin-project/specs-actors v0.9.13 // indirect
	github.com/filecoin-project/specs-actors/v5 v5.0.4 // indirect
	github.com/ipfs/go-block-format v0.0.3 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.5 // indirect
	github.com/ipfs/go-ipld-format v0.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.0.6 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.0.15 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20190809202753-05966cbd336a // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)

replace github.com/klauspost/cpuid/v2 => github.com/ipfs-force-community/cpuid/v2 v2.0.13-0.20220421095210-bfbeb72f34dd

replace github.com/whyrusleeping/cbor-gen => github.com/ipfs-force-community/cbor-gen v0.0.0-20220421100448-dc345220256c

replace github.com/minio/sha256-simd => github.com/ipfs-force-community/sha256-simd v1.0.1-0.20220421100150-fcbba4b6ea96

replace golang.org/x/crypto => github.com/ipfs-force-community/crypto v0.0.0-20220421095836-dd8044371872

replace github.com/ipfs/go-block-format => github.com/ipfs-force-community/go-block-format v0.0.4-0.20220425095807-073e9266335c

replace github.com/stretchr/testify => github.com/ipfs-force-community/testify v1.7.1-0.20220507025933-e761b418477e

replace github.com/davecgh/go-spew => github.com/ipfs-force-community/go-spew v1.1.2-0.20220507024706-1904e9f50471
