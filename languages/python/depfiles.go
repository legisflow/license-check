// Implement the depfiles interface for the python language.
// It contains multiple implementations of the DepFile interface, each of which
// uses different file format.
package python

import (
	"bufio"
	"context"
	"io/ioutil"
	"strings"

	"github.com/radiculaCZ/license-check/interfaces"
)

type RequirementsTxt struct {
}

// NewRequirementsTxt creates a new instance of the RequirementsTxt struct
func NewRequirementsTxt() interfaces.DepFile {
	return &RequirementsTxt{}
}

func (r *RequirementsTxt) GetDepFileType() string {
	return "python/requirements.txt"
}

func (r *RequirementsTxt) GetRepository() interfaces.PackageRepository {
	return NewPyPI("PyPI", PyPIURL)
}

func (r *RequirementsTxt) GetDependencies(ctx context.Context, fileName string) (<-chan string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	stringReader := strings.NewReader(string(data))
	bufScanner := bufio.NewScanner(stringReader)
	depChan := make(chan string)

	go func() {
		defer close(depChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !bufScanner.Scan() {
					return
				}
				depChan <- bufScanner.Text()
			}
		}
	}()

	return depChan, nil
}
