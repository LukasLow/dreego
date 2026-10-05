// Federkiel — kleines Client-Skript als statische Datei.
// Beweist, dass externes JS (go:embed) unter der CSP-nonce-Regel nicht nötig
// ist: diese Datei wird per <script src> geladen und braucht keine nonce,
// weil script-src 'nonce-…' nur Inline-Skripte betrifft — externe Dateien
// fallen unter default-src 'self'.
document.addEventListener("DOMContentLoaded", function () {
    var zahl = document.querySelector("[data-woerter]");
    if (!zahl) { return; }
    var aktion = document.querySelector("[data-woerter-action]");
    var aktuell = 0;
    aktion.addEventListener("click", function () {
        aktuell += 137;
        zahl.textContent = String(aktuell);
    });
});
