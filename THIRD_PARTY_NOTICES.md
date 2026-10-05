# Third-party notices

dreego stands on other people's work. This file lists every third-party
dependency, its license, and what it means for you.

---

## gomponents

- **Module:** `maragu.dev/gomponents`
- **Source:** https://github.com/maragudk/gomponents
- **Version pinned:** `v1.3.0`
- **License:** MIT — Copyright (c) Maragu AG

gomponents is the HTML component engine dreego builds on. dreego wraps it in
the `dreego/dom` package (a re-export) so that a normal dreego user never has to
import gomponents directly — but it is still there, under the hood.

**MIT means:** you may use, copy, modify, merge, publish, distribute,
sublicense and sell it, as long as the copyright notice and this license text
stay included. It is compatible with dreego's own MIT license. There is no
copyleft obligation: your application does not have to be open source.

The full MIT text of gompononents is reproduced below, as the license requires.

```
MIT License

Copyright (c) Maragu AG

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
