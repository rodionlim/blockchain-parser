package types

type Tracker struct {
	FirstParsedBlock uint64 // record keeping to notify from which block we started subscribing from
	Offset           int    // keep a tracker of the last transaction that was published to the subscriber
}
