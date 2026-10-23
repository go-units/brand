// Command mark writes the go-units mark, svg/color/go-units.svg:
//
//	cd mark && go run . > ../svg/color/go-units.svg
//
// go-units is units and physical constants, and the mark is the instrument that
// stands for measuring: a ruler. Its graduations are the ones a ruler has, a
// long one at the middle and short ones between.
//
// Nothing is placed by hand. The constants below give the size of the ruler and
// the number of divisions; the program derives the pitch, the length of each
// graduation and every coordinate from them, and centres the drawing on the
// square. Change the divisions and the mark redraws itself.
//
// The rest is the family grammar every go-* mark shares: a 256-unit rounded
// square (rx 56) filled with a diagonal two-stop gradient, a white line glyph
// at stroke-width 13 with round caps and joins. logos/brandgen renders the
// format matrix from the file this writes; the kit workflow runs this again
// and fails if a byte of svg/color differs.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	light = "#A3B84A" // the gradient's light stop
	dark  = "#5F6B12" // its dark stop

	stroke = 13.0 // the family's line width

	width     = 184.0 // the ruler, between stroke centres
	height    = 80.0
	corner    = 14.0 // the ruler's corner radius
	divisions = 6    // graduations are drawn at the divisions between the two ends

	longTick  = 42.0 // at the middle
	shortTick = 26.0 // everywhere else
)

// num prints a coordinate with at most one decimal and no trailing zero.
func num(v float64) string {
	s := fmt.Sprintf("%.1f", math.Round(v*10)/10)
	s = strings.TrimSuffix(s, ".0")
	if s == "-0" {
		s = "0"
	}
	return s
}

// tickLength is the length of the graduation at division i: the middle is the
// longest, the rest are short.
func tickLength(i int) float64 {
	if i == divisions/2 {
		return longTick
	}
	return shortTick
}

func main() {
	left, top := 128-width/2, 128-height/2
	pitch := width / divisions

	paths := []string{
		fmt.Sprintf(`<rect x="%s" y="%s" width="%s" height="%s" rx="%s"/>`,
			num(left), num(top), num(width), num(height), num(corner)),
	}
	// The graduations hang from the top edge; the two ends are the ruler's own.
	for i := 1; i < divisions; i++ {
		x := left + float64(i)*pitch
		paths = append(paths, fmt.Sprintf(`<path d="M%s %sV%s"/>`, num(x), num(top), num(top+tickLength(i))))
	}

	fmt.Fprintf(os.Stdout, `<svg viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="go-units">
  <defs><linearGradient id="ungrad" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="%s"/><stop offset="1" stop-color="%s"/></linearGradient></defs>
  <rect width="256" height="256" rx="56" fill="url(#ungrad)"/>
  <g fill="none" stroke="#FFFFFF" stroke-width="%s" stroke-linecap="round" stroke-linejoin="round">
    %s
  </g>
</svg>
`, light, dark, num(stroke), strings.Join(paths, "\n    "))
}
