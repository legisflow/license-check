// Test the requirements.txt parser
package python

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func stringPtr(s string) *string {
	return &s
}

func TestRequirementsTxtBasicParser(t *testing.T) {
	r := newRequirementsTxtParser()

	data, err := ioutil.ReadFile("testdata/requirements_basic.txt")
	if err != nil {
		t.Fatal(err)
	}

	req, err := r.ParseString("", string(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(req.Line) != 15 {
		// comments are excluded from parser output
		t.Fatalf("Expected 15 lines, got %d", len(req.Line))
	}

	type version struct {
		name     string
		operator string
		version  string
	}

	expectedLines := []version{
		{"contourpy", "==", "1.0.6"},
		{"cycler", ">", "0.11.0"},
		{"fonttools", ">=", "4.38.0"},
		{"kiwisolver", "<=", "1.4.4"},
		{"matplotlib", "<", "3.6.2"},
		{"numpy", "==", "1.23.5"},
		{"packaging", "==", "21.3"},
		{"pandas", "==", "1.5.1"},
		{"pillow", "==", "9.3.0"},
		{"pyparsing", "==", "3.0.9"},
		{"python-dateutil", "==", "2.8.2"},
		{"pytz", "==", "2022.6"},
		{"scipy", "==", "1.9.3"},
		{"seaborn", "==", "0.10.1"},
		{"six", "==", "1.16.0"},
	}

	for i, l := range req.Line {
		if l.Package == nil {
			t.Fatalf("Expected package, got nil")
		}

		if l.Package.Name != expectedLines[i].name {
			t.Fatalf("Expected name %s, got %s", expectedLines[i].name, l.Package.Name)
		}

		if l.Package.Versions == nil {
			t.Fatal("Expected versions, got nil")
		}

		if len(l.Package.Versions) != 1 {
			t.Fatalf("Expected 1 version, got %d", len(l.Package.Versions))
		}

		if l.Package.Versions[0].Operator != expectedLines[i].operator {
			t.Fatalf("Expected operator %s, got %s", expectedLines[i].operator, l.Package.Versions[0].Operator)
		}

		if l.Package.Versions[0].Value != expectedLines[i].version {
			t.Fatalf("Expected version %s, got %s", expectedLines[i].version, l.Package.Versions[0].Value)
		}
	}
}

func TestRequirementsTxtFullLine(t *testing.T) {
	line := `contourpy [bold, mega] >1.0.6, <=1.1.2; python_version >= "3.6", platform == "linux" # comment`

	r := newRequirementsTxtParser()
	req, err := r.ParseString("", line)
	if err != nil {
		t.Fatal(err)
	}

	if len(req.Line) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(req.Line))
	}

	if req.Line[0].Package == nil {
		t.Fatalf("Expected package, got nil")
	}

	if req.Line[0].Package.Name != "contourpy" {
		t.Fatalf("Expected package contourpy, got %s", req.Line[0].Package.Name)
	}

	if len(req.Line[0].Package.Extras) != 2 {
		t.Fatalf("Expected 2 extras, got %d", len(req.Line[0].Package.Extras))
	}

	if req.Line[0].Package.Extras[0] != "bold" {
		t.Fatalf("Expected extra bold, got %s", req.Line[0].Package.Extras[0])
	}

	if req.Line[0].Package.Extras[1] != "mega" {
		t.Fatalf("Expected extra mega, got %s", req.Line[0].Package.Extras[1])
	}

	if len(req.Line[0].Package.Versions) != 2 {
		t.Fatalf("Expected 2 versions, got %d", len(req.Line[0].Package.Versions))
	}

	if req.Line[0].Package.Versions[0].Operator != ">" {
		t.Fatalf("Expected operator >, got %s", req.Line[0].Package.Versions[0].Operator)
	}

	if req.Line[0].Package.Versions[0].Value != "1.0.6" {
		t.Fatalf("Expected version 1.0.6, got %s", req.Line[0].Package.Versions[0].Value)
	}

	if req.Line[0].Package.Versions[1].Operator != "<=" {
		t.Fatalf("Expected operator <=, got %s", req.Line[0].Package.Versions[1].Operator)
	}

	if req.Line[0].Package.Versions[1].Value != "1.1.2" {
		t.Fatalf("Expected version 1.1.2, got %s", req.Line[0].Package.Versions[1].Value)
	}

	if len(req.Line[0].Package.Environs) != 2 {
		t.Fatalf("Expected 2 environment, got %d", len(req.Line[0].Package.Environs))
	}

	if req.Line[0].Package.Environs[0].Name != "python_version" {
		t.Fatalf("Expected environment python_version, got %s", req.Line[0].Package.Environs[0].Name)
	}

	if req.Line[0].Package.Environs[0].Operator != ">=" {
		t.Fatalf("Expected operator >=, got %s", req.Line[0].Package.Environs[0].Operator)
	}

	if req.Line[0].Package.Environs[0].Value != "3.6" {
		t.Fatalf("Expected version 3.6, got %s", req.Line[0].Package.Environs[0].Value)
	}

	if req.Line[0].Package.Environs[1].Name != "platform" {
		t.Fatalf("Expected environment platform, got %s", req.Line[0].Package.Environs[1].Name)
	}

	if req.Line[0].Package.Environs[1].Operator != "==" {
		t.Fatalf("Expected operator ==, got %s", req.Line[0].Package.Environs[1].Operator)
	}

	if req.Line[0].Package.Environs[1].Value != "linux" {
		t.Fatalf("Expected version linux, got %s", req.Line[0].Package.Environs[1].Value)
	}
}

func TestRequirementsTxtAlterantiveDownloads(t *testing.T) {
	line := `urllib3 [security] @ https://github.com/urllib3/urllib3/archive/refs/tags/1.26.8.zip; platform_system != "AIX" # comment`

	r := newRequirementsTxtParser()
	req, err := r.ParseString("", line)
	if err != nil {
		t.Fatal(err)
	}

	if len(req.Line) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(req.Line))
	}

	if req.Line[0].Package == nil {
		t.Fatalf("Expected package, got nil")
	}

	if req.Line[0].Package.Name != "urllib3" {
		t.Fatalf("Expected package urllib3, got %s", req.Line[0].Package.Name)
	}

	if len(req.Line[0].Package.Extras) != 1 {
		t.Fatalf("Expected 1 extra, got %d", len(req.Line[0].Package.Extras))
	}

	if req.Line[0].Package.Extras[0] != "security" {
		t.Fatalf("Expected extra security, got %s", req.Line[0].Package.Extras[0])
	}

	if req.Line[0].Package.Download == nil {
		t.Fatalf("Expected download, got nil")
	}

	if *req.Line[0].Package.Download != "https://github.com/urllib3/urllib3/archive/refs/tags/1.26.8.zip" {
		t.Fatalf("Expected download URL not matching")
	}

	if len(req.Line[0].Package.Environs) != 1 {
		t.Fatalf("Expected 1 environment, got %d", len(req.Line[0].Package.Environs))
	}

	if req.Line[0].Package.Environs[0].Name != "platform_system" {
		t.Fatalf("Expected environment platform_system, got %s", req.Line[0].Package.Environs[0].Name)
	}

	if req.Line[0].Package.Environs[0].Operator != "!=" {
		t.Fatalf("Expected operator !=, got %s", req.Line[0].Package.Environs[0].Operator)
	}

	if req.Line[0].Package.Environs[0].Value != "AIX" {
		t.Fatalf("Expected version AIX, got %s", req.Line[0].Package.Environs[0].Value)
	}
}

func TestRequirementsTxtURL(t *testing.T) {
	line := `https://github.com/urllib3/urllib3/archive/refs/tags/1.26.8.zip`

	r := newRequirementsTxtParser()
	req, err := r.ParseString("", line)
	if err != nil {
		t.Fatal(err)
	}

	if len(req.Line) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(req.Line))
	}

	if req.Line[0].Download == nil {
		t.Fatalf("Expected download, got nil")
	}

	if *req.Line[0].Download != "https://github.com/urllib3/urllib3/archive/refs/tags/1.26.8.zip" {
		t.Fatalf("Expected download URL not matching")
	}
}

func TestRequirementsTxtLineVariants(t *testing.T) {
	lines := []string{
		`numpy`,
		`numpy[all]`,
		`numpy[all, noop]; python_version >= "3.6"`,
		`numpy >= 1.23.5, <= 1.23.6`,
		`--no-binary`,
		`-r requirements.txt`,
		`# comment`,
	}

	data := strings.Join(lines, "\n")

	r := newRequirementsTxtParser()
	req, err := r.ParseString("", data)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Requirement{
		Line: []Line{
			{
				Package: &Package{
					Name: "numpy",
				},
			},
			{
				Package: &Package{
					Name:   "numpy",
					Extras: []string{"all"},
				},
			},
			{
				Package: &Package{
					Name:   "numpy",
					Extras: []string{"all", "noop"},
					Environs: []*Environ{
						{
							Name:     "python_version",
							Operator: ">=",
							Value:    "3.6",
						},
					},
				},
			},
			{
				Package: &Package{
					Name: "numpy",
					Versions: []*Version{
						{
							Operator: ">=",
							Value:    "1.23.5",
						},
						{
							Operator: "<=",
							Value:    "1.23.6",
						},
					},
				},
			},
			{
				Command: &Command{
					LongName: stringPtr("no-binary"),
				},
			},
			{
				Command: &Command{
					ShortName: stringPtr("r"),
					Option:    stringPtr("requirements.txt"),
				},
			},
		},
	}

	if diff := cmp.Diff(expected, req); diff != "" {
		t.Fatalf("RequirementsTxt parser mismatch (-want +got):\n%s", diff)
	}
}
