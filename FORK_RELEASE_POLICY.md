# Fork Release Policy

This repository is the SatenRuiko-Lv0 fork of MetaCubeX/mihomo.

## Repository Identity

- The GitHub repository name must stay `mihomo`.
- The repository must remain a GitHub-native fork of `MetaCubeX/mihomo`.
- `origin` points to `https://github.com/SatenRuiko-Lv0/mihomo.git`.
- `upstream` points to `https://github.com/MetaCubeX/mihomo.git`.

## Base Rules

- Stable releases must be based on the latest official mihomo release tag.
- The official stable code line is `Meta`; use `Meta` and official `v*` release tags for stable release work.
- Treat `Meta` as the current official stable branch for release sync, patch replay, build, and publishing.
- Do not use the official `main` branch as a release base. It is not the active stable release branch for this fork workflow and must be treated as stale/archive-only unless upstream policy changes.
- Alpha prereleases must be based on the official `Alpha` branch.
- Local changes are maintained as a small compatibility patch set and reapplied after each upstream sync.
- Do not publish stable builds from an Alpha snapshot.

## Update Source

- Core update URLs must point to `https://github.com/SatenRuiko-Lv0/mihomo`.
- Release channel updates use the fork's latest stable release.
- Alpha channel updates use the fork's single rolling `Prerelease-Alpha` release.
- Do not restore updater URLs to `MetaCubeX/mihomo` when rebuilding from upstream.

## Release Assets

- Publish only Android arm64-v8a binaries.
- Assets are direct `.bin` files, not `.gz`, `.zip`, Docker images, or multi-platform packages.
- Stable asset names use `mihomo-android-arm64-v8-<version>.bin`.
- Alpha asset names use `mihomo-android-arm64-v8-alpha-<short_sha>.bin`.
- Each release also publishes `version.txt` and may publish `checksums.txt`.

## Release Shape

- Stable releases mirror official tag names, for example `v1.19.27`.
- Alpha uses only one rolling prerelease: `Prerelease-Alpha`.
- Do not create ad hoc prerelease tags such as `*-legacy-*`.
- Delete stale assets before uploading replacements with the same release tag.

## Compatibility Patch Scope

- The legacy H1 changes are protocol-shape compatibility patches for testing.
- The patch set currently covers no-TLS VMess WebSocket and VMess TCP+HTTP request serialization.
- Keep the patch set narrow and easy to reapply after upstream changes.
