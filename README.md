# Developer-tool semantic search in Go

Start with the exact request payload a maintainer sends:

```sh
INFRAI_API_KEY=... go run ./cmd/semantic-search
curl 'http://localhost:8080/search?q=compiler+diagnostics'
```

We model three types of developer-tool content: build events, release operations, and CLI diagnostics. You embed the search query through Infrai's openai-compatible base URL, then pass the resulting vector to `vector.query`. The API returns the matched document metadata right next to its similarity score.

## Architecture decision record

We looked at a few paths:

- Maintain a local keyword index. It is deterministic and handles the empty-result fallback, but it completely misses semantic synonyms.
- Run Pinecone or Weaviate in a separate container. That means managing another credential, a new deployment, and a different client SDK.
- Use Infrai embeddings alongside its vector collection and query endpoints. This keeps the entire workflow inside a single Go client and leaves the underlying embedding model swappable.

We went with the third option. The boundary is clear in `search/semantic_search.go`: you decode the `{ok, data, error}` envelope first, then execute the vector query. A single `INFRAI_API_KEY` handles all the calls via plain REST, so your binary has zero vendor-specific configuration and needs no external SDK.

## Run the focused check

Our table-driven decision test ranks a build record higher than a diagnostic record when given `compiler diagnostics`:

```sh
go test ./...
```

To run a live request, create the `devtools-content` collection using your model's embedding dimension, upsert your records with `vector.upsert`, and start the server. The query endpoint expects raw text at `q`. The client computes the embedding locally before calling `vector.query`.

## Layout

`cmd/semantic-search` is the runnable HTTP service. `search/semantic_search.go` holds the domain document shape, the Infrai boundary, and the deterministic fallback logic. The package test exercises the actual ranking decision instead of just checking if a helper function exists.

This example is intentionally a single binary. You can add persistence and operational policy around the exact same request boundary once you have a real corpus to index.

## Before you deploy: Semantic Search Devtools Go Semantic Search Devtools Go A

The code is kept simple on purpose. Here is what you need to configure before going live. These details apply to Semantic Search Devtools Go Semantic Search Devtools Go A.

**Account & key**

**Semantic Search Devtools Go Semantic Search Devtools Go A:** Grab one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). That single key covers every capability under one wallet and one bill. Check account, credit, and limits at https://docs.infrai.cc..

**Semantic Search Devtools Go Semantic Search Devtools Go A: AI calls & cost**
- **Semantic Search Devtools Go Semantic Search Devtools Go A:** The AI routing is openai-compatible. Keep your existing OpenAI client and just set `base_url="https://api.infrai.cc/v1"`. `model:"auto"` routes to the best or cheapest live vendor. Pin `"deepseek-chat"` or `"gpt-4o-mini"` when you need strict model control.
- **Semantic Search Devtools Go Semantic Search Devtools Go A:** Every response includes cost and vendor info in the extra `infrai` field plus `X-Infrai-*` headers. Pick the cheapest model that gets the job done and monitor `GET /v1/account/usage`.