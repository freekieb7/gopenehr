# TODO
- Improve validation at unmarshal step
- Improve encoder
- Add XML support (optional)
- Template support
- Audit logging
- Webhooks
- Docs

# Configurable

| Key           | Description |
| ------------- | --------- |
| DATABASE_URL  | REQUIRED: The postgres connection string, for example `postgres://{user}:{pass}@{host}/{name}`    |
| API_KEY       | OPTIONAL: Secret used to protect sensitive endpoints such as /openehr. General endpoints such as `/health` and `OPTIONS /openehr/v1` endpoints remain unprotected. | 
| APP_PORT      | OPTIONAL: Port used for the web server. Default `3000`. |
| LOG_LEVEL     | OPTIONAL: Manages the logs written to stdout. Options are: `DEBUG`, `INFO`, `WARN` or `ERROR`. Default `INFO`. |