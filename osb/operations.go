package osb

const (
	// Provisioning represents the "provisioning" operation
	Provisioning = "provisioning"
	// Updating represents the "updating" operation
	Updating = "updating"
	// Deprovisioning represents the "deprovisioning" operation
	Deprovisioning = "deprovisioning"
	// Binding represents the "binding" operation
	Binding = "binding"
	// Unbinding represents the "unbinding" operation
	Unbinding = "unbinding"
)

type ProvisioningState int

const (
	// StateInProgress represents the state of an operation that is still
	// pending completion
	StateInProgress ProvisioningState = iota
	// StateSucceeded represents the state of an operation that has
	// completed successfully
	StateSucceeded
	// StateFailed represents the state of an operation that has failed
	StateFailed
	// StateGone is a pseudo oepration state represting the "state"
	// of an operation against an entity that no longer exists
	StateGone
)

func (s ProvisioningState) String() string {
	switch s {
	case StateInProgress:
		return "in progress"
	case StateSucceeded:
		return "succeeded"
	case StateFailed:
		return "failed"
	case StateGone:
		return "gone"
	default:
		return "Unknown"
	}
}
