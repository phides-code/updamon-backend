# updamon-backend

An AWS Lambda serving a JSON HTTP API backed by DynamoDB. Each URL path maps to one table and one resource type (`/computers` is the example implementation). Add resources by registering handlers on the same Lambda.

```
API Gateway → Lambda (gateway) → resource handler → repository → DynamoDB
```

The gateway authenticates, routes on the first path segment, and delegates. Each resource owns its entity, validation, HTTP handler, and DynamoDB code under `internal/<resource>/`. Cross-cutting rules live in `domain` and `platform`.

## Quick links

| I want to…                          | Go here                                      |
| ----------------------------------- | -------------------------------------------- |
| Run tests / local API               | [Development](#development)                  |
| Understand `/computers`             | [Computers](#computers-computers)            |
| Understand `/sitreps`               | [Sitreps](#sitreps-sitreps)                  |
| Add a field to an existing resource | [Adding a field](#adding-a-field)            |
| Add a new table / URL prefix        | [docs/new-resource.md](docs/new-resource.md) |

## Project layout

```
cmd/lambda/main.go       Lambda entry → app.Build
internal/
  app/                   Composition root (wire repos + Register)
  computer/              Reference vertical slice (copy this)
  sitrep/                Sitrep vertical slice
  domain/                Shared errors, UUID rules, string/IPv4 validation
  gateway/               Auth gate + first-segment routing
  platform/              Response envelope, error mapping, logging, CF token
  testutil/              Shared test helpers and resource fixtures
template.yml             SAM: API, Lambda, tables
Makefile                 test, build, local, deploy
```

**Dependency direction:** `app.Build` → `gateway` → handler → repository → DynamoDB. Handlers never call DynamoDB directly. Resources do not import each other.

Copy `internal/computer/` for a new resource. Keep composition stubs in `internal/app/*_stub_test.go` (Go cannot import another package’s tests).

## API contract

### Authentication

Every request except `OPTIONS` needs header `X-CF-Token` (`platform.CFTTokenHeader`). SAM parameter `AwsCfToken` maps to env `AWS_CF_TOKEN` (`platform.CFTTokenEnvVar`).

Under `sam local` (`AWS_SAM_LOCAL` = `true` or `1`), the token check is skipped.

### Response envelope

```json
{ "data": { ... } | [ ... ] | null, "error": "message" | null }
```

Success sets `data` and leaves `error` null. Failure does the opposite.

| HTTP | `error`                 | Domain sentinel       | When                         |
| ---- | ----------------------- | --------------------- | ---------------------------- |
| 400  | `invalid json`          | `ErrInvalidJSON`      | Body is not JSON             |
| 400  | `invalid id`            | `ErrInvalidID`        | Path `{id}` is not a UUID    |
| 400  | `validation failed`     | `ErrValidationFailed` | Domain rule failed           |
| 401  | `unauthorized`          | `ErrUnauthorized`     | Missing/wrong token          |
| 404  | `not found`             | `ErrNotFound`         | Missing item / unknown route |
| 405  | `method not allowed`    | `ErrMethodNotAllowed` | Unsupported method           |
| 409  | `already exists`        | `ErrAlreadyExists`    | Duplicate create             |
| 500  | `internal server error` | —                     | Unexpected failure           |

Client-facing text is the sentinel’s `Error()` string (`domain/errors.go`), mapped in `platform/errors.go`. Prefer `ErrValidationFailed` for field rules; add a new sentinel only when you need a new HTTP status or message for every resource.

### Computers (`/computers`)

| Method   | Path              | Behavior                                    |
| -------- | ----------------- | ------------------------------------------- |
| `GET`    | `/computers`      | List all                                    |
| `GET`    | `/computers/{id}` | Get by UUID                                 |
| `POST`   | `/computers`      | Create (`id` and `createdOn` set by server) |
| `PUT`    | `/computers/{id}` | Update client-writable fields               |
| `DELETE` | `/computers/{id}` | Hard delete; returns the deleted item       |

**Item** (list returns an array of the same shape):

```json
{
    "id": "uuid",
    "hostname": "string",
    "ip": "192.168.1.10",
    "os": "string",
    "kernel": "string",
    "model": "string",
    "ram": "string",
    "cpu": "string",
    "storage": "string",
    "createdOn": 1717516800000
}
```

**Create / update body:** `{ "hostname", "ip", "os", "kernel", "model", "ram", "cpu", "storage" }`

**Validation**

- `hostname`, `os`, `kernel`, `model`, `ram`, `cpu`, `storage`: required, 1–100 Unicode characters (`domain.DefaultMinStringLength`–`DefaultMaxStringLength`)
- `ip`: required dotted-quad IPv4 address (`domain.ValidateIPv4`)
- Path `{id}`: UUID, or 400 `invalid id`

List scans the full table. DynamoDB pagination stays inside the repository; it is not exposed over HTTP.

### Sitreps (`/sitreps`)

| Method   | Path             | Behavior                                    |
| -------- | ---------------- | ------------------------------------------- |
| `GET`    | `/sitreps`       | List all                                    |
| `GET`    | `/sitreps/{id}`  | Get by UUID                                 |
| `POST`   | `/sitreps`       | Create (`id` and `createdOn` set by server) |

**Item** (list returns an array of the same shape):

```json
{
    "id": "uuid",
    "hostname": "string",
    "createdOn": 1717516800000
}
```

**Create body:** `{ "hostname": "string" }`

**Validation**

- `hostname`: required, 1–100 Unicode characters (`domain.DefaultMinStringLength`–`DefaultMaxStringLength`)
- Path `{id}`: UUID, or 400 `invalid id`

List scans the full table. DynamoDB pagination stays inside the repository; it is not exposed over HTTP.

## Development

Requires Go 1.23+ and [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html).

```bash
make test      # unit tests + coverage gates (see Makefile)
make build
make local     # API on :8000; no auth header required
```

```bash
curl http://localhost:8000/computers
```

**Deploy:** `make init` once, then `export AWS_CF_TOKEN=… && make deploy`. CI (`.github/workflows/go.yml`) tests, builds, and deploys on push to `main`.

## Adding a field

Example: add `origin` to computer. Prefer TDD: failing test → smallest fix → green. Paths below use computer; substitute `<resource>` for another package.

| Step | File(s)                                  | Do this                                                                                                                                                                                                                                                                         |
| ---- | ---------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1    | `internal/computer/computer_test.go`     | Extend local `validCreateInput` / `validUpdateInput`. Add a case that blanks (or otherwise breaks) the new field.                                                                                                                                                               |
| 2    | `internal/computer/computer.go`          | Add the field on `Computer` with `json` + `dynamodbav` tags. If clients set it, add it to `CreateInput` / `UpdateInput` and validate (`domain.ValidateRequiredString` / `ValidateIPv4`, or custom). Server-owned fields are **not** on inputs — set them in the handler. |
| 3    | `internal/testutil/computer_fixtures.go` | Add the field to `ComputerBody`, `ValidComputerBody`, `ComputerWithID`, and list fixtures if needed.                                                                                                                                                                            |
| 4    | `internal/computer/fixtures_test.go`     | If the field is required on create/update, extend `newComputerValidationBodies` only when you need a new _shape_; reuse existing empty/whitespace/too-long fixtures when the rule matches hostname.                                                                             |
| 5    | `internal/computer/dynamodb.go`          | Add `Attr…` constant. If PUT-updatable, add it to the Update `SET` / names / values maps (keep attribute names alphabetical).                                                                                                                                                   |
| 6    | `internal/computer/assert_test.go`       | Add `Attr…` to the expected key list in `assertComputerDataKeys` (alphabetical).                                                                                                                                                                                                |
| 7    | `internal/computer/handler_test.go`      | Success create/update: assert the new field when it appears in the response. Client-error rows if validation can fail on this field.                                                                                                                                            |
| 8    | `internal/computer/handler.go`           | Add the field to `writePayload`; copy into create/update inputs and the entity passed to the repo.                                                                                                                                                                              |
| 9    | `internal/computer/dynamodb_test.go`     | If PUT-updatable: include the attr in the update success `AssertUpdateSets` map. Create/Get usually pick the field up via fixtures automatically.                                                                                                                               |
| 10   | `internal/computer/mocks_test.go`        | Only if a hand-built `Computer{…}` omits the new field and a test compares full structs.                                                                                                                                                                                        |
| 11   | `README.md`                              | Update the computers item shape, create/update bodies, validation, and PUT behavior row.                                                                                                                                                                                        |

Skip DynamoDB Update changes (steps 5 and 9) for read-only or create-only fields. Run `make test` before opening a PR.

## Adding a new resource

Copy `internal/computer/` and follow the file-by-file checklist in **[docs/new-resource.md](docs/new-resource.md)**.

Use `domain.ErrValidationFailed` unless you add a new cross-cutting sentinel (see the errors table above). Do not put resource routes in `gateway_test.go` — keep those in the resource’s `router_test.go`.
