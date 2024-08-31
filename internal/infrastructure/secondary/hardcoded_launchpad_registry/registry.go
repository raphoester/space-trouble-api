package hardcoded_launchpad_registry

const (
	Texas      = "TEXAS"
	Florida    = "FLORIDA"
	California = "CALIFORNIA"
)

func New() *Registry {
	return &Registry{
		launchpads: map[string]struct{}{
			Texas:      {},
			Florida:    {},
			California: {},
		},
	}
}

type Registry struct {
	launchpads map[string]struct{}
}

func (r *Registry) LaunchpadExists(launchpadID string) bool {
	_, ok := r.launchpads[launchpadID]
	return ok
}
