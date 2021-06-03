# Changelog

## [v0.14.0] - 2021-05-28

### Added:

- New endpoints for data export in CSV formats


## [v0.13.1] - 2021-05-24

### Added:
- New route to the management-api of the study-service to run custom study rules


## [v0.13.0]

### Added

- Folder [docs](docs) with some initial content.

### Changed

- [SIGNIFICANT]: reCaptcha usaged can be configured for each instance separately. Also API accepts "Instance-Id" header attribute, that can be used to refer to a specific instance. Webapp core lib attempts to send the Instance-Id from version 1.0.17.
For more infos, see [here](docs/recaptcha-config.md)
Change is backwards compatible, previous env variables should still work.

- Updated dependencies (reflected in go.mod). Relevant internal changes:
    - Auto Email has label attribute.

- Small code improvements how endpoints handle for singup is implemented.