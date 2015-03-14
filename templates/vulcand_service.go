package templates

import "text/template"

func VulcandService() *Template {
	t, err := template.New("Vulcand v0.7.0 Service").Parse(`
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
ExecStartPost=/bin/sh -c "sleep 10; /usr/bin/etcdctl set /vulcand/upstreams/{{.Name}}/endpoints/{{.Name}}-{{.Version}}-{{.Timestamp}}-%i http://$COREOS_PRIVATE_IPV4:$(echo $(/usr/bin/docker port {{.Name}}-{{.Version}}-{{.Timestamp}}-%i 3000) | cut -d ':' -f 2)"
ExecStop=/bin/sh -c "/usr/bin/etcdctl rm '/vulcand/upstreams/{{.Name}}/endpoints/{{.Name}}-{{.Version}}-{{.Timestamp}}-%i' ; /usr/bin/docker rm -f {{.Name}}-{{.Version}}-{{.Timestamp}}-%i"
`)
	if err != nil {
		panic(err)
	}

	return &Template{Name: "Vulcand v0.7.0 Service", Content: t}
}
