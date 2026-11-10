# Contributing

Thanks for helping improve the DLC Desktop App. Bug reports, fixes and
documentation improvements are all welcome.

## Reporting problems

- **Bugs and feature ideas:** [open an issue](https://github.com/deerwalkrnd/dlc-desktop-app/issues/new/choose).
  Include the version shown on the app's **Library** page.
- **A video doesn't show up or is in the wrong place:** use the
  "Video naming problem" issue template and paste the exact file name.
- **Security issues:** please don't open a public issue. See
  [SECURITY.md](SECURITY.md).

## Making a change

1. Fork the repository and create a branch from `main`.
2. Set up the tools: Go 1.27+ (the exact toolchain is pinned in `go.mod` and
   downloaded automatically), Node.js 22.18+, `make`, and optionally `ffmpeg`
   for `make demo`.
3. Make your change. The [developer guide](docs/docs/developers.md) explains
   the architecture.
4. Run the checks CI runs:

   ```sh
   make lint
   make test
   ```

5. Open a pull request describing what changed and why. Screenshots help for
   UI changes.

## Ground rules

- **Never break existing file names.** Every name that parsed in version 1.x
  must keep parsing the same way. `internal/catalog/legacy_test.go` and
  `FuzzLegacyCompatibility` enforce this; if your change needs to touch those
  tests, explain why in the pull request.
- **No new Go dependencies** without discussing it in an issue first. The
  backend uses only the standard library.
- **The app must work fully offline.** No CDNs, web fonts, analytics or other
  network requests from the frontend.
- **Keep Windows 10 working.** Most installations are Windows lab computers.
- Update the docs in `docs/docs/` when behaviour that users or developers see
  changes, and add an entry to [CHANGELOG.md](CHANGELOG.md).

## Code style

- Go: `gofmt`, `go vet`.
- Frontend: Prettier and ESLint (`npm run format` and `npm run lint` in
  `web/`), and `svelte-check` must pass.

By contributing you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).
