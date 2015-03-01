# Deployster API v1

The Deployster API is the main entrypoint into controlling Deployster and executing deploys and tasks.

## Authentication

## Deploys resource

### Start a new deploy
Trigger an asynchronous deploy of a service, optionally cleaning up an old version upon completion.

`POST /v1/services/{name}/deploys`

```json
{
  "deploy": {
    "version": "abc123f",
    "destroy_previous": true, // (optional)
    "timestamp": "2006.01.02-15.04.05", // (optional)
    "instance_count": 4 // (optional)
  }
}
```

#### Deploy entity
  * `version` (string): the tagged version of the Docker container to deploy (required)
  * `destroy_previous` (boolean): clean up previous version after the new version has been deployed (optional, default `false`)
  * `timestamp` (string): a date formatted as `2006.01.02-15.04.05` to include with all instances of the deployment (optional, default `time.Now()`)
  * `instance_count` (integer): the number of instances of the deployment to be launched (optional, default is 0 which tells Deployster to use the number currently running of the previous version *or* 1 if unable to determine)

#### Response
A `201 Created` with an empty body will be returned when a deploy is successfully triggered.

##### Errors
  * `400 Bad Request`
    * Too many versions are running.  Destroying previous units is not supported when more than one version is currently running.
    * A greater number of instances than what was specified is already running.  Make sure this number is less than or equal to the number already running or disable destroying previous units.
  * `500 Internal Server Error` - any failure communicating with Fleet


### Shutdown a deployed service/version
Destroy all containers associated to a service's version, optionally locked to a specific timestamp.

`DELETE /v1/services/{name}/deploys/{version}`

#### Query parameters
  * `timestamp` (string): a date formatted as `2006.01.02-15.04.05` if only a certain deploy of a version should be shutdown (optional, default is all timestamps related to the given version)


## Tasks resource

### Launch a new task
Run a time-boxed task using an image of the given service and version.  Tasks must complete within 10 minutes or they will be forcefully stopped.

`POST /v1/services/{name}/tasks`

```json
{
  "task": {
    "version": "abc123f",
    "command": "bundle exec rake db:migrate"
  }
}
```

#### Task entity
  * `version` (string): the tagged version of the Docker container to use for running the task (required)
  * `command` (string): the command to launch the Docker container with (required)

#### Response
A `200 OK` with `text/plain` output of the running container will be streamed back via the response until the container exists.  The last line of output will be the exit code of the container (e.g. `Exited (0)`).

##### Errors
If an error occurs decoding the JSON or creating/running the container, a `500 Internal Server Error` will be returned in the response.  However, if an error occurs after this point, we've already sent a `200 OK` and started streaming the response body.  This means the task was successfully launched, but the task could have possibly errored out.  At the end of the task output, the exit code of the task will be printed so that it can be handled by the client if necessary.


## Units resource

### Retrieve service's units
Retrieve the status of all units associated to a given service.

`GET /v1/services/{name}/units`

#### Unit entity
  * `service` (string): name of the service
  * `instance` (string): the number of the instance
  * `version` (string): the tagged version of the Docker container
  * `current_state` (string): the current state of the systemd/Fleet unit
  * `desired_state` (string): the state that the systemd/Fleet unit is supposed to be
  * `machine_id` (string): the Fleet machine ID where the instance is running
  * `timestamp` (string): a date formatted as `2006.01.02-15.04.05` signifying when the unit was deployed

#### Response
A `200 OK` with an `application/json` output including an array of units.

##### Errors
A `500 Internal Server Error` will be returned for any failure communicating with Fleet.
