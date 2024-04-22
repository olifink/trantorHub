# Features

## Done 

* Allow Anonymous GET (not logged in)

# Plans

## Quick wins

* Token refresh handler
* Origin lock configuration
* Cli get token for user

## New feature ideas before 1.0

* Custom root redirect config
* Write back users file
* CLI options to add user to users file with bcrypt password hashes
* Register new user hander (and write back file)
* Static file serving 
* Markdown processing for static files (documentation portal)
* Default embedded login/error forms
* Auto TLS configuration with domain names list
* Umami statistics support

## Bigger ideas

* Multiple route endpoints
* Permissions to allow only some methods  based on groups (GET/HEAD, PUT/POST, DELETE)
* SQLite for users/groups (is this really worth it?)
* OpenTelemetry