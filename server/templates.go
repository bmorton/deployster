package server

// dockerUnitTemplate is the only currently supported Fleet unit file for
// launching new units.  It makes lots of assumptions about how the service is
// configured and stored.  These assumptions are essentially the conventions
// that power deployster and are described in more detail in the README.
//
// Additionally, we only store this unit template to make it easy to read and
// update.  We always convert this unit file to an array of fleet.UnitOption
// structs before sending it off to the Fleet client.
const dockerUnitTemplate = `
[Unit]
Description={{.Name}}-{{.Version}}-{{.Timestamp}}
After=docker.service

[Service]
EnvironmentFile=/etc/environment
User=core
TimeoutStartSec=0
ExecStartPre=/usr/bin/docker pull {{.ImagePrefix}}/{{.Name}}:{{.Version}}
ExecStartPre=-/usr/bin/docker rm -f {{.Name}}-{{.Version}}-{{.Timestamp}}-%i
ExecStart=/usr/bin/docker run --name {{.Name}}-{{.Version}}-{{.Timestamp}}-%i -p 3000 {{.ImagePrefix}}/{{.Name}}:{{.Version}}
ExecStartPost=/bin/sh -c "sleep 3; /usr/bin/etcdctl set /vulcand/upstreams/{{.Name}}/endpoints/{{.Name}}-{{.Version}}-{{.Timestamp}}-%i http://$COREOS_PRIVATE_IPV4:$(echo $(/usr/bin/docker port {{.Name}}-{{.Version}}-{{.Timestamp}}-%i 3000) | cut -d ':' -f 2)"
ExecStop=/bin/sh -c "/usr/bin/etcdctl rm '/vulcand/upstreams/{{.Name}}/endpoints/{{.Name}}-{{.Version}}-{{.Timestamp}}-%i' ; /usr/bin/docker rm -f {{.Name}}-{{.Version}}-{{.Timestamp}}-%i"
`
