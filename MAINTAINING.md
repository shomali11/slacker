Docs for slacker maintainers

## Versioning

Version numbers adhere to go [module versioning
numbering](https://go.dev/doc/modules/version-numbers).

Use the following criteria to identify which component of the version to change:

- If a breaking change is made, the major version must be incremented
- If any features are added, the minor version should be incremented
- If only fixes are added, the patch version should be incremented
- If no code changes are made, no version change should be made

When updating the major version we must also update our [go.mod](./go.mod)
module path to reflect this. For example the version 2.x.x module path should
end with `/v2`.

## Releases

Once all changes are merged to the master branch, a new release can be created
by performing the following steps:

- Identify new version number, `ie. v1.2.3`
- Use [golang.org/x/exp/cmd/gorelease](https://pkg.go.dev/golang.org/x/exp/cmd/gorelease) to ensure version is acceptable
    - If issues are identified, either fix the issues or change the target version
    - Example: `gorelease -base=<previous> -version=<new>`
- Tag commit with new version
- Push tag upstream

Once pushed, the [goreleaser](./.github/workflows/goreleaser.yaml) workflow is
triggered to create a new GitHub release from this tag along with a changelog
since the previous release.

### Changelogs

Changelog entries depend on commit subjects, which is why it is important that
we encourage well written commit messages.

Based on the commit message, we group changes together like so:

- `Features` groups all commits of type `feat`
- `Bug Fixes` groups all commits of type `fix`
- `Other` groups all other commits

Note that the `chore` and `docs` commit types are ignored and will not show up
in the changelog.

For more details on commit message formatting see the
[CONTRIBUTING](./CONTRIBUTING.md) doc.
