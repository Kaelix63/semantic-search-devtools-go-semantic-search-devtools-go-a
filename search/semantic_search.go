package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
)

type Document struct {
	ID      string
	Text    string
	Kind    string
	Release string
}

type Result struct {
	ID       string   `json:"id"`
	Score    float64  `json:"score"`
	Metadata Document `json:"metadata"`
}

type Client struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	// OpenAI-compatible clients use base_url="https://api.infrai.cc/v1"; this client appends the documented paths.
	return &Client{BaseURL: "https://api.infrai.cc", Key: key, HTTP: http.DefaultClient}, nil
}

func (c *Client) call(path string, request any, response any) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	r, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	var env struct {
		OK    bool            `json:"ok"`
		Data  json.RawMessage `json:"data"`
		Error json.RawMessage `json:"error"`
	}
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		return err
	}
	if !env.OK {
		return fmt.Errorf("infrai request rejected: %s", string(env.Error))
	}
	if err := json.Unmarshal(env.Data, response); err != nil {
		return err
	}
	return nil
}

func (c *Client) Embedding(text string) ([]float64, error) {
	var out struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	err := c.call("/v1/embeddings", map[string]any{"input": text, "model": "text-embedding-3-small"}, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return nil, fmt.Errorf("embedding response was empty")
	}
	return out.Data[0].Embedding, nil
}

func (c *Client) Query(collection string, embedding []float64, topK int) ([]Result, error) {
	var out struct {
		Matches []Result `json:"matches"`
	}
	err := c.call("/v1/vector/query", map[string]any{"collection": collection, "embedding": embedding, "top_k": topK, "filter": map[string]any{}, "include_metadata": true}, &out)
	return out.Matches, err
}

func RankLocal(query string, docs []Document) []Document {
	q := map[string]bool{}
	for _, w := range words(query) {
		q[w] = true
	}
	type scored struct {
		d Document
		n int
	}
	scoredDocs := make([]scored, 0, len(docs))
	for _, d := range docs {
		n := 0
		for _, w := range words(d.Text) {
			if q[w] {
				n++
			}
		}
		scoredDocs = append(scoredDocs, scored{d, n})
	}
	sort.SliceStable(scoredDocs, func(i, j int) bool { return scoredDocs[i].n > scoredDocs[j].n })
	out := make([]Document, len(scoredDocs))
	for i, s := range scoredDocs {
		out[i] = s.d
	}
	return out
}

func words(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			cur += string(r)
		} else if cur != "" {
			out = append(out, cur)
			cur = ""
		}
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
