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
	fileName string
	file     string
}

// NewRequirementsTxt creates a new instance of the RequirementsTxt struct
// It reads the file and stores the content in the struct
// It returns nil and error in case the file cannot be read
func NewRequirementsTxt(fileName string) (interfaces.DepFile, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &RequirementsTxt{
		fileName: fileName,
		file:     string(data),
	}, nil
}

func (r *RequirementsTxt) GetDepFileName() string {
	return r.fileName
}

func (r *RequirementsTxt) GetDepFileType() string {
	return "requirements.txt"
}

func (r *RequirementsTxt) GetDependencies(ctx context.Context) <-chan string {
	stringReader := strings.NewReader(r.file)
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

	return depChan
}
