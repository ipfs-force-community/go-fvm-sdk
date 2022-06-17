module hellocontract

go 1.16

replace (
	github.com/ipfs-force-community/go-fvm-sdk => ../..
	github.com/ipfs/go-block-format => github.com/ipfs-force-community/go-block-format v0.0.4-0.20220425095807-073e9266335c
	github.com/klauspost/cpuid/v2 => github.com/ipfs-force-community/cpuid/v2 v2.0.13-0.20220421095210-bfbeb72f34dd
	github.com/minio/sha256-simd => github.com/ipfs-force-community/sha256-simd v1.0.1-0.20220421100150-fcbba4b6ea96
	github.com/whyrusleeping/cbor-gen => github.com/ipfs-force-community/cbor-gen v0.0.0-20220421100448-dc345220256c
	golang.org/x/crypto => github.com/ipfs-force-community/crypto v0.0.0-20220421095836-dd8044371872
	hellocontract => ./
	lukechampine.com/blake3 => github.com/ipfs-force-community/blake3 v1.1.8-0.20220609024944-51450f2b2fc0
)

require (
	github.com/ipfs/go-cid v0.1.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-state-types v0.1.9
	github.com/filecoin-project/specs-actors/v8 v8.0.0-20220412224951-92abd0e6e7ae
	github.com/filecoin-project/venus v1.2.4
	github.com/ipfs-force-community/go-fvm-sdk v0.0.0-00010101000000-000000000000
)
