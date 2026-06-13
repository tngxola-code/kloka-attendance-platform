# Active Branches

This document tracks the active long-lived branches in the Kloka repository.

## Branch Overview

| Branch | Purpose                                           | Status |
| ------ | ------------------------------------------------- | ------ |
| `main` | Primary development branch and project foundation | Active |

## Current State

The project is currently in its foundation phase.

The `main` branch contains:

* Initial platform architecture
* Domain structure
* Database migration framework
* Authentication foundation
* HTTP server setup
* OpenAPI contract
* Core workforce domain scaffolding

Additional feature and hotfix branches will be created as development progresses.

## Future Branch Types

Short-lived branches will follow the naming conventions below:

### Feature Branches

```text
feature/attendance
feature/disputes
feature/leave-management
feature/payroll-reports
```

### Hotfix Branches

```text
hotfix/authentication
hotfix/database-migration
```

These branches are created from `main` and merged back through Pull Requests.

## Notes

* `main` is the source of truth.
* All development starts from `main`.
* Feature branches should be deleted after merge.
* Hotfix branches should be deleted after merge.
* Long-lived development branches are intentionally avoided.
