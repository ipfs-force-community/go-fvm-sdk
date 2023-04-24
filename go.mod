module github.com/ipfs-force-community/go-fvm-sdk

go 1.20

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-amt-ipld/v4 v4.0.0
	github.com/filecoin-project/go-crypto v0.0.1
	github.com/filecoin-project/go-hamt-ipld/v3 v3.1.0
	github.com/filecoin-project/go-state-types v0.9.9
	github.com/filecoin-project/specs-actors v0.9.13
	github.com/filecoin-project/specs-actors/v7 v7.0.0
	github.com/ipfs/go-cid v0.2.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/minio/sha256-simd v1.0.0
	github.com/multiformats/go-multihash v0.0.15
	github.com/stretchr/testify v1.7.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/exp v0.0.0-20221114191408-850992195362
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/filecoin-project/specs-actors/v5 v5.0.4 // indirect
	github.com/ipfs/go-block-format v0.0.3 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.6 // indirect
	github.com/ipsn/go-secp256k1 v0.0.0-20180726113642-9d62b9f0bc52 // indirect
	github.com/klauspost/cpuid/v2 v2.0.6 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20190809202753-05966cbd336a // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
)

replace (
	github.com/davecgh/go-spew => github.com/ipfs-force-community/go-spew v1.1.2-0.20220524052205-0034150c051a
	github.com/filecoin-project/go-address => github.com/ipfs-force-community/go-address v0.0.7-0.20230207015848-7a27d889c267
	// github.com/filecoin-project/go-state-types => github.com/ipfs-force-community/go-state-types v0.1.8-0.20220523102255-d400d329d9c1
	github.com/ipfs/go-block-format => github.com/ipfs-force-community/go-block-format v0.0.4-0.20220425095807-073e9266335c
	//remove anything about reflect
	github.com/ipfs/go-ipld-cbor => github.com/ipfs-force-community/go-ipld-cbor v0.0.7-0.20220713070731-f5190aacb1a4
	github.com/klauspost/cpuid/v2 => github.com/ipfs-force-community/cpuid/v2 v2.0.13-0.20220523085810-ac111993ce74
	github.com/minio/blake2b-simd => github.com/ipfs-force-community/blake2b-simd v0.0.0-20220523083450-6e9a68832d69
	github.com/minio/sha256-simd => github.com/ipfs-force-community/sha256-simd v1.0.1-0.20220421100150-fcbba4b6ea96
	github.com/polydawn/refmt => github.com/hunjixin/refmt v0.0.0-20220520091210-cb3c7d292019
	//remove http/file asset and reflect code
	github.com/stretchr/testify => github.com/ipfs-force-community/testify v1.7.1-0.20220616060316-ea4f53121ac3
	github.com/whyrusleeping/cbor-gen => github.com/ipfs-force-community/cbor-gen v0.0.0-20220421100448-dc345220256c
	//remove implement baseon specify platform only keep pure go implement
	golang.org/x/crypto => github.com/ipfs-force-community/crypto v0.0.0-20220523090957-2aff239c26f7
	lukechampine.com/blake3 => github.com/ipfs-force-community/blake3 v1.1.8-0.20220609024944-51450f2b2fc0
)
