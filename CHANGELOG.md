# Changelog


## ??? - 2024-01-23

### Added

- New endpoint for running custom study rules retrospectively on previous responses. Optional filter parameters are participantIDs, participant status, survey keys and range parameters for submission time of responses.


## ??? - 2023-08-02

### Added

- New endpoint for getting responses in JSON Format including pagination infos.
- For existing endpoints of getting responses (CSV/JSON/wide/long) pagination option can be set by parameter.
- New endpoints for
  - save study rules
  - delete study rules
  - get current study rules
  - get history of study rules with options of time filter and pagination.

## [v1.4.0] - 2023-07-14

### BREAKING CHANGE

- Changed URL for getting all participant states for a study by status, because previous url is used now with pagination and sorting option. To get all participant states (optionally filtered by study status) use the new endpoint `/v1/data/{studyKey}/participants/all`.

### Added

- New endpoint for getting participant state by ID.
- New endpoint for getting participant states by query with pagination and sorting option.

### Changed

- Add the possibility for the management-api to require mutual TLS authetication from the clients. This is done by setting the environment variable `REQUIRE_MUTUAL_TLS` with `true` and providing the path to the certificate and key files in the environment variables `MUTUAL_TLS_SERVER_CERT` and `MUTUAL_TLS_SERVER_KEY` as well as the `MUTUAL_TLS_CA_CERT` for the CA certificate. The client certificate must be signed by the same CA as the server certificate.
- The renew token endpoint is also exposed by the management-api. This endpoint can be used to renew the access token of a user.

## [v1.3.0] - 2023-03-20

### Changed

- Update dependency for study-service
- Send info about profile and account deletion to the study-service, so this can be handled in the study rules.
- During user migration, set timestamps to the middle of the day to improve privacy.

## [v1.2.1] - 2023-02-27

### Changed

- Updated project dependencies, especially relevant the messaging service to the latest version (to include the "until" flag).

## [v1.2.0] - 2022-11-02

### BREAKING CHANGE

- Changes in the survey management API to apply the survey history model changes.
  - This version requires the study-service with minimum of 1.3.0 version

### Changed

- Updating go version to minimum of 1.17
- Updating project dependencies

## [v1.1.1] - 2022-06-30

### Changed

- Replacing remaining standard logger with custom logger to standardise log output across projects.
- When logging in through SAML protocol, role and instance ID checks are now case insensitive to make process less error prone.

## [v1.1.0] - 2022-06-03

### Added

- New endpoints for managing researcher notification subscriptions.

### Changed

- Updated project dependencies.

## [v1.0.2] - 2022-03-06

### Changed

- Fix bug in report fetching endpoint

## [v1.0.1] - 2022-03-15

### Changed

- Updated study-service api

## [v1.0.0] - 2022-03-08

### Added

- New endpoints for temporary participants.
- New endpoint to be able to download response export in a JSON format.
- New endpoints for handling confidential responses (remove for profiles/researcher access).
- New endpoint for uploading participant file.
- New endpoint for running custom rules for a single participant.

### Changed

- To allow workflows where registration is triggered from a survey, study flow methods don't require the account to be confirmed yet.
- Profile creation does not require verified emails - needed for entry flow through surveys.
- Migration endpoint extended to be able to create reports as well.

## [v0.15.2] - 2022-02-16

### Changed

- Add option to set max grpc message size through environment variable

## [v0.15.1] - 2021-07-28

- Update dependency for messaging service, to use recent API.

## [v0.15.0] - 2021-07-07

### New

- Added new form of authentication for management-api: the service can be connected with a SAML based identity provider. Upon successful authentication, it calls the user-management-service's new endpoint for logging in with external identity provider and display the token provided by the that.
In this current concept, this workflow is set up to work with the Python scripts by copy pasting that token string into the terminal. Documentation about how to configure this new SAML interface can be found at [docs/saml-config.md](docs/saml-config.md).

## [v0.14.1] - 2021-06-03

### Changed

- Migration endpoint can receive a list of profile names, that will be used to create profiles for the user, and submit a migration survey for each of the new profiles. Profile names can be defined through the optional attribute `profileNames`, which is expected to be a string array.
Also `oldParticipantID` is renamed to `oldParticipantIDs` and expects a string array.
If `profileNames` is populated, `oldParticipantIDs` must contain at least the same number of entries (and any values without matching profile name, will be ignored.)

## [v0.14.0] - 2021-05-28

### Added

- New endpoints for data export in CSV formats

## [v0.13.1] - 2021-05-24

### Added

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
