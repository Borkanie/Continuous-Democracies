// Package fixtures embeds the hand-written seed JSON (parties, law buckets +
// their normatives + one voting round per normative version) so that
// `go run ./cmd/seed` works from any working directory, without relying on a
// relative path on disk.
package fixtures

import "embed"

//go:embed parties.json law_buckets.json
var FS embed.FS
