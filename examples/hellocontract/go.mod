module hellocontract

go 1.18

replace (
	erc20 => ./
	github.com/davecgh/go-spew => github.com/ipfs-force-community/go-spew v1.1.2-0.20220524052205-0034150c051a
	github.com/filecoin-project/go-address => github.com/ipfs-force-community/go-address v0.0.7-0.20220524010936-42617a156be1
	github.com/ipfs-force-community/go-fvm-sdk => ../..
	github.com/ipfs/go-block-format => github.com/ipfs-force-community/go-block-format v0.0.4-0.20220425095807-073e9266335c
	github.com/ipfs/go-ipld-cbor => github.com/ipfs-force-community/go-ipld-cbor v0.0.7-0.20220713070731-f5190aacb1a4
	github.com/klauspost/cpuid/v2 => github.com/ipfs-force-community/cpuid/v2 v2.0.13-0.20220523085810-ac111993ce74
	github.com/minio/blake2b-simd => github.com/ipfs-force-community/blake2b-simd v0.0.0-20220523083450-6e9a68832d69
	github.com/minio/sha256-simd => github.com/ipfs-force-community/sha256-simd v1.0.1-0.20220421100150-fcbba4b6ea96
	github.com/polydawn/refmt => github.com/hunjixin/refmt v0.0.0-20220520091210-cb3c7d292019
	github.com/stretchr/testify => github.com/ipfs-force-community/testify v1.7.1-0.20220616060316-ea4f53121ac3
	github.com/whyrusleeping/cbor-gen => github.com/ipfs-force-community/cbor-gen v0.0.0-20220421100448-dc345220256c
	golang.org/x/crypto => github.com/ipfs-force-community/crypto v0.0.0-20220523090957-2aff239c26f7
	lukechampine.com/blake3 => github.com/ipfs-force-community/blake3 v1.1.8-0.20220609024944-51450f2b2fc0
)

require (
	github.com/ipfs/go-cid v0.2.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20220323183124-98fa8256a799
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f
)

require (
	github.com/filecoin-project/go-address v0.0.6
	github.com/filecoin-project/go-state-types v0.1.12-alpha
	github.com/filecoin-project/specs-actors/v8 v8.0.0-20220412224951-92abd0e6e7ae
	github.com/filecoin-project/venus v1.2.4
	github.com/ipfs-force-community/go-fvm-sdk v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/filecoin-project/go-amt-ipld/v2 v2.1.1-0.20201006184820-924ee87a1349 // indirect
	github.com/filecoin-project/go-amt-ipld/v3 v3.1.0 // indirect
	github.com/filecoin-project/go-amt-ipld/v4 v4.0.0 // indirect
	github.com/filecoin-project/go-bitfield v0.2.4 // indirect
	github.com/filecoin-project/go-commp-utils v0.1.3 // indirect
	github.com/filecoin-project/go-commp-utils/nonffi v0.0.0-20220905160352-62059082a837 // indirect
	github.com/filecoin-project/go-crypto v0.0.1 // indirect
	github.com/filecoin-project/go-fil-commcid v0.1.0 // indirect
	github.com/filecoin-project/go-hamt-ipld v0.1.5 // indirect
	github.com/filecoin-project/go-hamt-ipld/v2 v2.0.0 // indirect
	github.com/filecoin-project/go-hamt-ipld/v3 v3.1.0 // indirect
	github.com/filecoin-project/go-jsonrpc v0.1.4-0.20210217175800-45ea43ac2bec // indirect
	github.com/filecoin-project/specs-actors v0.9.14 // indirect
	github.com/filecoin-project/specs-actors/v2 v2.3.6 // indirect
	github.com/filecoin-project/specs-actors/v3 v3.1.1 // indirect
	github.com/filecoin-project/specs-actors/v4 v4.0.1 // indirect
	github.com/filecoin-project/specs-actors/v5 v5.0.4 // indirect
	github.com/filecoin-project/specs-actors/v6 v6.0.1 // indirect
	github.com/filecoin-project/specs-actors/v7 v7.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/ipfs/go-block-format v0.0.3 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.6 // indirect
	github.com/ipfs/go-ipld-format v0.2.0 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-log/v2 v2.5.0 // indirect
	github.com/ipsn/go-secp256k1 v0.0.0-20180726113642-9d62b9f0bc52 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/libp2p/go-buffer-pool v0.0.2 // indirect
	github.com/libp2p/go-flow-metrics v0.0.3 // indirect
	github.com/libp2p/go-libp2p-core v0.14.0 // indirect
	github.com/libp2p/go-openssl v0.0.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr v0.5.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.1.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20201211092308-30ac6d18308e // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)
