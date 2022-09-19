module erc20

go 1.16

replace (
	erc20 => ./
	github.com/davecgh/go-spew => github.com/ipfs-force-community/go-spew v1.1.2-0.20220524052205-0034150c051a
	github.com/filecoin-project/go-address => github.com/ipfs-force-community/go-address v0.0.7-0.20220524010936-42617a156be1
	github.com/ipfs-force-community/go-fvm-sdk => ../..
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

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-state-types v0.1.12-alpha
	github.com/filecoin-project/specs-actors/v8 v8.0.1
	github.com/filecoin-project/venus v1.2.4
	github.com/ipfs-force-community/go-fvm-sdk v0.0.0-00010101000000-000000000000
	github.com/ipfs/go-cid v0.2.0
	github.com/stretchr/testify v1.7.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)
