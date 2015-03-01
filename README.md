# Deployster

Deployster is a Golang HTTP service for simplifying deploys to a CoreOS [Fleet cluster][fleet-cluster].  It is extremely opinionated in how you tag your Docker images, where you store them, and how the service's unit files are configured.

This project is also available as `bmorton/deployster` publicly on the [Docker Hub Registry][deployster-docker-hub].

Currently this project is in use for a few side projects, but is not in heavy production use.  At this point it's experimental and feedback is very welcomed and appreciated.  [Yammer][yammer] has been exploring a similar path for production and this will likely be used in some prototyping there, as well.


### Features
* Deploy a new version of a service from a Docker registry (with optionally destroying the previously running version after deploy)
* Shutdown a deployed version of a service
* List all units associated to a service
* Basic authentication and HTTPS support
* Custom task launching for doing things like [migrating a database][running-rails-migrations] using a given service image and version
* Source images from the public Docker Hub Registry or a private registry


### Requirements and limitations

To use Deployster, you'll need:

* **CoreOS cluster** - There are some tutorials for doing this on [DigitalOcean][digitalocean] and [Azure][azure].  Make sure to be using version 550.0.0 or greater of CoreOS so that Fleet's HTTP API is available for Deployster to use.
* **HTTP service exposed on port 3000 of container** - This should be configurable in the future too.
* **Stateless containers** - Linking in volumes is currently not supported.  Again, something for the future.
* **Vulcand running** - For [zero downtime deploys][zero-downtime] of new versions of services.
* **Automatic environment configuration** - As we're currently reusing the same unit file for all services, environment variables can't be passed to containers at boot, so containers need to use something like [etcd], [consul], or [confd][confd] to bootstrap themselves at launch.


#### Task limitations

Tasks are limited to 10 minutes of running time, after which they will be forcefully removed.  If the task is killed, it will have an exit code of 124.  In the future, this timeout may be configurable per task to override the default timeout.


### Getting started

There are two options for getting started with Deployster.  If all the above requirements are fulfilled, you can launch Deployster with Fleet on your own CoreOS cluster and start deploying right away!  If you're just looking for a quick demo to see what Deployster can do, we've provided a Vagrant environment and accompanying tutorial as well.

* [Getting Started with Vagrant][vagrant-guide]
* [Setting up Deployster on your own CoreOS cluster][setup-guide]


### Command line options

```ShellSession
$ deployster -h
Usage of deployster:
  -cert="": Path to certificate to be used for serving HTTPS
  -docker-hub-username="deployster": The username of the Docker Hub account that all deployable images are hosted under
  -key="": Path to private key to be used for serving HTTPS
  -listen="0.0.0.0:3000": Specifies the IP and port that the HTTP server will listen on
  -password="mmmhm": Password that will be used to authenticate with Deployster via HTTP basic auth
  -registry-url="": If using a private registry, this is the address:port of that registry (if supplied, docker-hub-username will be ignored)
  -username="deployster": Username that will be used to authenticate with Deployster via HTTP basic auth
```


### Notes

* For authenticating with the public Docker Hub Registry, follow [this CoreOS guide][registry-authentication].


### Contributing

Pull requests and bug reports are greatly appreciated and encouraged.  If you'd like to help out with any of the above items or have a feature that you think would be awesome for this project, we'd love your help!  To get the design conversation started, open a new GitHub issue with your ideas and we can hash out the details.


### License

Code and documentation copyright 2015 Brian Morton. Code released under the MIT license.

[fleet-cluster]: https://coreos.com/using-coreos/clustering/
[deployster-docker-hub]: https://registry.hub.docker.com/u/bmorton/deployster/
[yammer]: https://www.yammer.com
[running-rails-migrations]: http://guides.rubyonrails.org/active_record_migrations.html#running-migrations
[digitalocean]: https://www.digitalocean.com/community/tutorials/how-to-set-up-a-coreos-cluster-on-digitalocean
[azure]: https://coreos.com/docs/running-coreos/cloud-providers/azure
[zero-downtime]: https://coreos.com/blog/zero-downtime-frontend-deploys-vulcand/
[etcd]: https://github.com/coreos/etcd
[consul]: https://www.consul.io
[confd]: https://github.com/kelseyhightower/confd
[vagrant-guide]: https://github.com/bmorton/deployster/wiki/Getting-Started-with-Vagrant
[setup-guide]: https://github.com/bmorton/deployster/wiki/Setting-up-Deployster-on-your-own-CoreOS-cluster
[registry-authentication]: https://coreos.com/docs/launching-containers/building/registry-authentication/
