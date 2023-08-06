package interfaces

// PackageMeta represents a programming langugage package
// It contains the meta information about the package
// PackageMeta does not have to be an interface because
// it will not change and will not require different
// implementations
type PackageMeta struct {
	Author      string
	Name        string
	Version     string
	License     string
	Description string
	Homepage    string
	Repository  string
	Language    string
}
