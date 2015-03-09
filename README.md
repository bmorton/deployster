# Deployster

Deployster uses a convention-over-configuration approach to simplify deploying Docker containers to a CoreOS [Fleet cluster][fleet-cluster] with zero downtime.

It is implemented in Golang as an HTTP service exposing a REST API to automate interactions with Fleet, Docker, and Vulcand.  As part of its convention-based approach, Deployster is opinionated in how images are tagged and stored, the configuration of unit files, and the bootstrapping of containers.


### Comparison with vanilla Fleet

Compared to manually using Fleet to deploy services, Deployster offers conveniences to automate repetitive tasks and provide a layer of conventions to easily handle deployments of many services while making it effortless to deploy new services.

Here are some of the features that make this possible, all exposed through Deployster's REST API:

* Deploy any number of instances of a new service using only the name of the Docker image and the tag/version
* Facilitate new versions by starting up new version units and killing off the old units as the new ones come online
* Launch custom tasks using the same images (for doing things like [migrating a database][running-rails-migrations])
* [`deployctl`](https://github.com/bmorton/deployctl) utility for integrating with CI/CD and command-line workflows
* Authentication and HTTPS support


### Requirements and limitations

To use Deployster, you'll need:

* CoreOS cluster running 550.0.0 or greater (tutorials available for [DigitalOcean][digitalocean] and [Azure][azure])
* HTTP service exposed on port 3000 of image
* [Stateless containers][12-factor-processes]
* Vulcand (for [zero downtime deploys][zero-downtime] while cycling versions)
* Automatic environment configuration.  When your container launches, [etcd] will be available for [bootstrapping your environment][confd].


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

* This project is also available as `bmorton/deployster` publicly on the [Docker Hub Registry][deployster-docker-hub].
* For authenticating with the public Docker Hub Registry, follow [this CoreOS guide][registry-authentication].


### Disclaimer

Currently this project is in use for a few side projects, but is not in heavy production use.  At this point it's experimental and feedback is very welcomed and appreciated.  [Yammer][yammer] has been exploring a similar path for production and this will be used in prototyping.


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
[12-factor-processes]: http://12factor.net/processes
[zero-downtime]: https://coreos.com/blog/zero-downtime-frontend-deploys-vulcand/
[etcd]: https://github.com/coreos/etcd
[consul]: https://www.consul.io
[confd]: https://github.com/kelseyhightower/confd
[vagrant-guide]: https://github.com/bmorton/deployster/wiki/Getting-Started-with-Vagrant
[setup-guide]: https://github.com/bmorton/deployster/wiki/Setting-up-Deployster-on-your-own-CoreOS-cluster
[registry-authentication]: https://coreos.com/docs/launching-containers/building/registry-authentication/
