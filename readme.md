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
| DATABASE_URL          | REQUIRED: The postgres connection string, for example `postgres://{user}:{pass}@{host}/{name}`                                    |
| API_KEY               | OPTIONAL: Secret used to protect sensitive endpoints such as /openehr. General endpoints such as `/health` remain unprotected.    | 
| APP_PORT              | OPTIONAL: Port used for the web server. Default `3000`.                                                                           |
| LOG_LEVEL             | OPTIONAL: Manages the logs written to stdout. Options are: `DEBUG`, `INFO`, `WARN` or `ERROR`. Default `INFO`.                    |
| OAUTH_TRUSTED_ISSUERS | OPTIONAL: Comma seperated list of trusted token issuers                                                                           |
| OAUTH_AUDIENCE        | OPTIONAL: Restrict tokens based on `aud` claim. Not provided means no additional `aud` claim check will be performed.             |