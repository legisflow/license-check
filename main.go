package main

import (
	"fmt"
	"os"

	"github.com/radiculaCZ/license-check/core"
	"github.com/radiculaCZ/license-check/interfaces"
	"github.com/radiculaCZ/license-check/languages/python"
	"github.com/radiculaCZ/license-check/results"
	"github.com/urfave/cli/v2"
)

func main() {
	// Register all depfiles
	requrementsTxt := python.NewRequirementsTxt()

	depFiles := map[string]interfaces.DepFile{
		requrementsTxt.GetDepFileType(): requrementsTxt,
	}
	// End of registering depfiles

	// Register results
	licenseResult := results.NewLicenseResult()

	results := map[string]interfaces.Result{
		licenseResult.GetResultName(): licenseResult,
	}
	// End of registering results

	typeFlag := &cli.StringFlag{
		Name:     "type",
		Aliases:  []string{"t"},
		Required: true,
		Action: func(c *cli.Context, v string) error {
			if _, ok := depFiles[v]; !ok {
				return fmt.Errorf("Invalid type %s", v)
			}
			return nil
		},
	}

	typeFlag.Usage = "Allowed values: " + func() string {
		var types []string
		types = append(types, "\n")
		for k := range depFiles {
			types = append(types, k+"\n")
		}
		return fmt.Sprintf("%v", types)
	}()

	resultFlag := &cli.StringFlag{
		Name:     "result",
		Aliases:  []string{"r"},
		Required: true,
		Action: func(c *cli.Context, v string) error {
			if _, ok := results[v]; !ok {
				return fmt.Errorf("Invalid result %s", v)
			}
			return nil
		},
	}

	resultFlag.Usage = "Allowed values: " + func() string {
		var types []string
		types = append(types, "\n")
		for k := range results {
			types = append(types, k+"\n")
		}
		return fmt.Sprintf("%v", types)
	}()

	app := cli.App{
		Name:    "license-check",
		Version: "0.0.1",
		Usage:   "Check license of dependencies by passing the dependency file",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:      "file",
				Aliases:   []string{"f"},
				Usage:     "Path to the dependency file",
				Required:  true,
				TakesFile: true,
				Action: func(c *cli.Context, p cli.Path) error {
					if _, err := os.Stat(p); os.IsNotExist(err) {
						return fmt.Errorf("File %s does not exist", p)
					}
					return nil
				},
			},
			typeFlag,
			resultFlag,
		},
		Action: func(c *cli.Context) error {
			packages, err := core.DownloadDependencyInfo(depFiles[c.String("type")], c.Path("file"))
			if err != nil {
				return err
			}

			result := core.ProcessDependencyInfo(results[c.String("result")], packages)
			fmt.Println(result)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
