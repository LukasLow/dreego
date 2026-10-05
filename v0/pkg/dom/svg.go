package dom

// SVG helpers. These are hand-written (not part of the generated
// elements.go/attributes.go) because SVG is a separate namespace with its own
// element and attribute set. They follow the same model as everything else:
// every element is a View, every attribute is a View.
//
//	d.SVG(d.ViewBox("0 0 24 24"), d.Width("24"), d.Height("24"),
//	    d.Circle(d.Cx("12"), d.Cy("12"), d.R("9"),
//	        d.Fill("none"), d.Stroke("#0080ff"), d.StrokeWidth("2")),
//	)

// Circle rendert ein <circle>-Element.
func Circle(children ...View) View {
	return El("circle", children...)
}

// Ellipse rendert ein <ellipse>-Element.
func Ellipse(children ...View) View {
	return El("ellipse", children...)
}

// G rendert ein <g>-Element (SVG-Gruppe).
func G(children ...View) View {
	return El("g", children...)
}

// Line rendert ein <line>-Element.
func Line(children ...View) View {
	return El("line", children...)
}

// Path rendert ein <path>-Element.
func Path(children ...View) View {
	return El("path", children...)
}

// Polygon rendert ein <polygon>-Element.
func Polygon(children ...View) View {
	return El("polygon", children...)
}

// Polyline rendert ein <polyline>-Element.
func Polyline(children ...View) View {
	return El("polyline", children...)
}

// Rect rendert ein <rect>-Element.
func Rect(children ...View) View {
	return El("rect", children...)
}

// Defs rendert ein <defs>-Element.
func Defs(children ...View) View {
	return El("defs", children...)
}

// Use rendert ein <use>-Element.
func Use(children ...View) View {
	return El("use", children...)
}

// TextPath rendert ein <textPath>-Element.
func TextPath(children ...View) View {
	return El("textPath", children...)
}

// --- SVG-Attribute ---

// ViewBox setzt das viewBox-Attribut.
func ViewBox(v string) View { return Attr("viewBox", v) }

// X setzt das x-Attribut.
func X(v string) View { return Attr("x", v) }

// Y setzt das y-Attribut.
func Y(v string) View { return Attr("y", v) }

// X1 setzt das x1-Attribut.
func X1(v string) View { return Attr("x1", v) }

// Y1 setzt das y1-Attribut.
func Y1(v string) View { return Attr("y1", v) }

// X2 setzt das x2-Attribut.
func X2(v string) View { return Attr("x2", v) }

// Y2 setzt das y2-Attribut.
func Y2(v string) View { return Attr("y2", v) }

// Cx setzt das cx-Attribut.
func Cx(v string) View { return Attr("cx", v) }

// Cy setzt das cy-Attribut.
func Cy(v string) View { return Attr("cy", v) }

// Rx setzt das rx-Attribut.
func Rx(v string) View { return Attr("rx", v) }

// Ry setzt das ry-Attribut.
func Ry(v string) View { return Attr("ry", v) }

// R setzt das r-Attribut.
func R(v string) View { return Attr("r", v) }

// Points setzt das points-Attribut.
func Points(v string) View { return Attr("points", v) }

// D setzt das d-Attribut (Pfaddaten).
func D(v string) View { return Attr("d", v) }

// Fill setzt das fill-Attribut.
func Fill(v string) View { return Attr("fill", v) }

// FillRule setzt das fill-rule-Attribut.
func FillRule(v string) View { return Attr("fill-rule", v) }

// ClipRule setzt das clip-rule-Attribut.
func ClipRule(v string) View { return Attr("clip-rule", v) }

// Stroke setzt das stroke-Attribut.
func Stroke(v string) View { return Attr("stroke", v) }

// StrokeWidth setzt das stroke-width-Attribut.
func StrokeWidth(v string) View { return Attr("stroke-width", v) }

// StrokeLinecap setzt das stroke-linecap-Attribut.
func StrokeLinecap(v string) View { return Attr("stroke-linecap", v) }

// StrokeLinejoin setzt das stroke-linejoin-Attribut.
func StrokeLinejoin(v string) View { return Attr("stroke-linejoin", v) }

// StrokeDasharray setzt das stroke-dasharray-Attribut.
func StrokeDasharray(v string) View { return Attr("stroke-dasharray", v) }

// StrokeDashoffset setzt das stroke-dashoffset-Attribut.
func StrokeDashoffset(v string) View { return Attr("stroke-dashoffset", v) }

// StrokeOpacity setzt das stroke-opacity-Attribut.
func StrokeOpacity(v string) View { return Attr("stroke-opacity", v) }

// FillOpacity setzt das fill-opacity-Attribut.
func FillOpacity(v string) View { return Attr("fill-opacity", v) }

// Opacity setzt das opacity-Attribut.
func Opacity(v string) View { return Attr("opacity", v) }

// Transform setzt das transform-Attribut.
func Transform(v string) View { return Attr("transform", v) }

// Xmlns setzt das xmlns-Attribut.
func Xmlns(v string) View { return Attr("xmlns", v) }

// AriaLabel setzt das aria-label-Attribut.
func AriaLabel(v string) View { return Attr("aria-label", v) }

// AriaHidden setzt das aria-hidden-Attribut.
func AriaHidden(v string) View { return Attr("aria-hidden", v) }
