//nolint:unparam
package simulated

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"
)

func blakehash(data []byte) []byte {
	blake := blake2b.New256()
	return blake.Sum(data)
}

type emptyInterface struct {
	_    uintptr
	word unsafe.Pointer
}

func GetSlicePointerAndLen(data interface{}) (uintptr, uint32) {
	eface := (*emptyInterface)(unsafe.Pointer(&data))
	s := (*reflect.SliceHeader)(eface.word)
	return s.Data, uint32(s.Len)
}

func GetStringPointerAndLen(str string) (uintptr, uint32) {
	s := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return s.Data, uint32(s.Len)
}

// Generate a hash of length 32 bytes
func makeRandomness(dst int64, round int64, entropy []byte) []byte {
	dstbyte := [8]byte{}
	binary.BigEndian.PutUint64(dstbyte[0:8], abs(dst))
	roundbyte := [8]byte{}
	binary.BigEndian.PutUint64(roundbyte[0:8], abs(round))
	entropy = append(entropy, dstbyte[:]...)
	entropy = append(entropy, roundbyte[:]...)
	return blakehash(entropy)
}

func abs(i int64) uint64 {
	if i < 0 {
		return uint64(-i)
	}
	return uint64(i)
}

var (
	AccountCid          = mustParseCid("bafk2bzacebmfbtdj5vruje5auacrhhprcjdd6uclhukb7je7t2f6ozfcgqlu2")
	CronCid             = mustParseCid("bafk2bzacea4gwsbeux7z4yxvpkxpco77iyxijoyqaoikofrxdewunwh3unjem")
	InitCid             = mustParseCid("bafk2bzacebwkqd6e7gdphfzw2kdmbokdh2bly6fvzgfopxzy7quq4l67gmkks")
	MultisigCid         = mustParseCid("bafk2bzacea5zp2g6ag5qfuro7zw6kyku2swxs57wjxncaaxbih5iqflqy4ghm")
	PaymentchannelCid   = mustParseCid("bafk2bzaced47dbtbygmfwnyfsp5iihzhhdmnkpuyc5nlnfgc4mkkvlsgvj2do")
	RewardCid           = mustParseCid("bafk2bzacecmcagk32pzdzfg7piobzqhlgla37x3g7jjzyndlz7mqdno2zulfi")
	StoragemarketCid    = mustParseCid("bafk2bzacecxqgajcaednamgolc6wc3lzbjc6tz5alfrbwqez2y3c372vts6dg")
	StorageminerCid     = mustParseCid("bafk2bzaceaqwxllfycpq6decpsnkqjdeycpysh5acubonjae7u3wciydlkvki")
	StoragepowerCid     = mustParseCid("bafk2bzaceddmeolsokbxgcr25cuf2skrobtmmoof3dmqfpcfp33lmw63oikvm")
	SystemCid           = mustParseCid("bafk2bzaced6kjkbv7lrb2qwq5we2hqaxc6ztch5p52g27qtjy45zdemsk4b7m")
	VerifiedRegistryCid = mustParseCid("bafk2bzacectzxvtoselhnzsair5nv6k5vokvegnht6z2lfee4p3xexo4kg4m6")
)

var EmbeddedBuiltinActors = map[string]cid.Cid{
	"account":          mustParseCid("bafk2bzacebmfbtdj5vruje5auacrhhprcjdd6uclhukb7je7t2f6ozfcgqlu2"),
	"cron":             mustParseCid("bafk2bzacea4gwsbeux7z4yxvpkxpco77iyxijoyqaoikofrxdewunwh3unjem"),
	"init":             mustParseCid("bafk2bzacebwkqd6e7gdphfzw2kdmbokdh2bly6fvzgfopxzy7quq4l67gmkks"),
	"multisig":         mustParseCid("bafk2bzacea5zp2g6ag5qfuro7zw6kyku2swxs57wjxncaaxbih5iqflqy4ghm"),
	"paymentchannel":   mustParseCid("bafk2bzaced47dbtbygmfwnyfsp5iihzhhdmnkpuyc5nlnfgc4mkkvlsgvj2do"),
	"reward":           mustParseCid("bafk2bzacecmcagk32pzdzfg7piobzqhlgla37x3g7jjzyndlz7mqdno2zulfi"),
	"storagemarket":    mustParseCid("bafk2bzacecxqgajcaednamgolc6wc3lzbjc6tz5alfrbwqez2y3c372vts6dg"),
	"storageminer":     mustParseCid("bafk2bzaceaqwxllfycpq6decpsnkqjdeycpysh5acubonjae7u3wciydlkvki"),
	"storagepower":     mustParseCid("bafk2bzaceddmeolsokbxgcr25cuf2skrobtmmoof3dmqfpcfp33lmw63oikvm"),
	"system":           mustParseCid("bafk2bzaced6kjkbv7lrb2qwq5we2hqaxc6ztch5p52g27qtjy45zdemsk4b7m"),
	"verifiedregistry": mustParseCid("bafk2bzacectzxvtoselhnzsair5nv6k5vokvegnht6z2lfee4p3xexo4kg4m6"),
}

func mustParseCid(c string) cid.Cid {
	ret, err := cid.Decode(c)
	if err != nil {
		panic(err)
	}

	return ret
}

func actorTypeTostring(actorT types.ActorType) (string, error) {
	switch actorT {
	case types.System:
		return "system", nil
	case types.Init:
		return "init", nil
	case types.Cron:
		return "cron", nil
	case types.Account:
		return "account", nil
	case types.Power:
		return "power", nil
	case types.Miner:
		return "miner", nil
	case types.PaymentChannel:
		return "paymentchannel", nil
	case types.Multisig:
		return "multisig", nil
	case types.Reward:
		return "reward", nil
	case types.VerifiedRegistry:
		return "verifiedregistry", nil
	default:
		return "", ErrorNotFound
	}

}

func stringToactorType(str string) (actorT types.ActorType, err error) {
	switch str {
	case "system":
		return types.System, nil
	case "init":
		return types.Init, nil
	case "cron":
		return types.Cron, nil
	case "account":
		return types.Account, nil
	case "power":
		return types.Power, nil
	case "miner":
		return types.Miner, nil
	case "paymentchannel":
		return types.PaymentChannel, nil
	case "multisig":
		return types.Multisig, nil
	case "reward":
		return types.Reward, nil
	case "verifiedregistry":
		return types.VerifiedRegistry, nil
	default:
		return types.ActorType(0), ErrorNotFound
	}

}
