module gen

go 1.17

require (
	github.com/ipfs-force-community/go-fvm-sdk/gen v0.0.0-00010101000000-000000000000
	hellocontract v0.0.0-00010101000000-000000000000
)

replace (
	github.com/ipfs-force-community/go-fvm-sdk => ../../..
	github.com/ipfs-force-community/go-fvm-sdk/gen => ../../../gen
	hellocontract => ../
)

require (
	github.com/filecoin-project/go-address v0.0.6 // indirect
	github.com/filecoin-project/go-amt-ipld/v4 v4.0.0 // indirect
	github.com/filecoin-project/go-bitfield v0.2.4 // indirect
	github.com/filecoin-project/go-commp-utils v0.1.3 // indirect
	github.com/filecoin-project/go-commp-utils/nonffi v0.0.0-20220905160352-62059082a837 // indirect
	github.com/filecoin-project/go-fil-commcid v0.1.0 // indirect
	github.com/filecoin-project/go-hamt-ipld/v3 v3.1.0 // indirect
	github.com/filecoin-project/go-state-types v0.1.12-alpha // indirect
	github.com/filecoin-project/specs-actors v0.9.14 // indirect
	github.com/filecoin-project/specs-actors/v2 v2.3.6 // indirect
	github.com/filecoin-project/specs-actors/v5 v5.0.4 // indirect
	github.com/filecoin-project/specs-actors/v7 v7.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/ipfs-force-community/go-fvm-sdk v0.0.0-00010101000000-000000000000 // indirect
	github.com/ipfs/go-block-format v0.0.3 // indirect
	github.com/ipfs/go-cid v0.2.0 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.6 // indirect
	github.com/ipfs/go-ipld-format v0.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.1.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/polydawn/refmt v0.0.0-20201211092308-30ac6d18308e // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20220514204315-f29c37e9c44c // indirect
	golang.org/x/crypto v0.0.0-20211209193657-4570a0811e8b // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/tools v0.1.11 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)
