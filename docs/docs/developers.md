---
sidebar_position: 5
title: Developer guide
---

# Developer guide

## Architecture

The app is a single Go binary with the SvelteKit web app embedded in it.

```
DLC/ folder ──scan──▶ in-memory catalogue ──JSON API──▶ SvelteKit SPA (in the browser)
                                         └─/media/{id}─▶ <video> (HTTP range requests)
```

- **No database.** Every piece of information comes from file names, so the
  catalogue is rebuilt from the folder at startup (milliseconds for thousands
  of files) and on **Rescan**. Rescans build a new immutable index and swap
  it in atomically.
- **Stable IDs.** Lesson IDs are a hash of the file path, and subject slugs
  come from the name and syllabus, so URLs survive rescans and restarts.
- **No third-party Go dependencies.** Only the standard library is used.

| Path                  | Purpose                                                        |
| --------------------- | -------------------------------------------------------------- |
| `cmd/dlc-desktop-app` | Entry point: flags, library path, browser launch, shutdown     |
| `cmd/dlc-check`       | CLI that validates a video folder                              |
| `internal/catalog`    | Pure filename parser (with the legacy-compatibility guarantee) |
| `internal/library`    | Folder scan → immutable index; safe file opening (`os.Root`)   |
| `internal/server`     | HTTP API, media streaming, SPA serving, security middleware    |
| `web/`                | SvelteKit app; `web/embed.go` embeds `web/build`               |
| `scripts/`            | Demo-library generator and release README                      |

## Building

Requirements: Go (version in `go.mod`), Node.js 22.18+, `make`, and optionally
`ffmpeg` for the demo library.

```sh
make            # web app + bin/dlc-desktop-app + bin/dlc-check
make test       # Go tests (race detector) + frontend unit tests
make lint       # go vet, gofmt, svelte-check, eslint, prettier
make demo       # creates ./DLC with playable sample lessons (needs ffmpeg)
make run        # build and start against ./DLC
make release    # dist/*.zip for Windows and Linux
```

Without the frontend build, the Go binary still compiles and serves the API.
The browser then shows a page explaining how to build the web app.

## Development workflow

Run the API server and the Vite dev server side by side:

```sh
make demo                                     # once
go run ./cmd/dlc-desktop-app -open=false      # API on :3000
cd web && npm run dev                         # UI with hot reload on :5173
```

Vite proxies `/api` and `/media` to the Go server (see `web/vite.config.ts`).

## The filename parser and compatibility

`internal/catalog` first splits names on every `-`, exactly like version 1.x.
Only if that fails does it split on ` - `. As a result, every name that parsed
in 1.x parses to the same result. `legacy_test.go` keeps a verbatim copy of the
old parser, and `FuzzLegacyCompatibility` checks this property on arbitrary
input:

```sh
go test ./internal/catalog -run '^$' -fuzz FuzzLegacyCompatibility -fuzztime 60s
```

The only intentional difference is that lesson numbers such as `NaN`, `Inf`
and negative numbers are now rejected. 1.x turned them into nonsense chapter
numbers.

## Video player

`web/src/lib/player.ts` is an explicit state machine
(`idle → loading → paused ⇄ playing ⇄ buffering → ended`, plus `error`).
`VideoPlayer.svelte` translates `<video>` events into machine events and
renders from the state. Run its unit tests with `npm test` in `web/`.

## HTTP API

All responses are JSON. Errors look like `{"error": "message"}`.

| Method & path                              | Returns                                                                   |
| ------------------------------------------ | ------------------------------------------------------------------------- |
| `GET /api/classes`                         | `{classes: [{number, subjectCount, lessonCount, syllabi}]}`               |
| `GET /api/classes/{n}[?syllabus=old\|new]` | `{number, syllabi, subjects: [{slug, name, syllabus, chapterCount, lessonCount}]}` |
| `GET /api/classes/{n}/subjects/{slug}`     | `{class, slug, name, syllabus, chapters: [{number, title, lessons: [{id, number, title, teacher, mediaUrl}]}]}` |
| `GET /api/teachers`                        | `{teachers: [{name, lessonCount, subjects}]}`                             |
| `GET /api/library`                         | Folder, counts, scan time, `skipped` and `warnings` lists, version        |
| `POST /api/library/rescan`                 | Rescans, then returns the same as `GET /api/library`                      |
| `GET /media/{id}`                          | The video file, with range-request support                                |
| `GET /healthz`                             | `{ok, app, version}`                                                      |

## Security model

The server binds to `127.0.0.1` by default and:

- rejects requests whose `Host` header isn't a loopback name, which blocks
  DNS-rebinding attacks from websites,
- sends no CORS headers, so other websites can't read the API,
- rejects cross-site `POST` requests (`Origin` / `Sec-Fetch-Site` checks),
- serves videos only by ID, opening files through `os.Root` so that paths,
  `..` and symlinks can't escape the library folder, and never lists
  directories,
- sends `nosniff`, `X-Frame-Options: DENY` and `Referrer-Policy` headers.

With `-addr 0.0.0.0:…` (LAN mode) the Host check is turned off, since clients
connect by IP address.

## Releasing

Push a tag such as `v2.0.0`. The `Release` workflow runs the tests, builds the
Windows and Linux archives and publishes a GitHub release.
