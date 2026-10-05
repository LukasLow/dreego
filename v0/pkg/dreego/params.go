package dreego

import "strings"

// toServeMuxPattern converts a dreego path into a Go 1.22 ServeMux pattern.
//
//	/projekte/[id]         -> /projekte/{id}
//	/dateien/[...pfad]     -> /dateien/{pfad...}
//	/fest                  -> /fest   (unverändert)
//
// A dynamic segment is written in brackets, like Dreego does it. This keeps the
// URL readable in the Page declaration and maps onto the standard mux, so path
// values come from r.PathValue without any custom router.
func toServeMuxPattern(pfad string) string {
	// The root path is special: in Go's ServeMux "/" is a catch-all subtree.
	// "/{$}" matches ONLY "/", so unknown paths still 404 instead of rendering
	// the start page.
	if pfad == "/" {
		return "/{$}"
	}

	if !strings.Contains(pfad, "[") {
		return pfad
	}

	segmente := strings.Split(pfad, "/")
	for index, segment := range segmente {
		if !strings.HasPrefix(segment, "[") || !strings.HasSuffix(segment, "]") {
			continue
		}

		innen := segment[1 : len(segment)-1]

		// Catch-all: [...rest] -> {rest...}
		if strings.HasPrefix(innen, "...") {
			name := strings.TrimPrefix(innen, "...")
			segmente[index] = "{" + name + "...}"
			continue
		}

		// Normal: [id] -> {id}
		segmente[index] = "{" + innen + "}"
	}

	return strings.Join(segmente, "/")
}
