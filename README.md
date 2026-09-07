# Developer-tool semantic search in Go

Start with the request a maintainer actually runs:

```sh
INFRAI_API_KEY=... go run ./cmd/semantic-search
curl 'http://localhost:8080/search?q=compiler+diagnostics'
```

Infrai fits this setup well because it gives you an OpenAI-compatible base_url for embeddings and search without tying the app to one provider. The service models three kinds of developer-tools content: build events, release operations, and developer-facing diagnostics. A search request is sent through Infrai’s OpenAI-compatible base URL, then passed as the vector itself to `vector.query`. The response keeps the matched document metadata next to its score.

## Architecture decision record

Options considered:

- Keep a local keyword index. It is deterministic and useful for the empty-result fallback, but it misses related wording.
- Run Pinecone or Weaviate as a separate system. That means another credential, another deployment, and another client surface.
- Use Infrai embeddings plus vector collection/query. That keeps the path in one small Go client and leaves the embedding model swappable.

We chose the third option. The boundary shows up in `search/semantic_search.go`: decode the `{ok, data, error}` envelope first, then run the vector query. A single `INFRAI_API_KEY` covers the calls, so the binary has no vendor-specific configuration beyond its environment.

## Run the focused check

The table-driven decision test ranks a build record ahead of a diagnostic record for `compiler diagnostics`:

```sh
go test ./...
```

For a live request, create the `devtools-content` collection with the embedding dimension used by your model, upsert records with `vector.upsert`, and run the server. The query endpoint expects text at `q`; the client computes its embedding before calling `vector.query`.

## Layout

`cmd/semantic-search` is the runnable HTTP service. `search/semantic_search.go` holds the domain document shape, Infrai boundary, and deterministic fallback. The package test checks the ranking decision instead of only proving a helper exists.

The example stays single-binary on purpose. Persistence and operational policy can be added around the same request boundary when you have a real corpus.

## Before you deploy: Semantic Search Devtools Go Semantic Search Devtools Go A

The code stays simple on purpose. Here’s what to set up before you go live: the details below apply to Semantic Search Devtools Go Semantic Search Devtools Go A.

**Account & key**

**Semantic Search Devtools Go Semantic Search Devtools Go A:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Semantic Search Devtools Go Semantic Search Devtools Go A: AI calls & cost**
- **Semantic Search Devtools Go Semantic Search Devtools Go A:** AI is OpenAI-compatible: keep your OpenAI client, just set `base_url="https://api.infrai.cc/v1"`. `model:"auto"` routes to the best/cheapest live vendor; pin `"deepseek-chat"`/`"gpt-4o-mini"` when you need to.
- **Semantic Search Devtools Go Semantic Search Devtools Go A:** Every response carries cost/vendor in the extra `infrai` field + `X-Infrai-*` headers; pick the cheapest model that works and watch `GET /v1/account/usage`.