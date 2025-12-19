package hashtag

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Parser struct {
}

func span(tag []byte) int {
	if idx := bytes.IndexFunc(tag, unicode.IsSpace); idx >= 0 {
		tag = tag[:idx]
	}
	end := len(tag)

	hasNonDigit := false
	for i := 0; i < end; {
		r, size := utf8.DecodeRune(tag[i:end])
		if r == utf8.RuneError && size == 1 {
			break
		}
		if endOfHashtag(r) && !isEmoji(r) {
			end = i
			break
		}
		if !unicode.IsDigit(r) {
			hasNonDigit = true
		}
		i += size
	}

	if !hasNonDigit {
		return -1
	}
	return end
}

var _ parser.InlineParser = (*Parser)(nil)

var _hash = byte('#')

func (*Parser) Trigger() []byte {
	return []byte{_hash}
}

func (p *Parser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()

	if len(line) == 0 || line[0] != _hash {
		return nil
	}
	line = line[1:]

	end := span(line)
	if end < 0 {
		return nil
	}
	seg = seg.WithStop(seg.Start + end + 1) // + '#'

	n := Node{
		Tag: block.Value(seg.WithStart(seg.Start + 1)), // omit the "#"
	}
	n.AppendChild(&n, ast.NewTextSegment(seg))
	block.Advance(seg.Len())
	return &n
}

func endOfHashtag(r rune) bool {
	return !unicode.IsLetter(r) &&
		!unicode.IsDigit(r) &&
		r != '_' && r != '-' && r != '/'
}

func isEmoji(r rune) bool {
	return unicode.Is(unicode.So, r)
}
