package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"semantic-search-devtools/search"
)

func main() {
	if os.Getenv("INFRAI_API_KEY") == "" {
		log.Fatal("set INFRAI_API_KEY")
	}
	client, err := search.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	docs := []search.Document{{ID: "build-142", Text: "build event: compiler diagnostics for a failed module", Kind: "build", Release: "2026.09"}, {ID: "release-31", Text: "release operation promotes a version after checks", Kind: "release", Release: "2026.09"}, {ID: "diag-8", Text: "developer-facing diagnostics explain missing flags", Kind: "diagnostic", Release: "2026.09"}}
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "q is required", http.StatusBadRequest)
			return
		}
		emb, err := client.Embedding(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		matches, err := client.Query("devtools-content", emb, 5)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if len(matches) == 0 {
			for i, d := range search.RankLocal(query, docs) {
				matches = append(matches, search.Result{ID: d.ID, Score: float64(len(docs) - i), Metadata: d})
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"query":%q,"results":%v}`+"\n", query, matches)
	})
	log.Println("listening on http://localhost:8080/search?q=vector+search+endpoint")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
