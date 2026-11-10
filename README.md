# DLC Desktop App

Offline player for Deerwalk Learning Center lesson videos. It's a single
executable that scans a folder of `.mp4` files, organises them by class,
subject, chapter and lesson from their file names, and serves them to the
browser at `http://localhost:3000`.

**Full guide:** https://deerwalkrnd.github.io/dlc-desktop-app/

## Quick start (users)

1. Download the Windows zip from
   [Releases](https://github.com/deerwalkrnd/dlc-desktop-app/releases)
   (Windows 10 or later) and extract it.
2. Put your videos in the `DLC` folder next to `dlc-desktop-app.exe`.
3. Double-click `dlc-desktop-app.exe`.

Video names follow this format (sub-folders are fine):

```
TEACHER - LESSON NO - CHAPTER TITLE - LESSON TITLE - SUBJECT - OLD|NEW - CLASS.mp4

Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8.mp4
```

Added videos? Press **Rescan library** on the app's Library page. To check a
folder's file names before shipping it, run `dlc-check <folder>`.

## Development

```sh
make          # build web app + bin/dlc-desktop-app + bin/dlc-check
make test     # Go + frontend tests
make lint     # vet, gofmt, svelte-check, eslint, prettier
make demo     # create ./DLC with playable sample lessons (needs ffmpeg)
make run      # build and run against ./DLC
```

For hot reload, run `go run ./cmd/dlc-desktop-app -open=false` and
`npm run dev` in `web/` side by side.

See the [developer guide](docs/docs/developers.md) for the architecture, the
API and the filename-compatibility guarantee.

## Contributing

Bug reports and pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md),
and report security issues privately as described in [SECURITY.md](SECURITY.md).
Changes are listed in [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)
