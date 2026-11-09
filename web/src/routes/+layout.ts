// The app is a client-rendered SPA served by the Go binary: the data lives on
// the user's machine, so there is nothing to render on a server or at build
// time.
export const ssr = false;
export const prerender = false;
