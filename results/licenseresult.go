// Implements the Result interface for a license result
// The output for this implementation is
// <package name>: <license name>
// <package name>: <license name>
package results

import "github.com/radiculaCZ/license-check/interfaces"

type LicenseResult struct {
	name     string
	packages []interfaces.PackageMeta
}

func NewLicenseResult() *LicenseResult {
	return &LicenseResult{
		name:     "License",
		packages: []interfaces.PackageMeta{},
	}
}

func (l *LicenseResult) AddPackageMeta(packages []interfaces.PackageMeta) {
	l.packages = append(l.packages, packages...)
}

func (l *LicenseResult) GetResultName() string {
	return l.name
}

func (l *LicenseResult) GetResult() []byte {
	var result []byte
	for _, pkg := range l.packages {
		result = append(result, []byte(pkg.Name+": "+pkg.License+"\n")...)
	}
	return result
}
