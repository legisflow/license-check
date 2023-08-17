package python

import (
	"context"
	"io/ioutil"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/radiculaCZ/license-check/interfaces"
)

type Requirement struct {
	Line []*Line `@@*`
}

type Line struct {
	Package  *Package `  @@`
	Command  *Command `| @@`
	Download *string  `| @Download`
}

type Command struct {
	LongName  *string `   "--"  @Ident`
	ShortName *string `| "-" @Ident`
	Option    *string `@Ident`
}

type Package struct {
	Name     string     `@Ident`
	Extras   []string   `("[" @Ident ("," @Ident)* "]")?`
	Download *string    `("@" @Download)?`
	Versions []*Version `( @@ ( "," @@ )* )?`
	Environs []*Environ `(";" @@ ( "," @@ )* )?`
}

type Version struct {
	Operator string `@Operator`
	Value    string `@VersionValue`
}

type Environ struct {
	Name     string `@Ident`
	Operator string `@Operator`
	Value    string `@EnvironVersion`
}

type RequirementsTxt struct {
	lexer  *lexer.StatefulDefinition
	parser *participle.Parser[Requirement]
}

// NewRequirementsTxt creates a new instance of the RequirementsTxt struct
func NewRequirementsTxt() interfaces.DepFile {
	requirementLexer := lexer.MustSimple([]lexer.SimpleRule{
		{`Comment`, `#.*`},
		{`Command`, `--|-`},
		{`Download`, `((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[.\!\/\\w]*))?)`},
		{`Ident`, `[a-zA-Z_][a-zA-Z_0-9\-.]*`},
		{`Operator`, `==|!=|~=|>=|>|<=|<`},
		{`VersionValue`, `[0-9\.\*]+`},
		{`EnvironVersion`, `'(\\'|[^'])*'|"(\\"|[^"])*"`},
		{"Punct", `\[|]|[-!()+*=,;@]`},
		{`Whitespace`, `\s+`},
	})

	requirementParser := participle.MustBuild[Requirement](
		participle.Lexer(requirementLexer),
		participle.Unquote("EnvironVersion"),
		participle.UseLookahead(2),
		participle.Elide("Whitespace", "Comment"),
	)

	return &RequirementsTxt{
		lexer:  requirementLexer,
		parser: requirementParser,
	}
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

	req, err := r.parser.ParseString("", string(data))
	if err != nil {
		return nil, err
	}

	depChan := make(chan string)

	go func() {
		defer close(depChan)
		for _, line := range req.Line {
			select {
			case <-ctx.Done():
				return
			default:
				if line.Package != nil {
					depChan <- line.Package.Name
				}
			}
		}
	}()

	return depChan, nil
}
