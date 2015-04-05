## v0.4.0 - unreleased

Fixes:

  * Fixed incorrect default instance count in some cases ([#50][issue-50])
  * Refactored polling/previous version destruction to be more efficient and clean up the code ([#52][issue-52])


## v0.3.0 - March 8th, 2015

Features:

  * Deployment of multiple instances ([#30][issue-30])
  * Initial release of [`deployctl`](https://github.com/bmorton/deployctl)

Fixes:

  * Removed homegrown Fleet client and replaced with official CoreOS client ([#29][issue-29])
  * Fixed panic that occurred when Fleet units not managed by Deployster were running ([#39][issue-39])
  * Fixed bug with ignoring Fleet responses with multiple pages


## v0.2.0 - February 18th, 2015

Features:

  * HTTPS support
  * Custom task launching for doing things like migrating a database using a given service image and version ([#14][issue-14])
  * Private registry support ([#20][issue-20])
  * Vagrant environment ([#17][issue-17])

Fixes:

  * Improve code to be more idiomatic
  * More documentation and test coverage
  * Allow service names to contain hyphens ([#27][issue-27])


## v0.1.0 - January 5th, 2015

Features:

  * Deploy a new version of a service from the Docker registry (with optionally destroying the previously running version)
  * Shutdown a deployed version of a service
  * List all units associated to a service
  * HTTP basic authentication for all endpoints


[issue-14]: https://github.com/bmorton/deployster/pull/14
[issue-20]: https://github.com/bmorton/deployster/pull/20
[issue-17]: https://github.com/bmorton/deployster/pull/17
[issue-27]: https://github.com/bmorton/deployster/pull/27
[issue-29]: https://github.com/bmorton/deployster/pull/29
[issue-30]: https://github.com/bmorton/deployster/pull/30
[issue-39]: https://github.com/bmorton/deployster/pull/39
[issue-50]: https://github.com/bmorton/deployster/pull/50
[issue-52]: https://github.com/bmorton/deployster/pull/52
