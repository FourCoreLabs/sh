# FourCoreLabs/sh

A fork of [mvdan/sh](https://github.com/mvdan/sh) (BSD-3-Clause, see
`LICENSE`) for embedding the interpreter in sandboxed agents through
[FourCoreLabs/go-bash](https://github.com/FourCoreLabs/go-bash).

```
go get github.com/FourCoreLabs/sh@v1.0.0
```

The module path is `github.com/FourCoreLabs/sh`, so consumers require it
directly. No `replace` directive is needed, and none should be added.
Versions start at v1.0.0 and are independent of upstream's v3 numbering.

## Base

Upstream **v3.14.1** (`a3f0c75d21d9`), imported without history.
`git diff <import commit> HEAD` shows everything this fork changes.

| Commit | Change |
|---|---|
| Import mvdan/sh v3.14.1 | Upstream tree, verbatim |
| Declare module github.com/FourCoreLabs/sh | Module path and import rewrite. `moreinterp` (a separate upstream module) is left out |
| interp: add VFSPaths | POSIX path resolution for interpreters over a virtual filesystem |
| interp: contain panics in goroutine stages | `recoverStage`: a panic in a pipeline stage, background job or process substitution fails that stage (status 2) instead of the process |

## History of the fork

- The earlier `FourCoreLabs/sh` tag `v3.13.2-attack1` (on v3.13.1) carried
  VFSPaths plus literal quoted heredocs. Upstream v3.14 keeps quoted heredoc
  bodies literal itself, so only VFSPaths was carried over.
- `recoverStage` was an unreleased local patch in ATTACK's
  `third_party/mvdan-sh`.
- v3.14's `pattern` package fixes the `[![:space:]]` regexp panic (agentworld
  KNOWN_ISSUES #15). That fix is upstream code, so no patch is needed.

## Updating from upstream

There is no shared history, so updates are applied as diffs:

1. In an upstream clone: `git diff v3.14.1 <new tag> -- . ':!moreinterp' > up.patch`
2. Rewrite paths: `sed -i 's#mvdan.cc/sh/v3#github.com/FourCoreLabs/sh#g' up.patch`
3. Here: `git apply -3 up.patch`, resolve, run `go test ./...`, and commit as
   "Update to mvdan/sh <new tag>". Update the Base section above.

Patches that upstream would accept (recoverStage is a candidate) are worth
sending upstream, so this fork shrinks over time.
