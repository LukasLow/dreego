# Third-party notices

dreego stands on other people's work. This file lists every third-party
dependency, its license, and what it means for you.

---

## gomponents

- **Canonical module path:** `maragu.dev/gomponents` (Go requires this in imports — the declared path)
- **Fetched from (real source):** `github.com/maragudk/gomponents`, pinned to tag **`v1.3.0`** via a `replace` directive
- **License:** MIT — Copyright (c) **Maragu ApS**

dreego builds on gomponents and pins it explicitly:

```
require maragu.dev/gomponents v1.3.0
replace maragu.dev/gomponents => github.com/maragudk/gomponents v1.3.0
```

The canonical path must stay in imports (Go rejects a module whose declared
path differs), so the `replace` points the build at the GitHub repository and
the release tag. dreego wraps gomponents in the `dom` package so a normal user
never imports it directly.

**MIT means:** you may use, copy, modify, merge, publish, distribute,
sublicense and sell it, as long as the copyright notice and this license text
stay included. MIT (permissive) is compatible with dreego's own **MPL-2.0**
(weak copyleft): the MIT part stays MIT, and dreego's MPL applies only to
dreego's own files. There is no obligation to open your application.

The full MIT text of gomponents is reproduced below, as the license requires.

```
MIT License

Copyright (c) Maragu ApS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## No other runtime dependencies

dreego's core has **no** runtime dependencies besides gomponents (which itself
has none). Everything else — routing, session, CSRF, i18n, markdown, scoped CSS,
static, compression, error pages — is written in the Go standard library.
