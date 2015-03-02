# Deployster API v1

The Deployster API is the main entrypoint into controlling Deployster and executing deploys and tasks.

## Authentication
Requests are authenticated using HTTP Basic authentication.  Deployster is launched while specifying a username and password that must be provided with all requests.  If credentials are not provided, a `401 Unauthorized` will be returned.

## Deploys resource

### Start a new deploy
Trigger an asynchronous deploy of a service, optionally cleaning up an old version upon completion.

```http
POST /v1/services/{name}/deploys HTTP/1.1
Authorization: Basic dGVzdDp0ZXN0
Content-Type: application/json

{
  "deploy": {
    "version": "abc123f",
    "destroy_previous": true,
    "timestamp": "2006.01.02-15.04.05",
    "instance_count": 4
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

```http
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 02 Mar 2015 00:21:42 GMT
Content-Length: 0
```

##### Errors
  * `400 Bad Request`
    * Too many versions are running.  Destroying previous units is not supported when more than one version is currently running.
    * A greater number of instances than what was specified is already running.  Make sure this number is less than or equal to the number already running or disable destroying previous units.
  * `500 Internal Server Error` - any failure communicating with Fleet


### Shutdown a deployed service/version
Destroy all containers associated to a service's version, optionally locked to a specific timestamp.

```http
DELETE /v1/services/{name}/deploys/{version} HTTP/1.1
Authorization: Basic dGVzdDp0ZXN0
```

#### Query parameters
  * `timestamp` (string): a date formatted as `2006.01.02-15.04.05` if only a certain deploy of a version should be shutdown (optional, default is all timestamps related to the given version)

#### Response
A `204 No Content` will be returned if the destroy is successfully triggered.

```http
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Mon, 02 Mar 2015 00:24:41 GMT
```


## Tasks resource

### Launch a new task
Run a time-boxed task using an image of the given service and version.  Tasks must complete within 10 minutes or they will be forcefully stopped.  If forcefully stopped, an exit code of 127 will be provided at the end of the response.

```http
POST /v1/services/{name}/tasks HTTP/1.1
Authorization: Basic dGVzdDp0ZXN0
Content-Type: application/json

{
  "task": {
    "version": "abc123f",
    "command": "bundle check"
  }
}
```

#### Task entity
  * `version` (string): the tagged version of the Docker container to use for running the task (required)
  * `command` (string): the command to launch the Docker container with (required)

#### Response
A `200 OK` with `text/plain` output of the running container will be streamed back via the response until the container exists.  The last line of output will be the exit code of the container (e.g. `Exited (0)`).

```http
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Mon, 02 Mar 2015 00:30:17 GMT
Transfer-Encoding: chunked

The Gemfile's dependencies are satisfied

Exited (0)
```

##### Errors
If an error occurs decoding the JSON or creating/running the container, a `500 Internal Server Error` will be returned in the response.  However, if an error occurs after this point, we've already sent a `200 OK` and started streaming the response body.  This means the task was successfully launched, but the task could have possibly errored out.  At the end of the task output, the exit code of the task will be printed so that it can be handled by the client if necessary.


## Units resource

### Retrieve service's units
Retrieve the status of all units associated to a given service.

```http
GET /v1/services/{name}/units HTTP/1.1
Authorization: Basic dGVzdDp0ZXN0
```

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

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 02 Mar 2015 00:32:42 GMT
Content-Length: 216

{"units":[{"service":"hello-world","instance":"1","version":"0fbb804","current_state":"launched","desired_state":"launched","machine_id":"d6e4b05a215d4ac2839da17017ed1d59","deploy_timestamp":"2015.03.02-00.31.45"}]}
```

##### Errors
A `500 Internal Server Error` will be returned for any failure communicating with Fleet.


## cURL Examples

* `POST /v1/services/hello-world/deploys`
```ShellOutput
$ curl -v -u deployster:mmmhm -XPOST http://localhost:1234/v1/services/hello-world/deploys -H "Content-Type: application/json" -d '{"deploy":{"version":"0fbb804", "destroy_previous": true}}'
* About to connect() to localhost port 1234 (#0)
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 1234 (#0)
* Server auth using Basic with user 'deployster'
> POST /v1/services/hello-world/deploys HTTP/1.1
> Authorization: Basic ZGVwbG95c3RlcjptbW1obQ==
> User-Agent: curl/7.37.1
> Host: localhost:1234
> Accept: */*
> Content-Type: application/json
> Content-Length: 58
>
* upload completely sent off: 58 out of 58 bytes
< HTTP/1.1 201 Created
< Content-Type: application/json
< Date: Mon, 02 Mar 2015 00:21:42 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

* `DELETE /v1/services/hello-world/deploys/0fbb804`
```ShellOutput
$ curl -vvv -u deployster:mmmhm -XDELETE http://localhost:1234/v1/services/hello-world/deploys/0fbb804
* About to connect() to localhost port 1234 (#0)
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 1234 (#0)
* Server auth using Basic with user 'deployster'
> DELETE /v1/services/hello-world/deploys/0fbb804 HTTP/1.1
> Authorization: Basic ZGVwbG95c3RlcjptbW1obQ==
> User-Agent: curl/7.37.1
> Host: localhost:1234
> Accept: */*
>
< HTTP/1.1 204 No Content
< Content-Type: application/json
< Date: Mon, 02 Mar 2015 00:24:41 GMT
<
* Connection #0 to host localhost left intact
```

* `POST /v1/services/hello-world/tasks`
```ShellOutput
$ curl -vvv -u deployster:mmmhm -XPOST http://localhost:1234/v1/services/hello-world/tasks -H "Content-Type: application/json" -d '{"task":{"version":"0fbb804", "command": "bundle check"}}'
* About to connect() to localhost port 1234 (#0)
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 1234 (#0)
* Server auth using Basic with user 'deployster'
> POST /v1/services/hello-world/tasks HTTP/1.1
> Authorization: Basic ZGVwbG95c3RlcjptbW1obQ==
> User-Agent: curl/7.37.1
> Host: localhost:1234
> Accept: */*
> Content-Type: application/json
> Content-Length: 57
>
* upload completely sent off: 57 out of 57 bytes
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Date: Mon, 02 Mar 2015 00:30:17 GMT
< Transfer-Encoding: chunked
<
stdin: is not a tty
stdin: is not a tty
The Gemfile's dependencies are satisfied

Exited (0)
* Connection #0 to host localhost left intact
```

```ShellOutput
$ curl -v -u deployster:mmmhm http://localhost:1234/v1/services/hello-world/units
* About to connect() to localhost port 1234 (#0)
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 1234 (#0)
* Server auth using Basic with user 'deployster'
> GET /v1/services/hello-world/units HTTP/1.1
> Authorization: Basic ZGVwbG95c3RlcjptbW1obQ==
> User-Agent: curl/7.37.1
> Host: localhost:1234
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Mon, 02 Mar 2015 00:33:48 GMT
< Content-Length: 216
<
{"units":[{"service":"hello-world","instance":"1","version":"0fbb804","current_state":"launched","desired_state":"launched","machine_id":"d6e4b05a215d4ac2839da17017ed1d59","deploy_timestamp":"2015.03.02-00.31.45"}]}
* Connection #0 to host localhost left intact
```
