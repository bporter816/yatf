package internal

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Linter struct {
	files []string
}

func NewLinter(files []string) *Linter {
	return &Linter{
		files: files,
	}
}

func (l Linter) Lint() {
	for _, v := range l.files {
		b, err := os.ReadFile(v)
		if err != nil {
			panic(err)
		}

		file, diag := hclwrite.ParseConfig(b, "", hcl.InitialPos)
		fmt.Printf("%+v\n", file.Body().Blocks())
		if diag.HasErrors() {
			panic(diag.Errs())
		}

		LintTrailingCommas(file)

		f, err := os.Create(v)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		file.WriteTo(f)
	}
}
