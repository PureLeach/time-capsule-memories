# 0001 — Attachments bypass the API

**Status:** accepted

## Context

A capsule may carry up to three images of up to 5 MB each. Routing them through
the Go service would mean the API buffers multipart bodies, holds them in memory
or in a temporary file, and forwards them to the object store — doubling the
bandwidth and making request timeouts a function of the client's upload speed.

## Decision

The browser uploads directly to the object store. The API only issues a signed
POST target: `GET /generate-presigned-url?directory=<uuid>` returns a URL plus
the policy fields, and the browser submits the file as `multipart/form-data`.
All attachments of one capsule share a directory UUID, and the capsule row
stores that UUID in `files_folder_uuid`.

The signature covers a policy that pins the object key prefix, caps the body at
5 MB, and requires a content type starting with `image/`. Those limits are
therefore enforced by the storage server, not by the browser: the client-side
checks in `FormPage.vue` exist to give fast feedback, not to protect anything.
The target expires after five minutes.

## Consequences

- The API never touches attachment bytes on the way in, so upload size is
  bounded by storage configuration rather than by the service's memory.
- The object store must be reachable from the browser under the same hostname
  the backend signed for, which is why `MINIO_ENDPOINT` is a Traefik network
  alias in `docker-compose.yml` rather than an internal service name.
- Uploads are unauthenticated, like the rest of this service. Anyone who can
  reach the endpoint can consume storage, one signed 5 MB image at a time.
- Nothing garbage-collects the directories of capsules that were uploaded to but
  never submitted. A bucket lifecycle rule is the intended fix.
