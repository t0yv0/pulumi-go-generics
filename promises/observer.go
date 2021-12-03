package promises

// Collects events from the promise system to allow out-of-band error
// detection and idle state tracking.
type Observer interface {
	Created()
	Rejected(error)
	Resolved()
}
