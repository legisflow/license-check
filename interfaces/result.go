// The result defines the interface for output formats
package interfaces

// Result represents the output format
// Data can only be appended as this is used only
// for output representation.
// Any intermediate processing can be done directly on data structures
type Result interface {
	GetResultName() string
	GetResult() []byte
	AddPackageMeta([]PackageMeta)
}
