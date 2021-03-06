package rich

// three-valued switch {unset, off, on}
type Tristate byte

const (
	Unset Tristate = iota
	Off
	On
)
