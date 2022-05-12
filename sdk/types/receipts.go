package types

type RawBytes []byte

//type ErrorNumber uint32
//
//const (
//	/// A syscall parameters was invalid.
//	IllegalArgument ErrorNumber = 1
//	/// The actor is not in the correct state to perform the requested operation.
//	IllegalOperation ErrorNumber = 2
//	/// This syscall would exceed some system limit (memory, lookback, call depth, etc.).
//	LimitExceeded ErrorNumber = 3
//	/// A system-level assertion has failed.
//	///
//	/// # Note
//	///
//	/// Non-system actors should never receive this error number. A system-level assertion will
//	/// cause the entire message to fail.
//	AssertionFailed ErrorNumber = 4
//	/// There were insufficient funds to complete the requested operation.
//	InsufficientFunds ErrorNumber = 5
//	/// A resource was not found.
//	NotFound ErrorNumber = 6
//	/// The specified IPLD block handle was invalid.
//	InvalidHandle ErrorNumber = 7
//	/// The requested CID shape (multihash codec, multihash length) isn't supported.
//	IllegalCid ErrorNumber = 8
//	/// The requested IPLD codec isn't supported.
//	IllegalCodec ErrorNumber = 9
//	/// The IPLD block did not match the specified IPLD codec.
//	Serialization ErrorNumber = 10
//	/// The operation is forbidden.
//	Forbidden ErrorNumber = 11
//)

type Receipt struct {
	ExitCode   uint32   `json:"exit_code"`
	ReturnData RawBytes `json:"return_data"`
	GasUsed    int64    `json:"gas_used"`
}
