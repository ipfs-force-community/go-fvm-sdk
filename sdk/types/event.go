package types

type Flags uint64

const (
	FLAGINDEXEDKEY   = 0b00000001
	FLAGINDEXEDVALUE = 0b00000010
	FLAGINDEXEDALL   = FLAGINDEXEDKEY | FLAGINDEXEDVALUE
)

// Entry A key value entry inside an Event.
type Entry struct {
	/// A bitmap conveying metadata or hints about this entry.
	Flags Flags
	/// The key of this event.
	Key string
	/// The value's codec. Must be IPLDRAW (0x55) for now according to FIP-0049.
	Codec Codec
	/// Any DAG-CBOR encodeable type.
	Value RawBytes
}

// ActorEvent An event as originally emitted by the actor.
type ActorEvent struct {
	Entries []*Entry
}
