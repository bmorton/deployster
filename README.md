# Deployster

Deployster is a Golang HTTP service for simplifying deploys to a CoreOS [Fleet cluster][fleet-cluster].  It is extremely opinionated in how you tag your Docker images, where you store them, and how the service's unit files are configured.

This project is also available as `bmorton/deployster` publically on the [Docker Hub Registry][deployster-docker-hub].

Currently this project is in use for a few side projects, but is not in heavy production use.  At this point it's experimental and feedback is very welcomed and appreciated.  [Yammer][yammer] has been exploring a similar path for production and this will likely be used in some prototyping there, as well.


### Features
* Deploy a new version of a service from the Docker registry (with optionally destroying the previously running version)
* Shutdown a deployed version of a service
* List all units associated to a service
* Basic authentication and HTTPS support
* Custom task launching for doing things like [migrating a database][running-rails-migrations] using a given service image and version


### Requirements and limitations

To use Deployster, you'll need:

* **CoreOS cluster** - There are some tutorials for doing this on [DigitalOcean][digitalocean] and [Azure][azure].  Make sure to be using version 550.0.0 or greater of CoreOS so that Fleet's HTTP API is available for Deployster to use.
* **Images hosted on the public Docker Hub Registry** - They [can be private][registry-authentication], but for now they must come from the same user on the public Docker Hub Registry.  There's a todo item below to make this better.
* **HTTP service exposed on port 3000 of container** - This should be configurable in the future too.
* **Stateless containers** - Linking in volumes is currently not supported.  Again, something for the future.
* **Vulcand running** - For [zero downtime deploys][zero-downtime] of new versions of services.
* **Automatic environment configuration** - As we're currently reusing the same unit file for all services, environment variables can't be passed to containers at boot, so containers need to use something like [etcd], [consul], or [confd][confd] to bootstrap themselves at launch.


#### Task limitations

Tasks are limited to 10 minutes of running time, after which they will be forcefully removed.  If the task is killed, it will have an exit code of 124.  In the future, this timeout may be configurable per task to override the default timeout.


### Getting started

After the above requirements are fulfilled, you can launch Deployster with Fleet.

1. Using a unit file like this one (replacing `password` and `docker-hub-username`), run `fleetctl start deployster.service`

    ```
    [Unit]
    Description=Deployster
    After=docker.service

    [Service]
    TimeoutStartSec=0
    ExecStartPre=/usr/bin/docker pull bmorton/deployster
    ExecStartPre=-/usr/bin/docker rm -f deployster
    # For HTTPS, put your certificate and private key in /home/core/ssl and add:
    # `-v /home/core/ssl:/ssl` to docker options and `-cert=/ssl/server.crt -key=/ssl/server.key` to deployster options
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

3. Run a custom task, like migrating the database, using the same conventions

    ```ShellSession
    $ curl -XPOST http://localhost:3000/v1/services/railsapp/tasks -H "Content-Type: application/json" -d '{"task":{"version":"7bdae1c", "command":"bundle exec rake db:migrate"}}' -u deployster:DONTUSETHIS
    == 20150118005051 CreateUsers: migrating =====================================
    -- create_table(:users)
       -> 0.0017s
    == 20150118005051 CreateUsers: migrated (0.0018s) ============================

    Exited (0)
    ```

4. Deploy an updated version while automatically destroying the previous version once the new one is online

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

5. List units associated to a service

    ```ShellSession
    $ curl http://localhost:3000/v1/services/railsapp/units -u deployster:DONTUSETHIS
    {"units":[{"service":"railsapp","instance":"1","version":"7bdae1c","current_state":"launched","desired_state":"launched","machine_id":"8dcea1bd8c304e1bbe2c25dce526109c"}]}
    ```

6. Manually shutdown a version of a service

    ```ShellSession
    $ curl -XDELETE http://localhost:3000/v1/services/railsapp/deploys/7bdae1c -u deployster:DONTUSETHIS
    ```


### Command line options

```ShellSession
$ deployster -h
Usage of deployster:
  -cert="": Path to certificate to be used for serving HTTPS
  -docker-hub-username="": The username of the Docker Hub account that all deployable images are hosted under
  -key="": Path to private key to bse used for serving HTTPS
  -listen="0.0.0.0:3000": Specifies the IP and port that the HTTP server will listen on
  -password="mmmhm": Password that will be used to authenticate with Deployster via HTTP basic auth
  -username="deployster": Username that will be used to authenticate with Deployster via HTTP basic auth
```


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
[registry-authentication]: https://coreos.com/docs/launching-containers/building/registry-authentication/
[zero-downtime]: https://coreos.com/blog/zero-downtime-frontend-deploys-vulcand/
[etcd]: https://github.com/coreos/etcd
[consul]: https://www.consul.io
[confd]: https://github.com/kelseyhightower/confd
