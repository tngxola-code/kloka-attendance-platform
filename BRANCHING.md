# Branching Strategy

Kloka follows a lightweight GitHub Flow model designed to keep the main branch deployable while enabling rapid feature development.

## Branches

### main

The `main` branch contains production-ready code.

Rules:

* Always deployable
* Protected branch
* No direct commits
* Changes merged through Pull Requests only
* All automated checks must pass before merge

### feature/*

Feature branches are used for new functionality, enhancements, refactoring, and documentation updates.

Examples:

```text
feature/attendance-aggregation
feature/payroll-reporting
feature/bcea-leave-balances
feature/risk-rules
feature/openapi-updates
```

Branch from:

```text
main
```

Merge back into:

```text
main
```

### hotfix/*

Hotfix branches are used for urgent production fixes.

Examples:

```text
hotfix/jwt-expiry-validation
hotfix/payroll-download-bug
hotfix/geofence-calculation
```

Branch from:

```text
main
```

Merge back into:

```text
main
```

## Pull Requests

All changes must be submitted through Pull Requests.

Requirements:

* Clear title and description
* Build passes
* Tests pass
* No unresolved review comments
* Documentation updated when required
* OpenAPI updated when API contracts change

Recommended PR title format:

```text
feat: add attendance aggregation endpoint
fix: resolve payroll report generation failure
refactor: simplify worker repository queries
docs: update installation guide
```

## Commit Messages

Follow Conventional Commits where possible.

Examples:

```text
feat: implement trust scoring engine
feat: add BCEA leave balances

fix: resolve duplicate worker validation
fix: handle missing clock-out events

refactor: simplify attendance aggregation

docs: update README

test: add payroll integration tests
```

## Releases

Releases are created from `main`.

Recommended version format:

```text
v1.0.0
v1.1.0
v1.2.0
```

Versioning follows Semantic Versioning:

* MAJOR: Breaking API changes
* MINOR: New backward-compatible functionality
* PATCH: Bug fixes

## Development Workflow

1. Pull the latest `main`
2. Create a feature branch
3. Implement changes
4. Add or update tests
5. Open a Pull Request
6. Complete review and validation
7. Merge into `main`

Example:

```bash
git checkout main
git pull origin main

git checkout -b feature/attendance-exceptions

git add .
git commit -m "feat: implement attendance exceptions"

git push origin feature/attendance-exceptions
```
