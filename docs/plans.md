# Features



## Done - not released yet 

* [x] Allow Anonymous GET (not logged in)

## In progress

## Planned - Quick wins

* [ ] Token refresh handler
* [ ] Check for non-existant user in token validation
* [ ] Origin lock configuration
* [ ] Cli get token for user
* [ ] Pass anonymized user to backend service

## Mid-term - ideas before 1.0

* Custom root redirect config
* Write back users file
* CLI options to add user to users file with bcrypt password hashes
* Register new user hander (and write back file)
* Static file serving 
* Markdown processing for static files (documentation portal)
* Default embedded login/error forms
* Auto TLS configuration with domain names list
* Umami statistics support

## Long-term - future versions

* Multiple route endpoints
* Permissions to allow only some methods  based on groups (GET/HEAD, PUT/POST, DELETE)
* Allow users to manage/generate their own tokens
* SQLite, and other DBs for users/groups (is this really worth it?)
* OpenTelemetry throughput, limits
* Re-doc, OpenAPI generation?