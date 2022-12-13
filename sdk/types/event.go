package types

type Flags uint8

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
	/// Any DAG-CBOR encodeable type.
	Value RawBytes
}

// ActorEvent An event as originally emitted by the actor.
type ActorEvent struct {
	Entries []Entry
}
