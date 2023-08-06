// Used to get packages from PyPI (pypi.org)
// Uses the PyPI JSON API https://warehouse.pypa.io/api-reference/json.html
// using the url https://pypi.org/pypi/<package_name>/json
// The PyPI represent an implementation of interface  packagerepository
package python

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/radiculaCZ/license-check/interfaces"
)

const PyPIURL = "https://pypi.org/pypi/<package_name>/json"

type PyPI struct {
	name string
	url  string
}

type pypiResponse struct {
	// PyPI general package information
	Info pypiInfo `json:"info"`
}

// pypiInfo contains the general package information
// it should capture at least what is covered in PackageMeta
type pypiInfo struct {
	Author  string `json:"author,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	// Contains the pypi package classifiers, one of which is license
	Classifiers []string   `json:"classifiers"`
	Description string     `json:"description"`
	Homepage    string     `json:"home_page,omitempty"`
	ProjectURL  projectURL `json:"project_urls"`
}

type projectURL struct {
	Homepage string `json:"Homepage,omitempty"`
	Source   string `json:"Source,omitempty"`
}

// NewPyPI creates a new instance of the PyPI struct
// Nothing crazy, just struct init, no call to the API
func NewPyPI(name string, pypiURL string) interfaces.PackageRepository {
	return &PyPI{
		name: name,
		url:  pypiURL,
	}
}

func (p *PyPI) GetRepositoryName() string {
	return p.name
}

func (p *PyPI) GetPackageInfo(packageName string) (*interfaces.PackageMeta, error) {
	url := strings.Replace(p.url, "<package_name>", packageName, 1)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response pypiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return convertPyPIResponseToPackageMeta(response), nil
}

// convertPyPIResponseToPackageMeta converts the PyPI response to the PackageMeta struct
// this shields the changes in the PyPI API from the rest of the code
func convertPyPIResponseToPackageMeta(response pypiResponse) *interfaces.PackageMeta {
	meta := &interfaces.PackageMeta{}

	meta.Author = response.Info.Author
	meta.Name = response.Info.Name
	meta.Version = response.Info.Version
	meta.License = extractLicenseFromClassifiers(response.Info.Classifiers)
	meta.Description = response.Info.Description
	meta.Homepage = response.Info.ProjectURL.Homepage
	meta.Repository = response.Info.ProjectURL.Source
	meta.Language = "python"

	return meta
}

// extractLicenseFromClassifiers extracts the license from the classifier list
// walks over classifiers and return the license, which is the last item in the item
// starting with 'license ::'
// if the license cannot be found return empty string
func extractLicenseFromClassifiers(classifiers []string) string {
	for _, classifier := range classifiers {
		if strings.HasPrefix(strings.ToLower(classifier), "license ::") {
			tempLicense := strings.Split(classifier, "::")
			return strings.TrimSpace(tempLicense[len(tempLicense)-1])
		}
	}
	return ""
}
