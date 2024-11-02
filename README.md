# Service Weaver Proof of Concep

In this repository, I am going to become more familliar with [serviceweaver](https://serviceweaver.dev/).

## Setup

The following command did excuted to setup the project:

Installing the weaver command:

```SHELL
go install github.com/ServiceWeaver/weaver/cmd/weaver
```

Installing the weaver command to support google cloud (gke) or kuberneates (kube):

```SHELL
go install github.com/ServiceWeaver/weaver-gke/cmd/weaver-gke
go install github.com/ServiceWeaver/weaver-kube/cmd/weaver-kube
```

```SHELL
go mod init github.com/smhmayboudi/service-weaver-poc
go mod tidy
weaver generate .
go run .
```

run with configuration.

```SHELL
SERVICEWEAVER_CONFIG=weaver.toml go run .
```

check the status

```SHELL
weaver single status
```

show the dashboard

```SHELL
weaver single dashboard
```

for multi process:

```SHELL
weaver multi status
weaver multi dashboard
```

```SHELL
weaver multi deploy weaver.toml
weaver single deploy weaver.toml
```

## Profile

```SHELL
$ profile=$(weaver multi profile <deployment>) # Collect the profile.
$ go tool pprof -http=localhost:9000 $profile # Visualize the profile.
```