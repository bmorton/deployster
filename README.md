# Deployster

Deployster is a Golang HTTP service for simplifying deploys to a CoreOS [Fleet cluster][fleet-cluster].  It is extremely opinionated in how you tag your Docker images, where you store them, and how the service's unit files are configured.

This project is also available as `bmorton/deployster` publically on the [Docker Hub Registry][deployster-docker-hub].

Currently this project is in use for a few side projects, but is not currently in heavy production use.  [Yammer][yammer] has been exploring this path for production and this will likely be used in some prototyping there.


### Features
* Deploy a new version of a service from the Docker registry (with optionally destroying the previously running version)
* Shutdown a deployed version of a service
* List all units associated to a service


### Requirements and limitations

To use Deployster, you'll need:

* **CoreOS cluster** - There are some tutorials for doing this on [DigitalOcean][digitalocean] and [Azure][azure].  Make sure to be using version 550.0.0 or greater of CoreOS so that Fleet's HTTP API is available for Deployster to use.
* **Images hosted on the public Docker Hub Registry** - They [can be private][registry-authentication], but for now they must come from the same user on the public Docker Hub Registry.  There's a todo item below to make this better.
* **Expose an HTTP service on port 3000** - This should be configurable in the future too.
* **Containers are stateless** - Linking in volumes is currently not supported.  Again, something for the future.
* **Vulcand running** - For [zero downtime deploys][zero-downtime] of new versions of services.


### Getting started

After the above requirements are fulfilled, you can launch Deployster with Fleet.

1. Using a unit file like this one, run `fleetctl start deployster.service`

    ```
    [Unit]
    Description=Deployster
    After=docker.service

    [Service]
    EnvironmentFile=/etc/environment
    User=core
    TimeoutStartSec=0
    ExecStartPre=/usr/bin/docker pull bmorton/deployster
    ExecStartPre=-/usr/bin/docker rm -f deployster
    ExecStart=/usr/bin/docker run --name deployster -p 3000:3000 -v /var/run/fleet.sock:/var/run/fleet.sock bmorton/deployster -password=DONTUSETHIS -docker-hub-username=mycompany
    ExecStop=/usr/bin/docker rm -f deployster
    ```

2. Start up a new service that is available on the Docker Hub Registry at mycompany/railsapp:9f88701 (username/service:version)

    ```ShellSession
    $ curl -XPOST http://localhost:3000/v1/services/railsapp/deploys -H "Content-Type: application/json" -d '{"deploy":{"version":"9f88701", "destroy_previous": false}}' -u deployster:DONTUSETHIS
    $ fleetctl list-units
    UNIT                        MACHINE                  ACTIVE  SUB
    railsapp-9f88701@1.service  8dcea1bd.../100.10.11.1  active  running
    ```

3. Deploy an updated version while automatically destroying the previous version once the new one is online

    ```ShellSession
    $ curl -XPOST http://localhost:3000/v1/services/railsapp/deploys -H "Content-Type: application/json" -d '{"deploy":{"version":"7bdae1c", "destroy_previous": true}}' -u deployster:DONTUSETHIS
    $ fleetctl list-units
    UNIT                        MACHINE                  ACTIVE      SUB
    railsapp-7bdae1c@1.service  8dcea1bd.../100.10.11.1  activating  start-pre
    railsapp-9f88701@1.service  8dcea1bd.../100.10.11.1  active      running
    $ fleetctl list-units
    UNIT                        MACHINE                  ACTIVE  SUB
    railsapp-7bdae1c@1.service  8dcea1bd.../100.10.11.1  active  running
    ```

4. List units associated to a service

    ```ShellSession
    $ curl http://localhost:3000/v1/services/carousel/units -u deployster:DONTUSETHIS
    {"units":[{"service":"railsapp","instance":"1","version":"7bdae1c","current_state":"launched","desired_state":"launched","machine_id":"8dcea1bd8c304e1bbe2c25dce526109c"}]}
    ```

5. Manually shutdown a version of a service

    ```ShellSession
    $ curl -XDELETE http://localhost:3000/v1/services/railsapp/deploys/7bdae1c -u deployster:DONTUSETHIS
    ```


### Configurable options

```ShellSession
$ deployster -h
Usage of ./deployster:
  -docker-hub-username="": The username of the Docker Hub account that all deployable images are hosted under
  -listen="0.0.0.0:3000": Specifies the IP and port that the HTTP server will listen on
  -password="mmmhm": Password that will be used to authenticate with Deployster via HTTP basic auth
  -username="deployster": Username that will be used to authenticate with Deployster via HTTP basic auth
```


### Todo

* Move these todos to GitHub issues
* Test coverage (started, but needs more ;/)
* Tutorial/examples for how to set this up, what's required, the limitations, and how your Docker images should be configured
* Documentation
* HTTPS support
* Vagrantfile for easy experimentation and testing
* Allow tasks, such as `rake db:migrate` to be run before a deploy
* Allow multiple instances to be started at once
* Add support for multiple unit templates
* Add support for Docker containers that need volumes linked (not stateless)


### Contributing

Pull requests and bug reports are greatly appreciated and encouraged.  If you'd like to help out with any of the above items or have a feature that you think would be awesome for this project, we'd love your help!  To get the design conversation started, open a new GitHub issue with your ideas and we can hash out the details.


### License

Code and documentation copyright 2015 Brian Morton. Code released under the MIT license.

[fleet-cluster]: https://coreos.com/using-coreos/clustering/
[deployster-docker-hub]: https://registry.hub.docker.com/u/bmorton/deployster/
[yammer]: https://www.yammer.com
[digitalocean]: https://www.digitalocean.com/community/tutorials/how-to-set-up-a-coreos-cluster-on-digitalocean
[azure]: https://coreos.com/docs/running-coreos/cloud-providers/azure
[registry-authentication]: https://coreos.com/docs/launching-containers/building/registry-authentication/
[zero-downtime]: https://coreos.com/blog/zero-downtime-frontend-deploys-vulcand/
