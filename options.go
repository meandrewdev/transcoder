package transcoder

// Options ...
type Options interface {
	GetStrArguments() []string
	GetInputs() []string
}
