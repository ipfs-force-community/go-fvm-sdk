package types

import cid "github.com/ipfs/go-cid"

type InstallParams struct {
	Code []byte
}

type InstallReturn struct {
	CodeCid   cid.Cid
	Installed bool
}
