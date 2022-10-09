package exitcode

// Common error codes that may be shared by different actors.
// Actors may also define their own codes, including redefining these values.

const (
	// ErrIllegalArgument Indicates a method parameter is invalid.
	ErrIllegalArgument = FirstActorErrorCode + iota
	// ErrNotFound Indicates a requested resource does not exist.
	ErrNotFound
	// ErrForbidden Indicates an action is disallowed.
	ErrForbidden
	// ErrInsufficientFunds Indicates a balance of funds is insufficient.
	ErrInsufficientFunds
	// ErrIllegalState Indicates an actor's internal state is invalid.
	ErrIllegalState
	// ErrSerialization Indicates de/serialization failure within actor code.
	ErrSerialization
	// ErrUnhandledMessage Indicates the actor cannot handle this message.
	ErrUnhandledMessage
	// ErrUnspecified Indicates the actor failed with an unspecified error.
	ErrUnspecified
	// ErrAssertionFailed Indicates the actor failed a user-level assertion
	ErrAssertionFailed

	// Common error codes stop here.  If you define a common error code above
	// this value it will have conflicting interpretations
	FirstActorSpecificExitCode = ExitCode(32)
)
