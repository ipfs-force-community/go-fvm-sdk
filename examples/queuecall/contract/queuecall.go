package contract

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"

	"github.com/minio/blake2b-simd"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

const MIN_DELAY uint64 = 10      // seconds
const MAX_DELAY uint64 = 1000    // seconds
const GRACE_PERIOD uint64 = 1000 // seconds

const DEFAULTHAMTBITWIDTH = 3

type Call struct {
	To     address.Address
	Method abi.MethodNum
	Value  abi.TokenAmount
	Params []byte

	TimeStamp uint64
}

func (call Call) Key() (types.StringKey, []byte, error) {
	buf := bytes.NewBuffer(nil)
	err := call.MarshalCBOR(buf)
	if err != nil {
		return "", nil, err
	}

	data := buf.Bytes()
	hasher := blake2b.New256()
	hasher.Write(data)
	return types.StringKey(hex.EncodeToString(hasher.Sum(nil))), data, nil
}

type State struct {
	Owner  abi.ActorID
	Queues cid.Cid
}

func (e *State) Export() []interface{} {
	return []interface{}{
		Constructor,
		e.Queue,
		e.Execute,
	}
}

func Constructor(ctx context.Context) error {
	origin, err := sdk.Origin(ctx)
	if err != nil {
		return err
	}

	emptyMap, err := adt.MakeEmptyMap(adt.AdtStore(ctx), DEFAULTHAMTBITWIDTH)
	if err != nil {
		return err
	}

	emptyRoot, err := emptyMap.Root()
	if err != nil {
		return err
	}

	st := &State{
		Owner:  origin,
		Queues: emptyRoot,
	}
	_ = sdk.Constructor(ctx, st)
	return nil
}

func (e *State) Queue(ctx context.Context, call *Call) (types.CborString, error) {
	caller, err := sdk.Caller(ctx)
	if err != nil {
		return "", err
	}

	if caller != e.Owner {
		return "", errors.New("only onwner operate")
	}

	callMap, err := adt.AsMap(adt.AdtStore(ctx), e.Queues, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return "", err
	}

	timestamp, err := sdk.TipsetTimestamp(ctx)
	if err != nil {
		return "", err
	}

	// ---|------------|---------------|-------
	//  block    block + min     block + max
	if call.TimeStamp < timestamp+MIN_DELAY || call.TimeStamp > timestamp+MAX_DELAY {
		return "", fmt.Errorf("timestamp %d not in range (%d, %d)", call.TimeStamp, timestamp+MIN_DELAY, timestamp+MAX_DELAY)
	}

	key, data, err := call.Key()
	if err != nil {
		return "", err
	}

	has, err := callMap.Has(key)
	if err != nil {
		return "", err
	}
	if has {
		return "", fmt.Errorf("already queue this call")
	}

	err = callMap.Put(key, types.CBORBytes(data))
	if err != nil {
		return "", err
	}

	root, err := callMap.Root()
	if err != nil {
		return "", err
	}
	e.Queues = root
	sdk.SaveState(ctx, e)
	return types.CborString(key.Key()), nil
}

func (e *State) Execute(ctx context.Context, id *types.CborString) (*types.Receipt, error) {
	caller, err := sdk.Caller(ctx)
	if err != nil {
		return nil, err
	}

	if caller != e.Owner {
		return nil, err
	}

	callMap, err := adt.AsMap(adt.AdtStore(ctx), e.Queues, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, err
	}

	call := &Call{}
	found, err := callMap.Get(types.StringKey(*id), call)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("call %s not found", id)
	}

	timestamp, err := sdk.TipsetTimestamp(ctx)
	if err != nil {
		return nil, err
	}

	// ----|-------------------|-------
	//  timestamp    timestamp + grace period
	if timestamp < call.TimeStamp {
		return nil, fmt.Errorf("timestamp not passed block time(%d), call time (%d)", timestamp, call.TimeStamp)
	}

	if timestamp > call.TimeStamp+GRACE_PERIOD {
		return nil, fmt.Errorf("call time (%d) has expired %d", timestamp, call.TimeStamp)
	}

	return sdk.Send(ctx, call.To, call.Method, call.Params, call.Value)
}
