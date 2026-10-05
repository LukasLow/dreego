// Package dom re-exports gomponents so a dreego user never imports it directly.
//
// dreego builds on maragu.dev/gomponents (MIT). This package forwards its API
// under the dreego namespace: use dom.Text, dom.Div, dom.Class, dom.Map, …
// Node is a type alias, so dom.Node and gomponents Node are the same type.
package dom

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Type aliases (same types as gomponents).
type Node = g.Node
type Group = g.Group
type NodeType = g.NodeType
type Classes = c.Classes
type HTML5Props = c.HTML5Props

// Aliases for the components helpers.
var HTML5 = c.HTML5
var JoinAttrs = c.JoinAttrs

// Aliases for the core helpers.
var Attr = g.Attr
var El = g.El
var If = g.If
var Iff = g.Iff
var Raw = g.Raw
var Rawf = g.Rawf
var Text = g.Text
var Textf = g.Textf

// Aliases for the HTML elements and attributes.
var A = h.A
var Abbr = h.Abbr
var Accept = h.Accept
var Action = h.Action
var Address = h.Address
var Alt = h.Alt
var Area = h.Area
var Aria = h.Aria
var Article = h.Article
var As = h.As
var Aside = h.Aside
var Async = h.Async
var Audio = h.Audio
var AutoComplete = h.AutoComplete
var AutoFocus = h.AutoFocus
var AutoPlay = h.AutoPlay
var B = h.B
var Base = h.Base
var BlockQuote = h.BlockQuote
var Body = h.Body
var Br = h.Br
var Button = h.Button
var Canvas = h.Canvas
var Caption = h.Caption
var Charset = h.Charset
var Checked = h.Checked
var Cite = h.Cite
var CiteAttr = h.CiteAttr
var CiteEl = h.CiteEl
var Class = h.Class
var Code = h.Code
var Col = h.Col
var ColGroup = h.ColGroup
var ColSpan = h.ColSpan
var Cols = h.Cols
var Content = h.Content
var Controls = h.Controls
var CrossOrigin = h.CrossOrigin
var Data = h.Data
var DataAttr = h.DataAttr
var DataEl = h.DataEl
var DataList = h.DataList
var DateTime = h.DateTime
var Dd = h.Dd
var Defer = h.Defer
var Del = h.Del
var Details = h.Details
var Dfn = h.Dfn
var Dialog = h.Dialog
var Dir = h.Dir
var Disabled = h.Disabled
var Div = h.Div
var Dl = h.Dl
var Doctype = h.Doctype
var Download = h.Download
var Draggable = h.Draggable
var Dt = h.Dt
var Em = h.Em
var Embed = h.Embed
var EncType = h.EncType
var FieldSet = h.FieldSet
var FigCaption = h.FigCaption
var Figure = h.Figure
var Footer = h.Footer
var For = h.For
var Form = h.Form
var FormAction = h.FormAction
var FormAttr = h.FormAttr
var FormEl = h.FormEl
var FormEncType = h.FormEncType
var FormMethod = h.FormMethod
var FormNoValidate = h.FormNoValidate
var FormTarget = h.FormTarget
var H1 = h.H1
var H2 = h.H2
var H3 = h.H3
var H4 = h.H4
var H5 = h.H5
var H6 = h.H6
var HGroup = h.HGroup
var HTML = h.HTML
var Head = h.Head
var Header = h.Header
var Height = h.Height
var Hidden = h.Hidden
var Hr = h.Hr
var Href = h.Href
var I = h.I
var ID = h.ID
var IFrame = h.IFrame
var Img = h.Img
var Input = h.Input
var Ins = h.Ins
var Integrity = h.Integrity
var Kbd = h.Kbd
var Label = h.Label
var LabelAttr = h.LabelAttr
var LabelEl = h.LabelEl
var Lang = h.Lang
var Legend = h.Legend
var Li = h.Li
var Link = h.Link
var List = h.List
var Loading = h.Loading
var Loop = h.Loop
var Main = h.Main
var Mark = h.Mark
var Max = h.Max
var MaxLength = h.MaxLength
var Menu = h.Menu
var Meta = h.Meta
var Meter = h.Meter
var Method = h.Method
var Min = h.Min
var MinLength = h.MinLength
var Multiple = h.Multiple
var Muted = h.Muted
var Name = h.Name
var Nav = h.Nav
var NoScript = h.NoScript
var Object = h.Object
var Ol = h.Ol
var Open = h.Open
var OptGroup = h.OptGroup
var Option = h.Option
var Output = h.Output
var P = h.P
var Param = h.Param
var Pattern = h.Pattern
var Picture = h.Picture
var Placeholder = h.Placeholder
var PlaysInline = h.PlaysInline
var Popover = h.Popover
var PopoverTarget = h.PopoverTarget
var PopoverTargetAction = h.PopoverTargetAction
var Poster = h.Poster
var Pre = h.Pre
var Preload = h.Preload
var Progress = h.Progress
var Q = h.Q
var ReadOnly = h.ReadOnly
var ReferrerPolicy = h.ReferrerPolicy
var Rel = h.Rel
var Required = h.Required
var Role = h.Role
var RowSpan = h.RowSpan
var Rows = h.Rows
var S = h.S
var SVG = h.SVG
var Samp = h.Samp
var Scope = h.Scope
var Script = h.Script
var Search = h.Search
var Section = h.Section
var Select = h.Select
var Selected = h.Selected
var Sizes = h.Sizes
var SlotAttr = h.SlotAttr
var SlotEl = h.SlotEl
var Small = h.Small
var Source = h.Source
var Span = h.Span
var SpellCheck = h.SpellCheck
var Src = h.Src
var SrcSet = h.SrcSet
var Step = h.Step
var Strong = h.Strong
var Style = h.Style
var StyleAttr = h.StyleAttr
var StyleEl = h.StyleEl
var Sub = h.Sub
var Summary = h.Summary
var Sup = h.Sup
var TBody = h.TBody
var TFoot = h.TFoot
var THead = h.THead
var TabIndex = h.TabIndex
var Table = h.Table
var Target = h.Target
var Td = h.Td
var Template = h.Template
var Textarea = h.Textarea
var Th = h.Th
var Time = h.Time
var Title = h.Title
var TitleAttr = h.TitleAttr
var TitleEl = h.TitleEl
var Tr = h.Tr
var Type = h.Type
var U = h.U
var Ul = h.Ul
var Value = h.Value
var Var = h.Var
var Video = h.Video
var Wbr = h.Wbr
var Width = h.Width

// Map is written out because it is generic and cannot be aliased.
func Map[T any](ts []T, cb func(T) Node) Group { return g.Map(ts, cb) }
