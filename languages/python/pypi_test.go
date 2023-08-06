package python

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestExtractLicenseFromClassifiers tests the extractLicenseFromClassifiers function
// It does not cover every possible license, but it should be enough to cover the
// most common ones, it also does not cover the negative scenarios and all possible
// butchering of the license names
// Test private functions
func TestExtractLicenseFromClassifiers(t *testing.T) {
	classifiers := map[string][]string{
		"MIT License": {
			"Development Status :: 5 - Production/Stable",
			"Intended Audience :: Developers",
			"License :: OSI Approved :: MIT License",
			"Programming Language :: Python :: 3.6",
			"Programming Language :: Python :: 3.7",
		},
		"BSD License": {
			"Development Status :: 5 - Production/Stable",
			"Intended Audience :: Developers",
			"Intended Audience :: Information Technology",
			"License :: OSI Approved :: BSD License",
			"Operating System :: OS Independent",
			"Programming Language :: C",
			"Programming Language :: Cython",
			"Programming Language :: Python :: 2",
			"Programming Language :: Python :: 2.7",
			"Programming Language :: Python :: 3",
			"Programming Language :: Python :: 3.10",
			"Programming Language :: Python :: 3.11",
			"Programming Language :: Python :: 3.12",
			"Programming Language :: Python :: 3.6",
			"Programming Language :: Python :: 3.7",
			"Programming Language :: Python :: 3.8",
			"Programming Language :: Python :: 3.9",
			"Topic :: Software Development :: Libraries :: Python Modules",
			"Topic :: Text Processing :: Markup :: HTML",
			"Topic :: Text Processing :: Markup :: XML",
		},
	}
	for expectedLicense, classifiers := range classifiers {
		license := extractLicenseFromClassifiers(classifiers)
		if license != expectedLicense {
			t.Errorf("Expected %s, got %s", expectedLicense, license)
		}
	}
}

// TestGetPackageInfo tests whether the pypi JSON gets properly
// converted to the PackageMeta struct
// Uses the httptest package to mock the HTTP response and uses predownloaded
// JSON files to avoid hitting the PyPI API
func TestGetPackageInfo(t *testing.T) {
	data, err := os.ReadFile("testdata/pypi_response.json")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
	}))
	defer server.Close()

	pypi := NewPyPI("pypi", server.URL)
	meta, err := pypi.GetPackageInfo("beautifulsoup4")
	if err != nil {
		t.Fatal(err)
	}

	// hardcoded values based on the JSON file
	if meta.Author != "" {
		t.Errorf("Expected '', got %s", meta.Author)
	}
	if meta.Name != "beautifulsoup4" {
		t.Errorf("Expected beautifulsoup4, got %s", meta.Name)
	}
	if meta.Version != "4.12.2" {
		t.Errorf("Expected 4.12.2, got %s", meta.Version)
	}
	if meta.License != "MIT License" {
		t.Errorf("Expected MIT License, got %s", meta.License)
	}
	if !strings.HasPrefix(meta.Description, "Beautiful Soup is a library that") {
		t.Errorf("Expected 'Beautiful Soup is a library that', got %s", meta.Description)
	}
	if meta.Homepage != "https://www.crummy.com/software/BeautifulSoup/bs4/" {
		t.Errorf("Expected https://www.crummy.com/software/BeautifulSoup/bs4/, got %s", meta.Homepage)
	}
	if meta.Repository != "" {
		t.Errorf("Expected empty repository, got %s", meta.Repository)
	}
	if meta.Language != "python" {
		t.Errorf("Expected python, got %s", meta.Language)
	}

}
