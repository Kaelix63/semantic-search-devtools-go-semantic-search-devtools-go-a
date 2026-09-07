package search

import "testing"

func TestRankLocalReleaseDiagnostics(t *testing.T) {
	cases := []struct {
		name, query string
		want        []string
	}{{"compiler terms", "compiler diagnostics", []string{"build", "diag"}}}
	docs := []Document{{ID: "build", Text: "build event records compiler diagnostics", Kind: "build"}, {ID: "release", Text: "release operation promotes a version", Kind: "release"}, {ID: "diag", Text: "developer-facing diagnostics include compiler output", Kind: "diagnostic"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := RankLocal(tc.query, docs)
			if got[0].ID != tc.want[0] || got[1].ID != tc.want[1] {
				t.Fatalf("ranked %v, want %v", got, tc.want)
			}
		})
	}
}
