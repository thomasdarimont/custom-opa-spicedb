Custom Open Policy Agent with prototypical support for Authzed
---

This experiment adds support for querying relations from [Authzed](https://authzed.com/) / [SpiceDB](https://github.com/authzed/spicedb) via GRPC to check resource level permissions
as custom builtin commands for [Open Policy Agent](https://www.openpolicyagent.org/).

Currently only one command is supported:
```
authzed.check_permission("SUBJECT", "PERMISSION", "RESOURCE_ID") -> bool
```

# Build

Note this example uses Go 1.19

```
go get
go build
```

# Demo

> Start authzed demo environment
```
docker compose -f demo/docker-compose.yml up -d
```

> Run custom Open Policy Agent with authzed plugin enabled
```
./custom-opa-spicedb run \
  --set plugins.authzed.endpoint=localhost:50051 \
  --set plugins.authzed.token=foobar \
  --set plugins.authzed.insecure=true
```

> Query relations against authzed
> See the [example RBAC schema](./demo/schema-and-data.yml) for reference.
```
> authzed.check_permission("user:tom", "view", "document:firstdoc")
true
> authzed.check_permission("user:tom", "edit", "document:firstdoc")
true
> authzed.check_permission("user:fred", "edit", "document:firstdoc")
false
> exit
```

> Stop demo environment
```
docker compose -f demo/docker-compose.yml down
```
