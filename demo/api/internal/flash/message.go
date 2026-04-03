package flash

// Message carries the demo app flash payload across packages and transport layers.
type Message struct {
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
