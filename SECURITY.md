# Security policy

## Supported versions

Only the latest release receives security fixes.

## Reporting a vulnerability

Please report vulnerabilities privately through GitHub:
**Security → Report a vulnerability** on this repository
([direct link](https://github.com/deerwalkrnd/dlc-desktop-app/security/advisories/new)).
Don't open a public issue.

Include the version, your operating system, steps to reproduce and what an
attacker could achieve. We aim to acknowledge reports within a week.

## Scope

The app is a local server meant to be reached only from the same computer
(`127.0.0.1`). Things we consider vulnerabilities include:

- reading files outside the video library folder,
- a website being able to read the app's API or trigger actions in it
  (for example through DNS rebinding or cross-site requests),
- crashing the app or exhausting resources with crafted file names or
  requests.

Running the app with `-addr 0.0.0.0:…` deliberately shares the videos with the
whole network; anyone on that network watching them is expected behaviour.

The [developer guide](docs/docs/developers.md#security-model) describes the
protections the server has.
