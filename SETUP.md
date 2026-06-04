# Maintainer setup

One-time manual steps that require account/browser authorization (automation
can't perform these).

## Renovate (self-hosted)

Dependency updates run via the self-hosted **Renovate** workflow
(`.github/workflows/renovate.yml`) in this repo, which manages all three repos —
`tfkit/tfparams`, `tfkit/tfparams-action`, `tfkit/tfparams-orb`. No Mend App
install is needed.

1. Create a token with read + write + pull-request access to the three repos:
   - a **fine-grained PAT** (recommended): Contents = Read/Write,
     Pull requests = Read/Write, Workflows = Read/Write, scoped to the repos; or
   - a classic PAT with `repo` + `workflow`.
2. Add it as a repository secret named **`RENOVATE_TOKEN`**
   (Settings → Secrets and variables → Actions).
3. Renovate then runs weekly (Mon 02:00 UTC) and on demand
   (Actions → Renovate → Run workflow). Per-repo rules are in each repo's
   `.github/renovate.json` (non-major Go + GitHub Actions automerge once CI is
   green; major Go updates stay manual).

## CircleCI orb

See [`tfkit/tfparams-orb`](https://github.com/tfkit/tfparams-orb) → `SETUP.md`
for connecting CircleCI, the API token, and the `orb-publishing` context. The
namespace and orb are created automatically by the publish job.

## GitHub Marketplace (optional)

To list the GitHub Action on the Marketplace, edit a release of
[`tfkit/tfparams-action`](https://github.com/tfkit/tfparams-action) and tick
"Publish this Action to the GitHub Marketplace".
