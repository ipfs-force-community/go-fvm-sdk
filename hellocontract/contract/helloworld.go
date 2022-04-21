package contract

import (
	"bytes"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	cid "github.com/ipfs/go-cid"
)

type State struct {
	Count uint64
}

func (s *State) Save() cid.Cid {
	buf := bytes.NewBuffer([]byte{})
	err := s.MarshalCBOR(buf)
	if err != nil {
		sdk.Abort(20, fmt.Sprintf("failed to get root: %v", err))
	}
	stBytes := buf.Bytes()
	stCid, err := sdk.Put(0xb220, 32, types.DAG_CBOR, stBytes)
	if err != nil {
		sdk.Abort(20, fmt.Sprintf("failed to get root: %v", err))
	}

	err = sdk.SetRoot(stCid)
	if err != nil {
		sdk.Abort(20, fmt.Sprintf("failed to get root: %v", err))
	}
	return stCid
}

func NewState() *State {
	root, err := sdk.Root()
	if err != nil {
		sdk.Abort(20, fmt.Sprintf("failed to get root: %v", err))
	}

	data, err := sdk.Get(root)
	if err != nil {
		sdk.Abort(20, fmt.Sprintf("failed to get data: %v", err))
	}
	st := new(State)
	err = st.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		sdk.Abort(20, fmt.Sprintf("failed to get data: %v", err))
	}
	return st
}
