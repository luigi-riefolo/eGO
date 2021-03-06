# eGO - Enlightening Golang

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage Report](https://gitlab.com/pantomath-io/demo-tools/badges/master/coverage.svg)](https://gitlab.com/pantomath-io/demo-tools/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/pantomath-io/demo-tools)](https://goreportcard.com/report/gitlab.com/pantomath-io/demo-tools)


eGO is a project skeleton that uses Golang and gRPC.

The intention of this project is to propose a simple approach to the
implementation of a micro-service architecture based on Golang, gRPC, Docker
and Kubernetes.

Furthermore eGO is set up to use Prometheus and Grafana for monitoring.

Users have to only define their proto definitions and implement the respective
functions, all the rest (configuration, middlewares/interceptors, etc.) is
provided by eGO.

_Note:_ this project is still experimental and not ready for production.

### Requirements

[gRPC][1]
[Docker][8]

### Setup

If you are not using the Makefile, then source the environment variables:
```
source .env
```

Create a Kubernetes cluster using Minikube:

```
./scripts/setup_kubernetes.sh
```
or
```
make kubernetes
```

Create a private local registry:
```
./tools/registry/create_registry.sh --docker-user mario -e mario.mario@bros.com
```

Generate the gRPC configuration for a service:

```
./tools/grpc-generator/generate.sh --service omega --project "github.com/luigi-riefolo/eGO"
```
or
```
make config omega
```

Build a service:

```
./scripts/build.sh -s alfa
```
or
```
make run alfa
```

Build and run all the servers:
```
make run_all
```

Test gateway:
```
curl $(minikube service alfa --url)/v1/alfa/get
```

Start the gateway locally:

```
go run src/alfa/cmd/main.go service -config $CONFIG_FILE

BETA_SERVER_PORT=9090 go run src/alfa/cmd/main.go service -config conf/global_conf.toml

go run src/alfa/cmd/main.go service -config $CONFIG_FILE -opts BETA_SERVER_PORT=9090
```

### Credits

Thanks to all the Open Source projects that inspired eGO:

* [gRPC][1]
* [grpc-gateway-generator][3]
* [go-micro-services][4]
* [Golang Project Makefile Template][5]

More to be mentioned.

[1]: http://www.grpc.io/
[2]: http://www.github.com/luigi.riefolo/alfa/contributors
[3]: https://github.com/devsu/grpc-gateway-generator
[4]: https://github.com/harlow/go-micro-services
[5]: https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4
[6]: http://marselester.com/prometheus-on-kubernetes.html
[7]: https://github.com/marselester/prometheus-on-kubernetes
[8]: https://www.docker.com
[9]: https://github.com/giantswarm/kubernetes-prometheus
[10]: https://github.com/olivere/grpc
[11]: https://medium.com/@harlow/tracing-grpc-calls-in-golang-with-google-stackdriver-b22495763a06
[12]: https://github.com/kevinburke/proto-make-example/blob/master/Makefile
[13]: https://about.gitlab.com/2017/11/27/go-tools-and-gitlab-how-to-do-continuous-integration-like-a-boss
