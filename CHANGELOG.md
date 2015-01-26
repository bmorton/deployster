## v0.2.0 (unreleased)

Features:

  * HTTPS support
  * Custom task launching for doing things like migrating a database using a given service image and version ([#14][issue-14])

Fixes:

  * Improve code to be more idiomatic
  * More documentation and test coverage


## v0.1.0

Features:

  * Deploy a new version of a service from the Docker registry (with optionally destroying the previously running version)
  * Shutdown a deployed version of a service
  * List all units associated to a service
  * HTTP basic authentication for all endpoints

[issue-14]: https://github.com/bmorton/deployster/pull/14
