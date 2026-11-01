// Package fixtures embeds the hand-written seed JSON (parties, the real
// current Senate/Chamber of Deputies rosters, law buckets + their normatives
// + one voting round per normative version) so that `go run ./cmd/seed`
// works from any working directory, without relying on a relative path on
// disk.
//
// senators.json/deputies.json are scraped from senat.ro/FisaSenatori.aspx
// and cdep.ro's structura2015.de roster pages (real names, real parliamentary
// group, real electoral circumscription) - see cmd/seed/main.go's doc
// comment for what's real vs mocked.
package fixtures

import "embed"

//go:embed parties.json senators.json deputies.json law_buckets.json
var FS embed.FS
