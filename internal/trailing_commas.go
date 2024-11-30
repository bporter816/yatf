package internal

import (
	"bytes"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func lintTrailingCommas(b *hclwrite.Body) {
	if b == nil {
		return
	}

	for k, v := range b.Attributes() {
		tok := v.Expr().BuildTokens(nil)
		stack := make([]hclsyntax.TokenType, 0)

		// only loop until second to first character because we don't have to bounds check i-1
		// and nothing will need to be fixed there anyway
		for i := len(tok) - 1; i > 0; i-- {
			if tok[i].Type == hclsyntax.TokenCBrack || tok[i].Type == hclsyntax.TokenCBrace {
				stack = append(stack, tok[i].Type)
			}

			if tok[i].Type == hclsyntax.TokenNewline {
				hasComma := tok[i-1].Type == hclsyntax.TokenComma

				// list with no trailing comma
				if len(stack) > 0 && stack[len(stack)-1] == hclsyntax.TokenCBrack && !hasComma {
					// special case for first list item
					if tok[i-1].Type == hclsyntax.TokenOBrack {
						continue
					}

					tok = slices.Insert(tok, i, &hclwrite.Token{
						Type:         hclsyntax.TokenComma,
						Bytes:        []byte{','},
						SpacesBefore: 0,
					})
				}

				// map with trailing comma
				if len(stack) > 0 && stack[len(stack)-1] == hclsyntax.TokenCBrace && hasComma {
					tok = slices.Delete(tok, i-1, i)
				}
			}

			if tok[i].Type == hclsyntax.TokenOBrace {
				if stack[len(stack)-1] != hclsyntax.TokenCBrace {
					panic("bad format")
				}
				stack = slices.Delete(stack, len(stack)-1, len(stack))
			}
			if tok[i].Type == hclsyntax.TokenOBrack {
				if stack[len(stack)-1] != hclsyntax.TokenCBrack {
					panic("bad format")
				}
				stack = slices.Delete(stack, len(stack)-1, len(stack))
			}

		}
		b.SetAttributeRaw(k, tok)
	}

	for _, v := range b.Blocks() {
		if v != nil {
			lintTrailingCommas(v.Body())
		}
	}
}

func LintTrailingCommas(f *hclwrite.File) {
	lintTrailingCommas(f.Body())
}

func LintTrailingCommasString(s string) string {
	data, err := hclwrite.ParseConfig([]byte(s), "", hcl.InitialPos)
	if err != nil {
		panic(err)
	}
	LintTrailingCommas(data)
	var buf bytes.Buffer
	data.WriteTo(&buf)
	return string(buf.Bytes())
}
