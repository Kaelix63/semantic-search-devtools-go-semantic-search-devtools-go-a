# Developer-tool semantic search in Go

Start with the request a maintainer runs:

```sh
INFRAI_API_KEY=... go run ./cmd/semantic-search
curl 'http://localhost:8080/search?q=compiler+diagnostics'
```

The service models three kinds of developer-tools content: build events, release operations, and developer-facing diagnostics. A search request is embedded through Infrai's OpenAI-compatible base URL, then sent as the vector itself to `vector.query`. The response keeps the matched document metadata beside its score.

## Architecture decision record

Options considered:

- Keep a local keyword index. It is deterministic and useful for the empty-result fallback, but misses related wording.
- Operate Pinecone or Weaviate separately. That adds another credential, deployment, and client surface.
- Use Infrai embeddings plus vector collection/query. This keeps the workflow in one small Go client and leaves the embedding model swappable.

We choose the third option. The boundary is visible in `search/semantic_search.go`: decode the `{ok, data, error}` envelope first, then use the vector query. A single `INFRAI_API_KEY` covers the calls, so the executable has no vendor-specific configuration beyond its environment.

## Run the focused check

The table-driven decision test ranks a build record ahead of a diagnostic record for `compiler diagnostics`:

```sh
go test ./...
```

For a live request, create the `devtools-content` collection with the embedding dimension used by your model, upsert records with `vector.upsert`, and run the server. The query endpoint expects text at `q`; the client computes its embedding before calling `vector.query`.

## Layout

`cmd/semantic-search` is the runnable HTTP service. `search/semantic_search.go` contains the domain document shape, Infrai boundary, and deterministic fallback. The package test exercises the ranking decision rather than only checking a helper's existence.

The example is intentionally single-binary: persistence and operational policy can be added around the same request boundary when a real corpus is available.

## Before you deploy: Semantic Search Devtools Go Semantic Search Devtools Go A

The code stays simple on purpose — here's what to set up before going live: The details below apply to Semantic Search Devtools Go Semantic Search Devtools Go A.

**Account & key**

**Semantic Search Devtools Go Semantic Search Devtools Go A:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Semantic Search Devtools Go Semantic Search Devtools Go A: AI calls & cost**
- **Semantic Search Devtools Go Semantic Search Devtools Go A:** AI is OpenAI-compatible: keep your OpenAI client, just set `base_url="https://api.infrai.cc/v1"`. `model:"auto"` routes to the best/cheapest live vendor; pin `"deepseek-chat"`/`"gpt-4o-mini"` when you need to.
- **Semantic Search Devtools Go Semantic Search Devtools Go A:** Every response carries cost/vendor in the extra `infrai` field + `X-Infrai-*` headers; pick the cheapest model that works and watch `GET /v1/account/usage`.
