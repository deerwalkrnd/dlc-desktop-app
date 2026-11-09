# DLC Desktop App documentation

The user and developer guide, built with [Docusaurus](https://docusaurus.io/)
and published to https://deerwalkrnd.github.io/dlc-desktop-app/ by the
`Deploy documentation` workflow on every push to `main`.

```sh
npm ci
npm start      # live preview at http://localhost:3000/dlc-desktop-app/
npm run build  # production build into build/ (fails on broken links)
```

Pages live in `docs/`; the sidebar is generated from their `sidebar_position`.
