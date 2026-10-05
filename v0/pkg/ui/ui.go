// Package ui is dreego's component library: reusable, themed building blocks.
//
// Every component is a Go function that takes the request context and its props,
// and returns a dom view. Styling is scoped through c.Box, so two
// components can use the same class names without clashing, and the CSS is
// emitted once per page.
//
// Rule: variants are props, never forks. A new look is a new value of an
// existing prop, not a copied component.
//
// Rule: never emit a style="…" attribute. dreego's CSP (style-src 'nonce-…')
// allows <style nonce> blocks but blocks inline style attributes. Put layout in
// a scoped CSS class instead — that is also what keeps the components themeable.
package ui

// Variant selects a visual style for buttons and alerts.
type Variant string

const (
	Primary Variant = "primary"
	Ghost   Variant = "ghost"
)

// Tone selects a semantic colour for badges and alerts.
type Tone string

const (
	Info    Tone = "info"
	Success Tone = "success"
	Warning Tone = "warning"
	Danger  Tone = "danger"
)
