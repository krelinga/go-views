# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`github.com/krelinga/go-views` provides read-only views over Go collections. It has no dependencies beyond the standard library, and no `main` package — nothing to run, only tests.

The repo holds **two modules**, one per major version, following Go's major-subdirectory strategy:

- `github.com/krelinga/go-views` — v1, package `views`, files at the repo root. Described below.
- `github.com/krelinga/go-views/v2` — v2, package `views`, files in `v2/` with its own `go.mod`. An in-progress redesign; the sections below describe v1 unless stated otherwise.

## Commands

```bash
go test ./...                          # v1 only — does NOT reach v2/ (separate module)
cd v2 && go test ./...                 # v2
go test -run TestDictOfMap ./...       # run a single test
go test -run 'TestDictOfMap/nil_map'   # run a single subtest (spaces in name -> underscores)
go vet ./...
gofmt -l .                             # list unformatted files
```

## Architecture

`views.go` defines the three interfaces; everything else implements them.

- `List[T any]` — `Len()`, `Values() iter.Seq[T]`. Iteration order is explicitly *not* guaranteed by the contract.
- `Bag[T comparable]` — embeds `List[T]`, adds `Has(T) bool`.
- `Dict[K comparable, V any]` — embeds `List[V]` (so `Values()` yields the *values*), adds `Keys()`, `Get(K) (V, bool)`, `Has(K) bool`, `All() iter.Seq2[K, V]`.

Implementations follow the naming pattern `<Interface>Of<Backing>` (`ListOfSlice`, `BagOfMapKeys`, `DictOfMap`, …), one type per file named in snake_case after the type, with its test in a sibling `_test.go`. Because `Bag` is a superset of `List` and `Dict` of both, one struct often satisfies several interfaces; `ListOfMapKeys` is a type alias for `BagOfMapKeys` rather than a separate type.

Conventions that new implementations must preserve:

- Structs are thin wrappers with a single exported backing field — `S []T` for slices, `M map[K]V` for maps — and all methods use **value receivers**. Callers construct them inline: `views.DictOfMap[string, int]{M: m}`.
- The zero value (nil backing) must behave as an empty collection, never panic. Map-backed `Has`/`Get` explicitly nil-check before indexing.
- Prefer delegating to `slices`/`maps` stdlib helpers (`slices.Values`, `maps.All`, …) over hand-written loops.
- Doc comments on `Len`/`Has`/`Get` state the complexity (e.g. "This is O(n) in complexity."), which differs meaningfully across backings — `BagOfSlice.Has` is O(n), `BagOfMapKeys.Has` is O(1).

## Test conventions

Tests live in package `views_test` and import the package by its full module path. They are table-driven, and the table field holding the value under test is declared as the **interface** type (`views.List[string]`, `views.Dict[string, int]`), so the table entries double as a compile-time conformance check. Every table covers three cases: non-empty, empty, and nil backing. Since iteration order is unspecified, collected results and the expected slices are both sorted before `reflect.DeepEqual` — and expected values for empty cases are `nil`, matching what `slices.Collect` returns over an empty sequence.

## Releasing

Pushing a `v*` tag triggers `.github/workflows/tag.yaml`. It derives the target module from the tag's major version (`v0`/`v1` → repo root, `vN` → the `vN/` subdirectory), verifies that directory's `go.mod` declares the matching module path, tests every module in the repo, creates a GitHub release with generated notes, and warms the module proxy via `go list -m` from the released module's directory.

Tags carry no subdirectory prefix — v2 releases are tagged plain `v2.0.0`, not `v2/v2.0.0`, because the major-version suffix is excluded from a module's subdirectory name.
