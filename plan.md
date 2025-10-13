# Refactor Plan: Streamline File Upload Routines

## Current Flow Snapshot

- `BotAPI.UploadFilesWithContext` streams multipart bodies through an `io.Pipe`, mixes parameter writing with file upload, and re-implements response decoding logic already present in `MakeRequestWithContext` (`bot.go:178`-`bot.go:269`).
- `BotAPI.RequestWithContext` branches between the multipart path and the url-encoded path via `hasFilesNeedingUpload`, mutating params in-place when no uploads are required (`bot.go:331`-`bot.go:353`).
- Every `Fileable` config builds slices of `RequestFile`, pairing field names with `RequestFileData` implementations (`configs.go:143`-`configs.go:168`, e.g. `PhotoConfig.files` at `configs.go:485`-`configs.go:503`).
- `RequestFileData` implementers differentiate upload sources vs. remote references via `NeedsUpload`, but enforce correctness with `panic`-backed guard methods (`configs.go:170`-`configs.go:226`).
- Media group helpers duplicate file/param preparation (`configs.go:3460`-`configs.go:3526`), mirroring the logic scattered in per-config `files()` methods.

## Pain Points & Constraints

| Area | Observation | Impact |
| --- | --- | --- |
| File source abstraction | `RequestFileData` blends two lifecycles (upload vs. reference) and depends on run-time panics when the wrong accessor is used. | Hard to reason about, error-prone API for contributors, difficult to extend (e.g. new remote source types). |
| Request execution | Multipart and urlencoded paths duplicate request creation, logging, and error handling. Multipart builder is buried in a goroutine, complicating diagnostics and tests. | Higher maintenance cost, inconsistent debugging output, limited surface for unit testing. |
| Config proliferation | Each config manually assembles `[]RequestFile`, re-encoding the same conditional thumb / media handling. | Boilerplate everywhere; changes to the upload contract require sweeping edits. |
| Media group helpers | `prepareInputMediaForParams` and `prepareInputMediaForFiles` maintain parallel logic for the same state transitions. | Increases coupling between steps; diverging behavior is easy to introduce accidentally. |
| Testing depth | Upload pipeline largely covered by integration-style tests; no focused tests for multipart generation edge cases. | Refactors risky without targeted assertions; regressions in param encoding or file ordering may slip through. |

## Refactor Goals

1. **Single request pipeline**: centralize body selection (url-encoded vs. multipart) and response handling so both code paths share logging, error wrapping, and context propagation.
2. **Explicit file source model**: replace panic-driven `RequestFileData` with a type-safe structure that encodes upload vs. reference behavior and exposes helpers for serialization.
3. **Declarative file assembly**: give configs and media helpers a small builder API so they describe intent instead of manipulating slices directly.
4. **Composable multipart writer**: isolate multipart assembly (field writer, file part opener, resource closing) behind a reusable helper to simplify future features (e.g. streaming progress, instrumentation).
5. **Improved test surface**: add focused tests that assert multipart payload layout, field propagation, and error propagation without hitting the network.

## Implementation Plan

1. **Introduce core abstractions**
   - Create a `FileSource` struct (or interface pair) expressing `Kind` (`Upload`, `FileID`, `URL`, `Attach`) plus typed payload. Provide constructors replacing `FileBytes`, `FileReader`, `FilePath`, etc., but preserve existing API surface via thin adapters to avoid breaking callers immediately.
   - Implement helper methods (`ToParamValue`, `OpenUpload`) so consumers no longer call `NeedsUpload`/`SendData` directly.
   - Add unit tests covering each source flavor (missing today in `types_test.go` beyond interface assertions).

2. **Centralize request building**
   - Extract a `buildRequestPayload` helper returning a struct with headers and `io.ReadCloser`. Reuse it from both current request paths, preserving streaming behavior (keep `io.Pipe`, but move goroutine contents into a dedicated type with deterministic error returns).
   - Update `MakeRequestWithContext` and `UploadFilesWithContext` to delegate to this helper, eliminating duplicated logging / response parsing code.

3. **Rework `Request` decision logic**
   - Replace `hasFilesNeedingUpload` + manual param mutation with a higher-level `PayloadDescriptor` returned by configs (e.g. `{Params, Uploads, Inline}`), so the caller simply hands it to the request builder.
   - Ensure the decision boundary (multipart vs urlencoded) occurs in one place, using the new descriptor.

4. **Declarative config helpers**
   - Provide a small utility (e.g. `fileBuilder := UploadBuilder()`) that exposes methods such as `AddUpload(field, source)` and `AddReference(field, value)`.
   - Update representative configs (`PhotoConfig`, `AudioConfig`, `WebhookConfig`, media groups) to use the builder, reducing duplication and clarifying thumb handling.
   - For media groups, merge `prepareInputMediaForParams` and `prepareInputMediaForFiles` into a single routine that returns both the transformed media slice and associated file uploads, relying on the new builder.

5. **Testing & validation**
   - Add table-driven tests covering: mixed upload/reference payloads, failure to open files, reader closing semantics, and attach:// naming for media groups.
   - Run existing test suite plus `go vet`; ensure mocks for `HTTPClient` cover both request paths.
   - Provide temporary integration test (under `tests/`) or example verifying compatibility with Telegram API for at least one upload endpoint.

6. **Migration & cleanup**
   - After new abstractions land, deprecate direct use of old `RequestFileData` constructors, document migration in `docs/internals/uploading-files.md`.
   - Remove legacy helpers once downstream usage updated, ensuring public API changes are communicated via changelog.

## Risks & Mitigations

- **API compatibility**: Refactoring `RequestFileData` touches exported types. Mitigate by shipping adapters and marking them as deprecated only after confirming downstream adoption paths.
- **Multipart regressions**: Streamed uploads are sensitive to ordering and boundary handling. Mitigate via new unit tests and by preserving the existing streaming approach within the rebuilt helper.
- **Refactor scope creep**: Many configs share patterns; tackle them incrementally, updating one category (e.g. photo/audio/video) per PR to keep reviews manageable.

## Open Questions

- Should we introduce telemetry hooks (e.g. upload progress) while touching the multipart writer, or keep scope limited?
- Are there consumers relying on the exact panic messages from current `RequestFileData` implementations?
- Would switching to a generated boundary per request break any upstream mocks relying on fixed strings?

