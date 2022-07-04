Docs for slacker maintainers

## Releases

Releases of Slacker are handled by [goreleaser](https://goreleaser.com) and
Github actions. Simply tagging a release with a semver compatible version tag
(ie. vX.Y.Z) and pushing the tag will trigger a Github action to generate a
release. See the goreleaser [config](.goreleaser.yaml) and Github
[workflow](.github/workflows/.goreleaser.yaml) files.

### Changelogs

goreleaser handles generating our changelog based on the commit subject of each
commit.

Commits that start with `feat:` are grouped into a "Features" section, while
those that start with `fix:` will be grouped into a "Bug fixes" section. Commits
that begin with `chore:` or `docs:` will be excluded, and all others will be
added to an "Others" section in the changelog.

When reviewing pull requests or committing code, it is strongly encouraged to
use one of the aformentioned prefixes so that changelogs are nicely formatted
and organized.

## Commit Messages

To maintain a tidy changelog on release, we should encourage the use of the
following commit subject prefixes (see the Changelogs for details on how they
are used)

- `feat`: New features
- `fix`: Bug fixes
- `docs`: Usage documentation changes (ie. examples, README)
- `chore`: Housekeeping taks that don't touch code or usage docs
