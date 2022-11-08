// Package ferrors fvm errors
package ferrors

import (
	"fmt"
)

// ExitCode define error in fvm and custom actor
type ExitCode uint32

// Error return error message of exitcode
func (e ExitCode) Error() string {
	return fmt.Sprintf("%d", e)
}

// Is check whether error is exitcode
func (e ExitCode) Is(code error) bool {
	return e == code
}

// IsSystemError Returns true if the error code is in the range of exit codes reserved for the VM (including Ok).
func (e ExitCode) IsSystemError() bool {
	return uint32(e) < FIRST_USER_EXIT_CODE
}

// nolint
const (
	// Exit codes which originate inside the VM.
	// These values may not be used by actors when aborting.

	// The code indicating successful execution.
	OK ExitCode = 0
	// Indicates the message sender doesn't exist.
	SYS_SENDER_INVALID ExitCode = 1
	// Indicates that the message sender was not in a valid state to send this message.
	// Either:
	// - The sender's nonce nonce didn't match the message nonce.
	// - The sender didn't have the funds to cover the message gas.
	SYS_SENDER_STATE_INVALID ExitCode = 2
	// Indicates failure to find a method in an actor.
	SYS_INVALID_METHOD ExitCode = 3 // FIXME: reserved
	// Indicates the message receiver trapped (panicked).
	SYS_ILLEGAL_INSTRUCTION ExitCode = 4
	// Indicates the message receiver doesn't exist and can't be automatically created
	SYS_INVALID_RECEIVER ExitCode = 5
	// Indicates the message sender didn't have the requisite funds.
	SYS_INSUFFICIENT_FUNDS ExitCode = 6
	// Indicates message execution (including subcalls) used more gas than the specified limit.
	SYS_OUT_OF_GAS ExitCode = 7
	// SYS_RESERVED_8 ExitCode = ExitCode::new(8);
	// Indicates the message receiver aborted with a reserved exit code.
	SYS_ILLEGAL_EXIT_CODE ExitCode = 9
	// Indicates an internal VM assertion failed.
	SYS_ASSERTION_FAILED ExitCode = 10
	// Indicates the actor returned a block handle that doesn't exist
	SYS_MISSING_RETURN ExitCode = 11
	// SYS_RESERVED_12 ExitCode = ExitCode::new(12);
	// SYS_RESERVED_13 ExitCode = ExitCode::new(13);
	// SYS_RESERVED_14 ExitCode = ExitCode::new(14);
	// SYS_RESERVED_15 ExitCode = ExitCode::new(15);

	// The lowest exit code that an actor may abort with.
	FIRST_USER_EXIT_CODE uint32 = 16

	// Standard exit codes according to the built-in actors' calling convention.
	// Indicates a method parameter is invalid.
	USR_ILLEGAL_ARGUMENT ExitCode = 16
	// Indicates a requested resource does not exist.
	USR_NOT_FOUND ExitCode = 17
	// Indicates an action is disallowed.
	USR_FORBIDDEN ExitCode = 18
	// Indicates a balance of funds is insufficient.
	USR_INSUFFICIENT_FUNDS ExitCode = 19
	// Indicates an actor's internal state is invalid.
	USR_ILLEGAL_STATE ExitCode = 20
	// Indicates de/serialization failure within actor code.
	USR_SERIALIZATION ExitCode = 21
	// Indicates the actor cannot handle this message.
	USR_UNHANDLED_MESSAGE ExitCode = 22
	// Indicates the actor failed with an unspecified error.
	USR_UNSPECIFIED ExitCode = 23
	// Indicates the actor failed a user-level assertion
	USR_ASSERTION_FAILED ExitCode = 24
	// RESERVED_25 ExitCode = 25
	// RESERVED_26 ExitCode = 26
	// RESERVED_27 ExitCode = 27
	// RESERVED_28 ExitCode = 28
	// RESERVED_29 ExitCode = 29
	// RESERVED_30 ExitCode = 30
	// RESERVED_31 ExitCode = 31
)

// ErrorNumber  When a syscall fails, it returns an `ErrorNumber` to indicate why
type ErrorNumber uint32

// Error return error message of ErrorNumber
func (e ErrorNumber) Error() string {
	return EnToString(uint32(e))
}

// EnToString return string message of ErrorNumber
func EnToString(code uint32) string {
	switch code {
	case 1:
		return "illegal argument"
	case 2:
		return "illegal operation"
	case 3:
		return "limit exceeded"
	case 4:
		return "filecoin assertion failed"
	case 5:
		return "insufficient funds"
	case 6:
		return "resource not found"
	case 7:
		return "invalid ipld block handle"
	case 8:
		return "illegal cid specification"
	case 9:
		return "illegal ipld codec"
	case 10:
		return "serialization error"
	case 11:
		return "operation forbidden"
	case 12:
		return "buffer too small"
	}
	return "other error"
}

// Is check whether error is ErrorNumber
func (e ErrorNumber) Is(code error) bool {
	return e == code
}

// nolint
const (
	// A syscall parameters was invalid.
	IllegalArgument ErrorNumber = 1
	// The actor is not in the correct state to perform the requested operation.
	IllegalOperation ErrorNumber = 2
	// This syscall would exceed some system limit (memory, lookback, call depth, etc.).
	LimitExceeded ErrorNumber = 3
	// A system-level assertion has failed.
	//
	// # Note
	//
	// Non-system actors should never receive this error number. A system-level assertion will
	// cause the entire message to fail.
	AssertionFailed ErrorNumber = 4
	// There were insufficient funds to complete the requested operation.
	InsufficientFunds ErrorNumber = 5
	// A resource was not found.
	NotFound ErrorNumber = 6
	// The specified IPLD block handle was invalid.
	InvalidHandle ErrorNumber = 7
	// The requested CID shape (multihash codec, multihash length) isn't supported.
	IllegalCid ErrorNumber = 8
	// The requested IPLD codec isn't supported.
	IllegalCodec ErrorNumber = 9
	// The IPLD block did not match the specified IPLD codec.
	Serialization ErrorNumber = 10
	// The operation is forbidden.
	Forbidden ErrorNumber = 11
	// The passed buffer is too small.
	BufferTooSmall ErrorNumber = 12
)

// SysCallError Fvm error number include error code and error message
type SysCallError struct {
	code    ErrorNumber
	message string
}

// NewSysCallError new fvm error number from error code and message
func NewSysCallError(code ErrorNumber, msg string) SysCallError {
	return SysCallError{code, msg}
}

// Error return error message for fvm error number
func (e SysCallError) Error() string {
	return fmt.Sprintf("%s %d", e.message, e.code)
}

// Unwrap return inner error code
func (e SysCallError) Unwrap() error {
	return e.code
}
