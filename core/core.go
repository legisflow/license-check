// Core contains the main logic for pulling package info from the the package index,
// processing the downloaded data, and storing the result
package core

import (
	"context"

	"github.com/radiculaCZ/license-check/interfaces"
)

func DownloadDependencyInfo(depFile interfaces.DepFile, repo interfaces.PackageRepository) ([]interfaces.PackageMeta, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	packages := make([]interfaces.PackageMeta, 0)

	// iterate over the list of dependencies and download their meta info
	for dep := range depFile.GetDependencies(ctx) {
		// download the meta info
		meta, err := repo.GetPackageInfo(dep)
		if err != nil {
			return nil, err
		}
		// store the meta info
		packages = append(packages, *meta)
	}

	return packages, nil
}
