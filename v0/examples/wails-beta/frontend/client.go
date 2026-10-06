package frontend

// clientJS is the page's client module: a plain-JavaScript port of the Wails
// vanilla template's main.ts. It calls the same Go service (GreetService.Greet)
// and listens for the same "time" event, but imports the Wails runtime and the
// generated binding directly as ES modules (both are served in-process) instead
// of being bundled by Vite.
//
// It is emitted as <script type="module" nonce="…"> by page.go, so it runs as a
// module and may use static imports.
const clientJS = `
import { Events, WML } from "/wails/runtime.js";
import { GreetService } from "/static/bindings/wails-dreego/index.js";

// Wire up data-wml-openURL links (logos + footer "Docs" link) once the DOM is ready.
WML.Enable();

const greetButton = document.getElementById("greet");
const nameElement = document.getElementById("name");
const resultElement = document.getElementById("result");
const timeElement = document.getElementById("time");
const titleNameElement = document.querySelector(".title-name");
const toastElement = document.getElementById("toast");
let toastTimer;

// Crossfade the framework word in the heading ("Wails + dreego") to the name the
// user entered ("Wails + <name>"): the old word fades out while the new one fades
// in over the same spot.
function swapTitleName(name) {
	const current = titleNameElement.querySelector(".title-name-text:not(.is-outgoing)");
	if (!current || current.textContent === name) {
		return;
	}
	const incoming = document.createElement("span");
	incoming.className = "title-name-text is-entering";
	incoming.textContent = name;
	current.classList.add("is-outgoing");
	titleNameElement.appendChild(incoming);
	// Force a reflow so the transitions run from the starting state.
	void incoming.offsetWidth;
	incoming.classList.remove("is-entering");
	current.classList.add("is-leaving");
	current.addEventListener("transitionend", () => current.remove(), { once: true });
}

// Pop the toast with the message Go returned, then auto-dismiss it.
function showToast(message) {
	resultElement.innerText = message;
	toastElement.classList.add("is-visible");
	clearTimeout(toastTimer);
	toastTimer = setTimeout(() => toastElement.classList.remove("is-visible"), 4000);
}

greetButton.addEventListener("click", async () => {
	let name = nameElement.value;
	if (!name) {
		name = "anonymous";
	}
	swapTitleName(name);
	try {
		showToast(await GreetService.Greet(name));
	} catch (err) {
		console.error(err);
	}
});

Events.On("time", (time) => {
	// The full RFC1123 stamp is too wide for the footer on a phone, so on narrow
	// screens (matching the CSS breakpoint) we show just the clock time.
	const full = time.data;
	const compact = (full.match(/\d{1,2}:\d{2}:\d{2}/) || [full])[0];
	timeElement.innerText = window.matchMedia("(max-width: 640px)").matches ? compact : full;
});
`
