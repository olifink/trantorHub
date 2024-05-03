# Features



## Done - not released yet 

* [x] Allow Anonymous GET (not logged in)
* [x] Read `JwtSecret` from ENV if not configured otherwise
* [x] Deny no-longer existing user in of valid tokens
* [x] Pass identity header to downstream service (pseudonymisation)
* [x] rework handler structures


## In progress

* [ ] re-factor web ui and redirect flow

## Planned - Quick wins

* [ ] User list and lookup as map
* [ ] User management functions (add, remove)
* [ ] Write back users file
* [ ] Token refresh handler
* [ ] CLI way to get token for user
* [ ] Reset password

## Mid-term - ideas before 1.0

* [ ] Role for users and downstream identity
* [ ] Web frontend support
* [ ] Audience from origin and origin validation
* [ ] Custom root redirect config
* [ ] CLI options to add user to users file with bcrypt password hashes
* [ ] Register new user hander (and write back file)
* [ ] Static file serving 
* [ ] Markdown processing for static files (documentation portal)
* [ ] Default embedded login/error forms
* [ ] Auto TLS configuration with domain names list
* [ ] Umami statistics support


## Long-term - future versions

* OpenID Connect federation
* Origin lock configuration
* Multiple route endpoints
* Permissions to allow only some methods  based on groups (GET/HEAD, PUT/POST, DELETE)
* Allow users to manage/generate their own tokens
* SQLite, and other DBs for users/groups (is this really worth it?)
* Catch-all for intrusion detection
* OpenTelemetry throughput, limits
* Re-doc, OpenAPI generation?