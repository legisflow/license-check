package main

import (
	"context"
	"fmt"

	"github.com/radiculaCZ/license-check/core"
	"github.com/radiculaCZ/license-check/languages/python"
	"github.com/radiculaCZ/license-check/results"
)

func main() {
	fileName := "requirements.txt"
	depFile, err := python.NewRequirementsTxt(fileName)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for dep := range depFile.GetDependencies(ctx) {
		println(dep)
	}

	packagesMeta, err := core.DownloadDependencyInfo(depFile, python.NewPyPI("pypi", python.PyPIURL))
	if err != nil {
		panic(err)
	}

	result := results.NewLicenseResult()
	result.AddPackageMeta(packagesMeta)

	fmt.Println(string(result.GetResult()))
}
