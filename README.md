# Deployster

Deployster is a Golang HTTP service for simplifying deploys to a fleet cluster.  It is extremely opinionated in how you tag your Docker images, where you store them, and how the service's unit files are configured.


### Features
* Deploy a new version of a service from the Docker registry
```ShellSession
$ curl -XPOST http://localhost:3000/v1/services/carousel/deploys -H "Content-Type: application/json" -d '{"deploy":{"version":"9f88701", "destroy_previous": true}}'
$ fleetctl list-units
UNIT        MACHINE       ACTIVE  SUB
carousel-9f88701@1.service  8dcea1bd.../100.78.68.84  active  running
```

* Shutdown a deployed version of a service
```ShellSession
$ curl -XDELETE http://localhost:3000/v1/services/carousel/deploys/9f88701
```

* List all units associated to a service

```ShellSession
$ curl http://localhost:3000/v1/services/carousel/units
{"units":[{"service":"carousel","instance":"1","version":"9f88701","current_state":"launched","desired_state":"launched","machine_id":"8dcea1bd8c304e1bbe2c25dce526109c"}]}
```


### Configurable options

```ShellSession
$ ./deployster -h
Usage of ./deployster:
  -docker-hub-username="": The username of the Docker Hub account that all deployable images are hosted under
  -listen="0.0.0.0:3000": Specifies the IP and port that the HTTP server will listen on
  -password="mmmhm": Password that will be used to authenticate with Deployster via HTTP basic auth
  -username="deployster": Username that will be used to authenticate with Deployster via HTTP basic auth
```


### Example unit for starting Deployster

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
ExecStart=/usr/bin/docker run --name deployster -p 3000 bmorton/deployster
ExecStartPost=/bin/sh -c "sleep 5; /usr/bin/etcdctl set /vulcand/upstreams/deployster/endpoints/deployster http://$COREOS_PRIVATE_IPV4:$(echo $(/usr/bin/docker port deployster 3000) | cut -d ':' -f 2)"
ExecStop=/bin/sh -c "/usr/bin/etcdctl rm '/vulcand/upstreams/deployster/endpoints/deployster' ; /usr/bin/docker rm -f deployster"
```


### Todo

* Test coverage
* Tutorial/examples for how to set this up, what's required, the limitations, and how your Docker images should be configured
* HTTPS support
* Vagrantfile for easy experimentation and testing
* Allow tasks, such as `rake db:migrate` to be run before a deploy
* Allow multiple instances to be started at once
* Add support for multiple unit templates
* Add support for Docker containers that need volumes linked (not stateless)
