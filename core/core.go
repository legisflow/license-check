// Core contains the main logic for pulling package info from the the package index,
// processing the downloaded data, and storing the result
package core

import (
	"context"

	"github.com/radiculaCZ/license-check/interfaces"
)

func DownloadDependencyInfo(depFile interfaces.DepFile, file string) ([]interfaces.PackageMeta, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	packages := make([]interfaces.PackageMeta, 0)

	repo := depFile.GetRepository()

	// iterate over the list of dependencies and download their meta info
	dependencies, err := depFile.GetDependencies(ctx, file)
	if err != nil {
		return nil, err
	}
	for dep := range dependencies {
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

func ProcessDependencyInfo(result interfaces.Result, packages []interfaces.PackageMeta) string {
	result.AddPackageMeta(packages)
	return string(result.GetResult())
}
