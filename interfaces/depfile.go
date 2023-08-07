// The interfaces package contains all interfaces used by the application
// It does not define any implementation specific, as this needs to be done
// based on a specific programming language
// The constructors will also be implemented based on the programming language
package interfaces

import "context"

// DepFile represents a programming language dependency file
// They contain a list of packages (libraries) that are required
type DepFile interface {
	// GetDepFileType returns the type of the dependency file
	// e.g. npm, composer, pip, etc.
	GetDepFileType() string
	// GetDependencies returns a channel of package names
	// The path to the depfile is passed as second argument
	// The channel is closed when all package names are sent
	// The context is passed in a case the function needs to be cancelled
	GetDependencies(context.Context, string) (<-chan string, error)
	// GetRepository returns the repository that is used to get package information
	// Ussually the repository is language specific, but it might be also
	// file format specific (altough I am not aware of any such case)
	GetRepository() PackageRepository
}
