# REST API

Authenticate every API request with **HTTP Basic Auth** using the dashboard credentials from `DASHBOARD_CLIENT_USERNAME` and `DASHBOARD_CLIENT_PASSWORD`.

Unless noted otherwise, responses use **`Content-Type: application/json`**.

Timestamps in JSON bodies are **Unix time in seconds** (integer).

**Table of contents:**

- [REST API](#rest-api)
  - [Headers](#headers)
  - [Check](#check)
  - [Get All Messages](#get-all-messages)
  - [Schedule Message](#schedule-message)
  - [Retry Message](#retry-message)
  - [System Errors](#system-errors)

## Headers

| Field           | When                         | Required | Notes                                                                 |
| --------------- | ---------------------------- | -------- | --------------------------------------------------------------------- |
| `Authorization` | All documented endpoints     | Yes      | `Authorization: Basic <base64(username:password)>`                      |
| `Content-Type`  | `POST` requests with a body | Yes      | Must be `application/json`.                                           |

`GET` requests only need a valid `Authorization` header.

## Check

GET: `/check`

Verifies that the system is up and running.

**Example:**

```http
GET /check HTTP/1.1
Host: localhost:9866
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
      "default_numbers": [
        "120363352351961275@g.us"
      ]
    },
    "ts": 1735432224
}
```

[Back to Top](#rest-api)

## Get All Messages

GET: `/messages`

Returns messages stored by the scheduler, including status and timestamps.

### Query parameters

| Parameter | Required | Description |
| --------- | -------- | ----------- |
| `status`  | No       | Filter by status: `scheduled`, `sent`, or `failed`. Omit to return all messages. |

### Known limitation

Filtering with **`status=failed` currently returns `400` with `invalid status`** due to a redundant validation check in `serveGetMessages` (`internal/driver/rest.go`). Intended behavior is to return only failed messages. Track progress in [issue #17](https://github.com/ghazlabs/wa-scheduler/issues/17).

**Example (all messages):**

```http
GET /messages HTTP/1.1
Host: localhost:9866
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
```

**Example (scheduled only):**

```http
GET /messages?status=scheduled HTTP/1.1
Host: localhost:9866
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
```

The shape of each object in `data` is illustrated below (four separate valid JSON examples).

Scheduled (not sent yet; `sent_at` is null):

```json
{
    "id": "1da2f3e4-5b6c-7d8e-9a0b-c1d2e3f4g5h6",
    "content": "Job alert for Software Engineer at Invertase...",
    "recipient_numbers": ["120363352351961275@g.us"],
    "scheduled_sending_at": 1735432224,
    "sent_at": null,
    "retried_count": 0,
    "status": "scheduled",
    "reason": null,
    "created_at": 1735432224,
    "updated_at": 1735432224
}
```

Sent (`sent_at` set):

```json
{
    "id": "2b3c4d5e-6f7g-8h9i-0j1k-l2m3n4o5p6q7",
    "content": "Job alert for Software Engineer at dev.to...",
    "recipient_numbers": ["120363352351961274@g.us", "120363352351961275@g.us"],
    "scheduled_sending_at": 1735432224,
    "sent_at": 1735432224,
    "retried_count": 0,
    "status": "sent",
    "reason": null,
    "created_at": 1735432224,
    "updated_at": 1735432224
}
```

Scheduled again after a retry (`retried_count` increased):

```json
{
    "id": "3c4d5e6f-7g8h-9i0j-1k2l-m3n4o5p6q7r8",
    "content": "Job alert for Software Engineer at dev.to...",
    "recipient_numbers": ["120363352351961274@g.us", "120363352351961275@g.us"],
    "scheduled_sending_at": 1735432224,
    "sent_at": null,
    "retried_count": 1,
    "status": "scheduled",
    "reason": null,
    "created_at": 1735432224,
    "updated_at": 1735432224
}
```

Failed (`reason` may explain the failure):

```json
{
    "id": "4d5e6f7g-8h9i-0j1k-2l3m-n4o5p6q7r8s9",
    "content": "Job alert for Software Engineer at dev.to...",
    "recipient_numbers": ["120363352351961274@g.us", "120363352351961275@g.us"],
    "scheduled_sending_at": 1735432224,
    "sent_at": null,
    "retried_count": 3,
    "status": "failed",
    "reason": "session expired",
    "created_at": 1735432224,
    "updated_at": 1735432224
}
```

**Success Response** wraps an array of objects like the examples above:

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": [],
    "ts": 1735432224
}
```

`data` is populated with message objects when messages exist.

[Back to Top](#rest-api)

## Schedule Message

POST: `/messages`

Schedules a WhatsApp message for a future time.

**Body Payload:**

| Field                  | Type            | Required | Description                                            |
| ---------------------- | --------------- | -------- | ------------------------------------------------------ |
| `recipient_numbers`    | Array of string | Yes      | Recipient WhatsApp IDs (private JIDs or group `@g.us`). |
| `content`              | String          | Yes      | Message body to send.                                   |
| `scheduled_sending_at` | Number          | Yes      | Unix timestamp (seconds) when the message should send. |

**Example:**

```http
POST /messages HTTP/1.1
Host: localhost:9866
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
Content-Type: application/json

{
    "recipient_numbers": [
        "120363352351961274@g.us",
        "120363352351961275@g.us"
    ],
    "content": "Job alert for Software Engineer at Invertase...",
    "scheduled_sending_at": 1735432224
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "ts": 1735432224
}
```

[Back to Top](#rest-api)

## Retry Message

POST: `/messages/{id}/retry`

Retries a failed message. If `scheduled_sending_at` is omitted from the body, the send is scheduled for **now**.

**Body Payload:**

| Field                  | Type   | Required | Description                                                                                          |
| ---------------------- | ------ | -------- | ---------------------------------------------------------------------------------------------------- |
| `scheduled_sending_at` | Number | No       | Unix timestamp (seconds) for when to send. If omitted, defaults to the current time (immediate send). |

**Example:**

```http
POST /messages/2b3c4d5e-6f7g-8h9i-0j1k-l2m3n4o5p6q7/retry HTTP/1.1
Host: localhost:9866
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
Content-Type: application/json

{
    "scheduled_sending_at": 1735432224
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "ts": 1735432224
}
```

[Back to Top](#rest-api)

## System Errors

Possible error responses from the API:

- Invalid Credentials

  ```json
  HTTP/1.1 401 Unauthorized
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_INVALID_CREDENTIALS",
    "msg": "invalid credentials",
    "ts": 1735432224
  }
  ```

  Authentication credentials were rejected.

- Session Expired

  ```json
  HTTP/1.1 500 Internal Server Error
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_SESSION_EXPIRED",
    "msg": "session expired",
    "ts": 1735432224
  }
  ```

  The WhatsApp session used by the publisher expired. Re-authenticate the publisher; scheduling may halt until then.

- Bad Request

  ```json
  HTTP/1.1 400 Bad Request
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_BAD_REQUEST",
    "msg": "missing `scheduled_sending_at`",
    "ts": 1735432224
  }
  ```

  The request was invalid. Inspect `msg` for detail.

- Internal Server Error

  ```json
  HTTP/1.1 500 Internal Server Error
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_INTERNAL_ERROR",
    "msg": "unable to reach WhatsApp publisher: connection refused",
    "ts": 1735432224
  }
  ```

  An unexpected server-side failure. Inspect `msg` for detail.

  - ID Not found

  ```json
  HTTP/1.1 404 Not Found
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MESSAGE_NOT_FOUND",
    "msg": "message id not found",
    "ts": 1735432224
  }
  ```

  The specified message ID could not be found. Verify that the message exists and that the correct ID was provided.

[Back to Top](#rest-api)
