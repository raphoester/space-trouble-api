package hardcoded_launchpad_registry

const (
	Texas      = "TEXAS"
	Florida    = "FLORIDA"
	California = "CALIFORNIA"
	Other      = "5e9e4501f509094ba4566f84" // hack to work with spaceX API
)

func New() *Registry {
	return &Registry{
		launchpads: map[string]struct{}{
			Texas:      {},
			Florida:    {},
			California: {},
			Other:      {},
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
