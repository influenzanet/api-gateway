# Survey response downloads

These endpoints buffer the whole export and then send it with a `Content-Length`:

- `/v1/data/:studyKey/survey/:surveyKey/response` (wide CSV)
- `/v1/data/:studyKey/survey/:surveyKey/response/long-format` (long CSV)
- `/v1/data/:studyKey/survey/:surveyKey/response/json` (flat JSON)

Add `stream=true` to send chunks as they arrive, without buffering in the gateway.
Errors before the first chunk still return a JSON error. If the export fails after
the download has started, the gateway aborts the connection. Clients must treat an
incomplete transfer (e.g. `unexpected EOF`) as a failed export.
