// PackageRepository represent a package index where packages are stored
// on the internet
package interfaces

// PackageRepository represents an online package index
// It is used to search for packages and their versions and licenses
// some also provide CVE checks
type PackageRepository interface {
	// GetRepositoryName returns the name of the repository
	GetRepositoryName() string
	// GetPackageInfo returns the meta information about a package
	// If nothing cannot be returned, the nil is returned and error
	GetPackageInfo(packageName string) (*PackageMeta, error)
}
