package printer

import (
	"go/ast"
	"go/token"
)

// Is the struct short enough to fit on one line
// For now limit to anonymous structs
// Part of vertical_compression patchset
func (p *printer) compressibleStruct(list []*ast.Field) bool {
	const maxSize = 30
	cum := p.indent*p.Config.Tabwidth + p.pos.Column + len("{ ")
	for i, f := range list {
		if f.Tag != nil || f.Comment != nil {
			return false
		}
		if f.Names != nil {
			return false
		} // only anonymous
		namesSize := identListSize(f.Names, maxSize)
		if namesSize > 0 {
			namesSize += 1
		}
		cum += p.nodeSize(f.Type, maxSize) + namesSize
		if i > 0 {
			cum += len("; ")
		}
	}
	cum += len(" }")
	return cum < 70
}

// Compress structs
// Part of vertical_compression patchset
func (p *printer) compressedStruct(lbrace token.Pos, list []*ast.Field) {
	p.print(lbrace, blank, token.LBRACE, blank)
	for i, f := range list {
		if i > 0 {
			p.print(token.SEMICOLON, blank)
		}
		for j, n := range f.Names {
			if j > 0 {
				p.print(token.COMMA, blank)
			}
			p.expr(n)
		}
		if len(f.Names) > 0 {
			p.print(blank)
		}
		p.expr(f.Type)
	}
	p.print(blank, token.RBRACE)
}

// from vertical_compression patchset
func (p *printer) compressedBlock(b *ast.BlockStmt) bool {
	var did_compress bool
	if len(b.List) == 1 {
		did_compress = p.openCompressedBlock(b)
		if did_compress {
			p.print(blank, b.Rbrace, token.RBRACE)
		}
	}
	return did_compress
}

// from vertical_compression patchset
// compressed block with no closing Rbrace, for use with if/else
func (p *printer) openCompressedBlock(b *ast.BlockStmt) bool {
	// eg. the statement ('if ... {', 'for ... {', etc)
	st_pos, st_end := b.Pos(), b.End()
	_ = st_end
	// the expression block
	e_pos, e_end := b.List[0].Pos(), b.List[0].End()
	// check that expression isn't multi-line
	if (p.lineFor(e_end) - p.lineFor(e_pos)) == 0 {
		// if short enough, keep on same line
		// starting column of block (opening {) - sans tab chars
		col := p.posFor(st_pos).Column - p.indent
		// indent in displayed chars
		indent := p.indent * p.Config.Tabwidth
		end_bracket_w_padding := 3 // hardcoded below
		block_len := int(e_end-e_pos) + end_bracket_w_padding
		if (block_len + indent + col) < 78 {
			p.print(b.Pos(), token.LBRACE, blank)
			p.stmt(b.List[0], true)
			return true
		}
	}
	return false
}
