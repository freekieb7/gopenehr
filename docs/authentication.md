````markdown
# Authentication using OAuth 2.0 PKCE with CareBase24 Smart Auth

This document explains how to implement the **OAuth 2.0 PKCE flow** using the CareBase24 **Smart Auth** service (`https://smart-auth.dev.carebase24.io`).

---

## 1. Prerequisites

- API secret to access Smart Auth for client registration (provided by Code24).
- API reference documentation for your environment.
- Access to a development or production Smart Auth environment.

### 1.1 API Secret

The API secret is used for **system-to-system** operations like client registration or launch token creation — **not** for end-user login.

**Authentication header:**
```http
Authorization: Bearer <API secret>
````

**Base URL (development):**

```
https://smart-auth.dev.carebase24.io
```

> ⚠️ The API secret may change over time. Update it in your environment or Postman variables when needed.

---

## 2. Preparing the Authentication Flow

### 2.1 Launch URI

Authentication is initiated through a **launch URI**, which includes the following parameters:

* `iss` — base URL of Smart Auth (e.g. `https://smart-auth.dev.carebase24.io`)
* `client_id` — registered OAuth client ID
* `launch` — one-time launch token

**Example launch URI:**

```
https://app.dev.carebase24.io/launch?iss=https%3A%2F%2Fsmart-auth.dev.carebase24.io&client_id=fb9ad137-36e3-4531-8166-428382a6d084&launch=8390d844-2075-4be4-ab15-4efefbf20ef2
```

Your application should parse `iss`, `client_id`, and `launch` to initiate the PKCE flow.

---

### 2.2 Development Usage

For development, you can manually create launch tokens.

#### 2.2.1 Client Registration

```http
POST https://smart-auth.dev.carebase24.io/api/clients
Authorization: Bearer <API secret>
Content-Type: application/json

{
    "name": "My App",
    "description": "Development client",
    "logo_uri": "https://picsum.photos/200",
    "is_confidential": true,
    "redirect_uris": [
        "https://app.dev.carebase24.io/redirect"
    ],
    "allowed_origins": ["*"],
    "meta_data": {
        "smart": true
    }
}
```

The response will include:

```json
{
    "id": "fb9ad137-36e3-4531-8166-428382a6d084",
    "client_secret": "...",
    "redirect_uris": [...]
}
```

#### 2.2.2 Obtain a Launch Token

```http
POST https://smart-auth.dev.carebase24.io/api/smart/launch
Authorization: Bearer <API secret>
Content-Type: application/x-www-form-urlencoded

client_id=<client_id>
flow_type=practitioner
party_id=<practitioner_uid>
```

**Response:**

```json
{
    "launch": "8390d844-2075-4be4-ab15-4efefbf20ef2",
    "expires_in": 600
}
```

#### 2.2.3 Construct Launch URI (Optional)

```
https://app.dev.carebase24.io/launch?iss=https%3A%2F%2Fsmart-auth.dev.carebase24.io&client_id=<client_id>&launch=<launch_token>
```

---

## 3. Initiate Authentication Flow

### 3.1 Configuration

* `launch_uri` — base of the launch URI
* `redirect_uri` — must match the registered client’s redirect URI

### 3.2 Discover OAuth Endpoints

```http
GET {iss}/.well-known/smart-configuration
```

**Example response:**

```json
{
  "authorization_endpoint": "https://smart-auth.dev.carebase24.io/oauth/authorize",
  "token_endpoint": "https://smart-auth.dev.carebase24.io/oauth/token"
}
```

You can also use:

```http
GET {iss}/.well-known/openid-configuration
```

---

## 3.3 OAuth2 PKCE Flow

### 3.3.1 Generate PKCE Parameters

* `code_verifier`: random string (43–128 chars)
* `code_challenge`: `Base64URL(SHA256(code_verifier))`

**Example:**

```
code_verifier  = dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk
code_challenge = E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM
```

---

### 3.3.2 Authorize Request

```http
GET https://smart-auth.dev.carebase24.io/oauth/authorize
```

| Parameter               | Description                            |
| ----------------------- | -------------------------------------- |
| `response_type`         | `code`                                 |
| `client_id`             | From launch URI                        |
| `scope`                 | `offline_access openid profile launch` |
| `redirect_uri`          | As registered                          |
| `code_challenge`        | Generated above                        |
| `code_challenge_method` | `S256`                                 |
| `audience`              | As required by resource server         |
| `launch`                | From launch URI                        |

**Redirect response example:**

```
https://app.dev.carebase24.io/redirect?code=authorization_code
```

---

### 3.3.3 Token Exchange

```http
POST https://smart-auth.dev.carebase24.io/oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code
code=<authorization_code>
code_verifier=<code_verifier>
client_id=<client_id>
redirect_uri=<redirect_uri>
```

**Token response:**

```json
{
    "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImFkMmY1M2FjNTAxNTFmNTRjOTJmNDQ3YzFlYzE5OTU1MTc1ZGRiOWMifQ...",
    "refresh_token": "f0dc3130fd4aeb4d65406811fabff26a7dbc3e3e",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

Use the `access_token` to authenticate subsequent API calls.

---

## 4. Example API Request

```http
GET https://api.dev.carebase24.io/v1/user
Authorization: Bearer <access_token>
```

---

## 5. Test Users

For development, the following test users are available:

| Name              | UID                                    | Purpose                       |
| ----------------- | -------------------------------------- | ----------------------------- |
| Code24 CareBase24 | `eec91c3f-363a-4fe0-a4bd-9bc8b022879b` | Full access                   |
| John              | `1c6d654b-430f-34fe-8029-acf2ea1284e5` | Limited patient range         |
| Jane              | `cab51da4-e206-3209-bad7-8c55495a4e02` | Group-based scenario          |
| Jimmy             | `37840551-4ba3-337d-b732-2480ae4a311b` | No group membership           |
| Jannette          | `85796d1b-488e-3393-8584-7edbd97b8097` | Disabled/unaffiliated account |

---

## 6. Error Format

Errors are returned in a standardized JSON structure:

```json
{
    "title": "400 Bad Request",
    "description": "The server cannot or will not process the request due to an apparent client error.",
    "message": "Missing \"party_id\" or \"client_id\" parameter."
}
```

---

## ✅ Summary

* Use `https://smart-auth.dev.carebase24.io` for the OAuth 2.0 PKCE flow.
* Register your OAuth client once using the API secret.
* Generate a launch token to start the flow.
* Use the returned `access_token` to access CareBase24 APIs.
* Follow SMART on FHIR/OIDC best practices.

---

## 🔗 Useful Endpoints

| Purpose                | URL                                     |
| ---------------------- | --------------------------------------- |
| Smart Configuration    | `GET /.well-known/smart-configuration`  |
| OpenID Configuration   | `GET /.well-known/openid-configuration` |
| Authorization Endpoint | `GET /oauth/authorize`                  |
| Token Endpoint         | `POST /oauth/token`                     |
| Client Registration    | `POST /api/clients`                     |
| Launch Token Creation  | `POST /api/smart/launch`                |

---

© CareBase24 – Authentication Guide

```

---

Would you like me to also make a **Postman collection** for this flow (with client registration, launch, authorize, and token requests preconfigured)?
```
