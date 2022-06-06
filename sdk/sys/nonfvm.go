//go:build !tinygo.wasm
// +build !tinygo.wasm

package sys

func vmAbort(code uint32, msgOff uintptr, msgLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func actorResolveAddress(ret uintptr, addr_off uintptr, addr_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func actorGetActorCodeCid(ret uintptr, addr_off uintptr, addr_len uint32, obuf_off uintptr, obuf_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func actorResolveBuiltinActorType(ret uintptr, cid_off uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func actorGetCodeCidForType(ret uintptr, typ int32, obuf_off uintptr, obuf_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func actorNewActorAddress(ret uintptr, obuf_off uintptr, obuf_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func actorCreateActor(actor_id uint64, typ_off uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }


func vmContext(ret uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }


func cryptoVerifySignature(ret uintptr, sigOff uintptr, sigLen uint32, addrOff uintptr, addrLen uint32, plainTextOff uintptr, plainTextLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoHashBlake2b(ret uintptr, dataOff uintptr, dataLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoComputeUnsealedSectorCid(ret uintptr, proofType int64, piecesOff uintptr, pieceLen uint32, cidPtr uintptr, cidLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoVerifySeal(ret uintptr, infoOff uintptr, infoLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoVerifyPost(ret uintptr, infoOff uintptr, infoLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoVerifyConsensusFault(ret uintptr, h1Off uintptr, h1Len uint32, h2Off uintptr, h2Len uint32, extraOff uintptr, extraLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoVerifyAggregateSeals(ret uintptr, aggOff uintptr, aggLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoVerifyReplicaUpdate(ret uintptr, repOff uintptr, repLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func cryptoBatchVerifySeals(batchOff uintptr, batLen uint32, resultOff uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }


func debugEnabled(ret uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func debugLog(message uintptr, message_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func gasCharge(name_off uintptr, name_len uint32, amount uint64) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func ipldOpen(ret uintptr, cid uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func ipldCreate(ret uintptr, codec uint64, data uintptr, len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func ipldRead(ret uintptr, id uint32, offset uint32, obuf uintptr, max_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func ipldStat(ret uintptr, id uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func ipldLink(ret uintptr, id uint32, hash_fun uint64, hash_len uint32, cid uintptr, cid_max_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func networkBaseFee(ret uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func networkTotalFilCircSupply(ret uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func getChainRandomness(ret uintptr, tag int64, epoch int64, entropy_off uintptr, entropy_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func getBeaconRandomness(ret uintptr, tag int64, epoch int64, entropy_off uintptr, entropy_len uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func sysSend(ret uintptr, recipient_off uintptr, recipient_len uint32, method uint64, params uint32, value_hi uint64, value_lo uint64) uint32 { panic("ignore this error, just implement nonfvm for ide working") }

func sselfRoot(ret uintptr, cid uintptr, cidMaxLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func sselfSetRoot(cid uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func selfCurrentBalance(ret uintptr) uint32 { panic("ignore this error, just implement nonfvm for ide working") }
func selfDestruct(addrOff uintptr, addrLen uint32) uint32 { panic("ignore this error, just implement nonfvm for ide working") }