package props

// OnceProp is included on the first visit. Subsequent requests that list
// the prop key in X-Inertia-Except-Once-Props will exclude it, allowing
// the client to reuse the previously loaded value.
type OnceProp struct {
	Value any
}

// Once creates a prop that is resolved only on the first visit.
func Once(val any) OnceProp {
	return OnceProp{Value: val}
}

func (p OnceProp) Prop() any { return p.Value }
