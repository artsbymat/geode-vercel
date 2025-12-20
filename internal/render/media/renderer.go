package media

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Renderer struct{}

var _ renderer.NodeRenderer = (*Renderer)(nil)

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *Renderer) renderImage(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*ast.Image)
	if !ok {
		return ast.WalkContinue, nil
	}

	if tweet, user, ok := tweetURL(n.Destination); ok {
		_, _ = w.WriteString(`<blockquote class="twitter-tweet"><a href="`)
		_, _ = w.Write(util.URLEscape([]byte(tweet), true /* resolve references */))
		_, _ = w.WriteString(`\">@`)
		_, _ = w.Write(util.EscapeHTML([]byte(user)))
		_, _ = w.WriteString(`</a></blockquote>`)
		return ast.WalkSkipChildren, nil
	}

	if id, isShort, ok := youtubeID(n.Destination); ok {
		altRaw := nodeText(src, n)
		_, width, hasWidth := parseAltAndWidth(altRaw)

		embed := "https://www.youtube.com/embed/" + id
		wAttr, hAttr := 560, 315
		if isShort {
			wAttr, hAttr = 315, 560
		}
		if hasWidth {
			wAttr = width
			if isShort {
				hAttr = (width * 16) / 9
				if hAttr <= 0 {
					hAttr = 560
				}
			} else {
				hAttr = (width * 9) / 16
				if hAttr <= 0 {
					hAttr = 315
				}
			}
		}

		_, _ = w.WriteString(`<iframe src="`)
		_, _ = w.Write(util.URLEscape([]byte(embed), true /* resolve references */))
		_, _ = w.WriteString(`"`)
		_, _ = w.WriteString(` width="`)
		_, _ = w.WriteString(strconv.Itoa(wAttr))
		_, _ = w.WriteString(`" height="`)
		_, _ = w.WriteString(strconv.Itoa(hAttr))
		_, _ = w.WriteString(`"`)
		_, _ = w.WriteString(` frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen loading="lazy"></iframe>`)
		return ast.WalkSkipChildren, nil
	}

	altRaw := nodeText(src, n)
	alt, width, hasWidth := parseAltAndWidth(altRaw)

	_, _ = w.WriteString(`<img src="`)
	_, _ = w.Write(util.URLEscape(n.Destination, true /* resolve references */))
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(util.EscapeHTML(alt))
	_, _ = w.WriteString(`"`)

	if len(n.Title) > 0 {
		_, _ = w.WriteString(` title="`)
		_, _ = w.Write(util.EscapeHTML(n.Title))
		_, _ = w.WriteString(`"`)
	}
	if hasWidth {
		_, _ = w.WriteString(` width="`)
		_, _ = w.WriteString(strconv.Itoa(width))
		_, _ = w.WriteString(`"`)
	}

	_, _ = w.WriteString(`>`)
	return ast.WalkSkipChildren, nil
}

func youtubeID(dest []byte) (id string, isShort bool, ok bool) {
	u, err := url.Parse(strings.TrimSpace(string(dest)))
	if err != nil || u == nil {
		return "", false, false
	}

	host := strings.ToLower(u.Host)
	path := strings.Trim(u.Path, "/")

	if host == "youtu.be" {
		if path == "" {
			return "", false, false
		}
		return path, false, true
	}

	if host == "www.youtube.com" || host == "youtube.com" || host == "m.youtube.com" {
		if path == "watch" {
			v := u.Query().Get("v")
			if v == "" {
				return "", false, false
			}
			return v, false, true
		}
		if v, ok := strings.CutPrefix(path, "embed/"); ok {
			if v == "" {
				return "", false, false
			}
			return v, false, true
		}

		if v, ok := strings.CutPrefix(path, "shorts/"); ok {
			if v == "" {
				return "", false, false
			}
			if i := strings.IndexByte(v, '/'); i >= 0 {
				v = v[:i]
			}
			return v, true, true
		}
	}

	return "", false, false
}

func tweetURL(dest []byte) (tweet string, user string, ok bool) {
	u, err := url.Parse(strings.TrimSpace(string(dest)))
	if err != nil || u == nil {
		return "", "", false
	}

	host := strings.ToLower(u.Host)
	if host != "x.com" && host != "www.x.com" && host != "twitter.com" && host != "www.twitter.com" {
		return "", "", false
	}

	path := strings.Trim(u.Path, "/")
	// Expected: <user>/status/<id>
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return "", "", false
	}
	if parts[1] != "status" {
		return "", "", false
	}
	if parts[0] == "" || parts[2] == "" {
		return "", "", false
	}

	user = parts[0]
	id := parts[2]
	return "https://twitter.com/" + user + "/status/" + id, user, true
}

func parseAltAndWidth(label []byte) (alt []byte, width int, ok bool) {
	s := strings.TrimSpace(string(label))
	if s == "" {
		return nil, 0, false
	}

	if i := strings.LastIndexByte(s, '|'); i >= 0 {
		left := strings.TrimSpace(s[:i])
		right := strings.TrimSpace(s[i+1:])
		if w, err := strconv.Atoi(right); err == nil && w > 0 {
			if left == "" {
				return nil, w, true
			}
			return []byte(left), w, true
		}
	}

	if w, err := strconv.Atoi(s); err == nil && w > 0 {
		return nil, w, true
	}

	return []byte(s), 0, false
}

func nodeText(src []byte, n ast.Node) []byte {
	var buf bytes.Buffer
	writeNodeText(src, &buf, n)
	return buf.Bytes()
}

func writeNodeText(src []byte, dst *bytes.Buffer, n ast.Node) {
	switch n := n.(type) {
	case *ast.Text:
		_, _ = dst.Write(n.Segment.Value(src))
	case *ast.String:
		_, _ = dst.Write(n.Value)
	default:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			writeNodeText(src, dst, c)
		}
	}
}
