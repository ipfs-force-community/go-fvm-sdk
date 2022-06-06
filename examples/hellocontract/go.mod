module hellocontract

go 1.16

replace (
	github.com/ipfs/go-block-format => github.com/ipfs-force-community/go-block-format v0.0.4-0.20220425095807-073e9266335c
	github.com/klauspost/cpuid/v2 => github.com/ipfs-force-community/cpuid/v2 v2.0.13-0.20220421095210-bfbeb72f34dd
	github.com/minio/sha256-simd => github.com/ipfs-force-community/sha256-simd v1.0.1-0.20220421100150-fcbba4b6ea96
	github.com/whyrusleeping/cbor-gen => github.com/ipfs-force-community/cbor-gen v0.0.0-20220421100448-dc345220256c
	golang.org/x/crypto => github.com/ipfs-force-community/crypto v0.0.0-20220421095836-dd8044371872
	hellocontract => ./
)

require (
	github.com/ipfs/go-cid v0.1.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-crypto v0.0.1 // indirect
	github.com/filecoin-project/specs-actors v0.9.14 // indirect
	github.com/filecoin-project/specs-actors/v2 v2.3.6 // indirect
	github.com/ipfs-force-community/go-fvm-sdk/sdk v0.0.0-20220606080235-f183d12b8045
	github.com/ipfs/go-ipld-cbor v0.0.6 // indirect
	github.com/ipfs/go-ipld-format v0.2.0 // indirect
	github.com/multiformats/go-multihash v0.1.0 // indirect
	github.com/polydawn/refmt v0.0.0-20201211092308-30ac6d18308e // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/warpfork/go-wish v0.0.0-20200122115046-b9ea61034e4a // indirect
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
