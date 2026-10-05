# Abhängigkeiten und Lizenz

## dreego selbst

dreego steht unter der **MPL-2.0** (Mozilla Public License 2.0, siehe `LICENSE`).
Schwaches Copyleft: Änderungen **an dreego selbst** müssen offen bleiben, deine
Anwendung, die dreego *benutzt*, aber **nicht**.

## Keine Laufzeit-Abhängigkeiten

dreego kompiliert mit **null** externen Go-Modulen. `go.mod` enthält nur den
Modulnamen und die Go-Version. Routing, Session, CSRF, i18n, Markdown, Scoped
CSS, Static, Gzip, Fehlerseiten und die DOM-Schicht sind Go-Standardbibliothek
plus dreego-eigener Code.

## gomponents — Inspiration, keine Abhängigkeit

| | |
|---|---|
| Projekt | https://github.com/maragudk/gomponents |
| Lizenz | MIT — Copyright (c) Maragu ApS |
| Rolle | **Vorbild** für „HTML als Go-Funktionen" — **nicht** eingebunden |

dreego's `dom`-Paket (`pkg/dom`) ist eine **eigene** Umsetzung dieser Idee. Sie
wurde bewusst neu geschrieben, damit sie zu dreegos Modell passt — insbesondere
so, dass **jedes Attribut ein Knoten** ist. Dadurch kann eine spätere
htmx-Schicht (`hx-*`) im **selben** Modell leben, ohne Sonderfall.

Der Typ heißt **`View`** (statt `Node`) — klarer und selbsterklärend. Die
HTML-Elemente/-Attribute erzeugt `internal/gen_dom` aus der HTML-Spezifikation
(Element-/Attributnamen sind nicht urheberrechtlich geschützt).

**Warum kein Fork / kein kopierter Code?** Ein Fork hätte den `Node`-Typ und
gomponents' feste Attribut-Form mitgeschleppt — beides passt nicht zu
dreego + htmx. Eine eigene Implementierung ist ehrlicher und sauberer. Kein
gomponents-Code ist kopiert; die Nennung ist freier Dank.

Volltext: [THIRD_PARTY_NOTICES.md](../THIRD_PARTY_NOTICES.md).
