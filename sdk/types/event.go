package types

type Flags uint8

const (
	FLAG_INDEXED_KEY   = 0b00000001
	FLAG_INDEXED_VALUE = 0b00000010
	FLAG_INDEXED_ALL   = FLAG_INDEXED_KEY | FLAG_INDEXED_VALUE
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
