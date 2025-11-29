# TODO
- Improve validation at unmarshal step
- Improve encoder
- Add XML support (optional)
- Template support
- Docs
- Add improved rest endpoints

# Configurable

| Key                   | Description                                                                                                                       |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| DATABASE_URL          | REQUIRED: The postgres connection string. Example: `postgres://{user}:{pass}@{host}/{name}`.                                      |
| LOG_LEVEL             | OPTIONAL: Manages the logs written to stdout. Options are: `DEBUG`, `INFO`, `WARN` or `ERROR`. Default `INFO`.                    |
| APP_PORT              | OPTIONAL: Port used for the web server. Default `3000`.                                                                           |
| API_KEY               | OPTIONAL: Enabled API Key protection for endpoints. General endpoints such as `/health` remain unprotected.                       | 
| OAUTH_TRUSTED_ISSUERS | OPTIONAL: Enables OAuth JWT validation with OPENID support. Comma seperated list of trusted token issuers.                        |
| OAUTH_AUDIENCE        | OPTIONAL: Restrict tokens based on `aud` claim. Not provided means no additional `aud` claim check will be performed.             |
| OTEL_ENDPOINT         | OPTIONAL: Enabled OpenTelemetry with GRPC. Example `localhost:4317`.                                                              |
| OTEL_INSECURE         | OPTIONAL: Allows insecure exporter connection. Default `false`.                                                                   | 