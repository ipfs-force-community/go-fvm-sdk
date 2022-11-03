//go:build !tinygo.wasm
// +build !tinygo.wasm

package sys

func vmAbort(code uint32, msgOff uintptr, msgLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func actorResolveAddress(ret uintptr, addrOff uintptr, addrLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func actorGetActorCodeCid(ret uintptr, actorID uint64, oBufOff uintptr, oBufLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func actorResolveBuiltinActorType(ret uintptr, cidOff uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func actorGetCodeCidForType(ret uintptr, typ int32, oBufOff uintptr, oBufLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func actorNewActorAddress(ret uintptr, oBufOff uintptr, oBufLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func actorCreateActor(actorID uint64, typOff uintptr, predictableAddrOff uintptr, predictableAddrLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func vmContext(ret uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func cryptoVerifySignature(ret uintptr, sigType uint32, sigOff uintptr, sigLen uint32, addrOff uintptr, addrLen uint32, plainTextOff uintptr, plainTextLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoHashBlake2b(ret uintptr, hashCode uint64, dataOff uintptr, dataLen uint32, digestOff uintptr, digestLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoComputeUnsealedSectorCid(ret uintptr, proofType int64, piecesOff uintptr, pieceLen uint32, cidPtr uintptr, cidLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoVerifySeal(ret uintptr, infoOff uintptr, infoLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoVerifyPost(ret uintptr, infoOff uintptr, infoLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoVerifyConsensusFault(ret uintptr, h1Off uintptr, h1Len uint32, h2Off uintptr, h2Len uint32, extraOff uintptr, extraLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoVerifyAggregateSeals(ret uintptr, aggOff uintptr, aggLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoVerifyReplicaUpdate(ret uintptr, repOff uintptr, repLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func cryptoBatchVerifySeals(batchOff uintptr, batLen uint32, resultOff uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func debugEnabled(ret uintptr) uint32 {
	return 0 //log usually used as global variable, stop panic when use this package
}
func debugLog(message uintptr, messageLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func gasCharge(nameOff uintptr, nameLen uint32, amount uint64) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func ipldOpen(ret uintptr, cid uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func ipldCreate(ret uintptr, codec uint64, data uintptr, len uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func ipldRead(ret uintptr, id uint32, offset uint32, obuf uintptr, maxLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func ipldStat(ret uintptr, id uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func ipldLink(ret uintptr, id uint32, hashFun uint64, hashLen uint32, cid uintptr, cidMaxLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func networkBaseFee(ret uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func networkTotalFilCircSupply(ret uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func getChainRandomness(ret uintptr, tag int64, epoch int64, entropyOff uintptr, entryopyLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func getBeaconRandomness(ret uintptr, tag int64, epoch int64, entropyOff uintptr, entryopyLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func sysSend(ret uintptr, recipientOff uintptr, recipientLen uint32, method uint64, params uint32, valueHI uint64, valueLow uint64) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func sselfRoot(ret uintptr, cid uintptr, cidMaxLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func sselfSetRoot(cid uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func selfCurrentBalance(ret uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func selfDestruct(addrOff uintptr, addrLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func actorLookupAddress(ret uintptr, actorID uint64, addrBufOff uintptr, addrBufLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func actorBalanceOf(ret uintptr, actorID uint64) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func debugStoreArtifact(nameOff uintptr, nameLen uint32, dataOff uintptr, dataLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func gasAvailable(ret uintptr) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
func networkTipsetTimestamp(ret uintptr) uint64 {
	panic("ignore this error, just implement nonfvm for ide working")
}

func networkTipsetCid(ret uintptr, epoch uint64, retOff uintptr, retLen uint32) uint32 {
	panic("ignore this error, just implement nonfvm for ide working")
}
