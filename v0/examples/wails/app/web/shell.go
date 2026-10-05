package web

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Shell is the page shell. It is a plain function: head content in, body
// content in, whole document out.
func Shell(c *dreego.Ctx, head d.View, body d.View) d.View {
	kopf := d.Group([]d.View{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		head,
	})
	return c.Document("en", kopf, body)
}

// shellCSS is the shared look of both pages. It is inlined (AddCritical) so the
// window paints themed immediately — no white flash on open.
const shellCSS = `
.shell { color-scheme: dark; --ink:#e7ebf5; --muted:#9aa3b8; --mint:#9fffd7; --violet:#b3a6ff;
  position: fixed; inset: 0; box-sizing: border-box; min-width: 20rem; overflow: auto;
  padding: 1.5rem clamp(1rem, 4vw, 2rem);
  background: radial-gradient(circle at 15% 8%, rgba(120,90,220,.22), transparent 40%),
              radial-gradient(circle at 85% 85%, rgba(35,141,115,.18), transparent 42%), #0b0d14;
  color: var(--ink);
  font-family: Inter, ui-sans-serif, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
.card { box-sizing:border-box; max-width:34rem; margin:0 auto; padding:clamp(1.35rem,4vw,2rem);
  border:1px solid rgba(255,255,255,.1); border-radius:1.5rem;
  background:linear-gradient(150deg, rgba(255,255,255,.07), rgba(255,255,255,.025));
  box-shadow:0 1.5rem 5rem rgba(0,0,0,.42), inset 0 1px rgba(255,255,255,.07); }
.eyebrow { margin:0 0 .5rem; color:var(--mint); font-size:.7rem; font-weight:800; letter-spacing:.16em; text-transform:uppercase; }
h1 { margin:0; font-size:clamp(1.6rem,5vw,2.4rem); font-weight:720; letter-spacing:-.04em; line-height:1.05; }
.intro { margin:.85rem 0 1.4rem; color:var(--muted); font-size:.9rem; line-height:1.55; }
code { font-family:ui-monospace,"SFMono-Regular",Consolas,monospace; font-size:.82em; background:rgba(255,255,255,.07); padding:.08rem .35rem; border-radius:.3rem; }
.greet { display:grid; gap:.5rem; }
.greet label { font-size:.78rem; font-weight:650; color:var(--muted); }
.row { display:flex; gap:.6rem; }
input { flex:1; min-width:0; padding:.65rem .8rem; border:1px solid rgba(255,255,255,.16); border-radius:.7rem; background:rgba(0,0,0,.25); color:var(--ink); font:inherit; }
button { padding:.65rem 1.1rem; border:0; border-radius:.7rem; background:var(--mint); color:#07130f; font:inherit; font-weight:750; cursor:pointer; transition:transform 150ms ease; }
button:hover { transform:translateY(-1px); }
input:focus-visible, button:focus-visible, a:focus-visible { outline:.2rem solid var(--violet); outline-offset:.18rem; }
.result { min-height:1.4rem; margin:.3rem 0 0; font-size:.95rem; font-weight:650; color:var(--mint); }
.topbar { display:flex; align-items:center; justify-content:space-between; max-width:40rem; margin:0 auto 1.25rem; }
.brand { display:inline-flex; align-items:center; gap:.55rem; font-size:.85rem; font-weight:750; letter-spacing:.03em; }
.brand-mark { display:grid; width:1.7rem; height:1.7rem; place-items:center; border:1px solid rgba(159,255,215,.42); border-radius:.5rem; background:rgba(159,255,215,.1); color:var(--mint); }
.native-badge { color:var(--muted); font-size:.68rem; font-weight:650; letter-spacing:.08em; text-transform:uppercase; }
footer { display:flex; max-width:34rem; align-items:center; justify-content:space-between; gap:1rem; margin:1rem auto 0; color:#747d97; font-size:.68rem; }
footer a { color:var(--mint); font-weight:700; text-decoration:none; }
.skip-link { position:absolute; left:-9999px; }
.skip-link:focus { position:fixed; top:.5rem; left:.5rem; z-index:100; padding:.5rem 1rem; background:var(--mint); color:#07130f; border-radius:.4rem; }
.shell:focus { outline:none; }
@media (max-width:27rem) { .native-badge, footer > span { display:none; } footer { justify-content:center; } }
@media (prefers-reduced-motion:reduce) { *, *::before, *::after { transition:none !important; } }
`
