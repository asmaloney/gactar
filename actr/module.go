package actr

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	BufferInterface

	GetModuleName() string
}
