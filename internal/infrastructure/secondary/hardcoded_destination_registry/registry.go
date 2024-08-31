package hardcoded_destination_registry

const (
	Mars         = "MARS"
	Moon         = "MOON"
	Pluto        = "PLUTO"
	AsteroidBelt = "ASTEROID_BELT"
	Europa       = "EUROPA"
	Titan        = "TITAN"
	Ganymede     = "GANYMEDE"
)

func New() *Registry {
	return &Registry{
		destinations: map[string]struct{}{
			Mars:         {},
			Moon:         {},
			Pluto:        {},
			AsteroidBelt: {},
			Europa:       {},
			Titan:        {},
			Ganymede:     {},
		},
	}
}

type Registry struct {
	destinations map[string]struct{}
}

func (r *Registry) DestinationExists(destinationID string) bool {
	_, ok := r.destinations[destinationID]
	return ok
}
