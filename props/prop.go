package props

// Proper is implemented by prop type wrappers to provide their
// underlying value during resolution.
type Proper interface {
	Prop() any
}

// TryProper extends Proper for values whose resolution may fail.
type TryProper interface {
	TryProp() (any, error)
}

// resolve evaluates a raw prop value. It unwraps Proper/TryProper
// wrappers first, then evaluates func() any / func() (any, error)
// lazy values.
func resolve(val any) (any, error) {
	// Unwrap prop type wrappers to get the inner value.
	switch v := val.(type) {
	case TryProper:
		return v.TryProp()
	case Proper:
		val = v.Prop()
	}

	// Evaluate lazy function values.
	switch v := val.(type) {
	case func() (any, error):
		return v()
	case func() any:
		return v(), nil
	default:
		return val, nil
	}
}
