## v0.2.0

Features:

  * HTTPS support
  * Custom task launching for doing things like migrating a database using a given service image and version ([#14][issue-14])
  * Private registry support ([#20][issue-20])
  * Vagrant environment ([#17][issue-17])

Fixes:

  * Improve code to be more idiomatic
  * More documentation and test coverage
  * Allow service names to contain hyphens ([#27][issue-27])


## v0.1.0

Features:

  * Deploy a new version of a service from the Docker registry (with optionally destroying the previously running version)
  * Shutdown a deployed version of a service
  * List all units associated to a service
  * HTTP basic authentication for all endpoints

[issue-14]: https://github.com/bmorton/deployster/pull/14
[issue-20]: https://github.com/bmorton/deployster/pull/20
[issue-17]: https://github.com/bmorton/deployster/pull/17
[issue-27]: https://github.com/bmorton/deployster/pull/27
