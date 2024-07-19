<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="/docs/logo.png">
  <img alt="tiny corp logo" src="/docs/logo.png" width="50%" height="50%">
</picture>

tiny-is : A fun project where I'm building an Identity and Access Management (IAM) product from scratch.

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
- Run the server
```bash
make run
```

### User Management:
- Add users
- Basic user authentication

### Application Management:
- Basic application management (client_id, client_secret, redirect_uris, grant_types)

### OAuth2.0
- Authorization Code Grant (with PKCE)
- Refresh Token Grant
- Client Credentials Grant

### Token Management
- jwt access and refresh tokens (EdDSA)
- token revocation

## Session
- in-memory session storage


