# Changelog

All notable changes are recorded here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses
[semantic versioning](https://semver.org/).

## [Unreleased]

## [2.0.0] - 2026-11-10

### Added
- Search across lesson titles, chapters, subjects, teachers, class and lesson
  numbers, with live results and highlighting (`/` jumps to the search box).
- Watch progress: resume where you left off, "Watched" marks, per-chapter
  progress and a **Continue watching** list on the home page.
- New video player with keyboard shortcuts, playback speed, previous/next
  lesson and clear error messages.
- `dlc-check` tool to validate a folder's file names before shipping it.
- **Library** page showing skipped files and warnings, with **Rescan library**.
- Videos in sub-folders, `.MP4` extensions and hyphens inside names
  (`Ram-Prasad`) are now accepted.
- Windows executables have the Deerwalk icon, version details (right-click →
  Properties → Details) and a Windows 10/11 manifest.
- The console window is titled with the app's address and a reminder that
  closing it stops the app.
- The project is open source under the MIT License.

### Changed
- Rewritten as a single executable with the web app embedded; no database.
  The catalogue is built from file names in memory at startup and on rescan.
- A lower-case `old` (or `Old`) syllabus now means the old syllabus. Version
  1.x treated it as new.
- New Deerwalk light-blue design across the app and player.
- Frontend upgraded to SvelteKit 3 and Svelte 5; backend to Go 1.27.

### Fixed
- Clicking inside the console window on Windows no longer freezes the app
  (QuickEdit mode is turned off).

### Removed
- Windows 7 support. Windows 10 or later is required.
- The `dlc.db` / `dlc.sqlite` database file is no longer created or used and
  can be deleted.

### Security
- The server only answers requests addressed to `localhost`/`127.0.0.1`
  (blocks DNS rebinding), rejects cross-site `POST` requests, sends no CORS
  headers and serves videos by ID through `os.Root`, so paths can't escape the
  library folder.

[Unreleased]: https://github.com/deerwalkrnd/dlc-desktop-app/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/deerwalkrnd/dlc-desktop-app/compare/1.1.0...v2.0.0
