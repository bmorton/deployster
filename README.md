# Deployster

Deployster is a Golang HTTP service for simplifying deploys to a fleet cluster.  It is extremely opinionated in how you tag your Docker images, where you store them, and how the service's unit files are configured.


### Features
* Deploy a new version of a service from the Docker registry
```ShellSession
$ curl -XPOST http://localhost:3000/v1/services/carousel/deploys -H "Content-Type: application/json" -d '{"deploy":{"version":"9f88701"}}'
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


### Todo

* Remove hardcoding of `mmmhm` user in docker template
* Authentication
* Shut down previous version after deploy is complete
* Allow tasks, such as `rake db:migrate` to be run before a deploy
* Allow multiple instances to be started at once
* Add support for multiple unit templates
