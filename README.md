<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="/docs/logo.png">
  <img alt="tiny corp logo" src="/docs/logo.png" width="50%" height="50%">
</picture>

tiny-is : A fun project where I'm building a light-weight framework that provides implementations of the [OAuth 2.1](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-v2-1-10) and [OpenID Connect 1.0](https://openid.net/specs/openid-connect-core-1_0.html) specifications and other related specifications. The framework follows the [OAuth 2.0 Security Best Current Practice](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)

</div>

### Stack:
- Golang
- SQLite
- HTMX

### Run Locally:

- Create sqlite database
```bash
make create_db
```
- Generate Keys
```bash
make generate_jwt_key
make generate_server_keypair
```

- Run the server
```bash
make run
```

### OAuth2.1
- Authorization Code Grant with PKCE
  - Only support with PKCE for better security
- Refresh Token Grant
- Client Credentials Grant
- Authorization Server Metadata

### Token Management
- JWT access and refresh tokens (EdDSA)
- Token revocation

### User Management:
- Add users
- Basic user authentication

### Application Management:
- Basic application management (client_id, client_secret, redirect_uris, grant_types)

## Session
- in-memory session storage


