package main

const DOCKER_UNIT_TEMPLATE = `
[Unit]
Description={{.Name}}
After=docker.service

[Service]
EnvironmentFile=/etc/environment
User=core
TimeoutStartSec=0
ExecStartPre=/usr/bin/docker pull mmmhm/{{.Name}}:{{.Version}}
ExecStartPre=-/usr/bin/docker rm -f {{.Name}}-{{.Version}}-%i
ExecStart=/usr/bin/docker run --name {{.Name}}-{{.Version}}-%i -p 3000 mmmhm/{{.Name}}:{{.Version}}
ExecStartPost=/bin/sh -c "sleep 15; /usr/bin/etcdctl set /vulcand/upstreams/{{.Name}}/endpoints/{{.Name}}-{{.Version}}-%i http://$COREOS_PRIVATE_IPV4:$(echo $(/usr/bin/docker port {{.Name}}-{{.Version}}-%i 3000) | cut -d ':' -f 2)"
ExecStop=/bin/sh -c "/usr/bin/etcdctl rm '/vulcand/upstreams/{{.Name}}/endpoints/{{.Name}}-{{.Version}}-%i' ; /usr/bin/docker rm -f {{.Name}}-{{.Version}}-%i"
`
