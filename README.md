# kube
kube is a kubernetes deployment framework currently under development. There is no guarantee about API stability.

## kube.Coder
package kube defines an interface called `Coder`, which is currently implemented by following providers:
* Persistent volume provider (pv)
* Persistent volume claim provider (pvc)
* Daemon set provider (ds)
* Pods provider (pod)
* Services provider (svc)
* Job provider (job)
* Namespace (ns)
* Sleeper (sleep) // a no op implementation

In other words, a `kube.Coder` instance can be obtained for kubernetes functions listed above.

## Daisy-Chain
The fact that various providers can speak the same language of `kube.Coder` interface, allows us to daisy-chain their
execution using `context.Context`. While explicit dependency management may not be required, passing context between
`kube.Coder` instances allows for executing interesting patters such as:
* Wait for an action completion
* Boot up or shutdown on events

## Examples
Pl. see [this example](https://github.com/sdeoras/kube/blob/master/examples/daisy-chain/main.go) for more details

## Helper Functions
several package level functions are defined to help work with coder objects
### async creation
several coders can deploy their create functions in async manner. Returned
context is `done` when each of the input coders finish their contexts
```go
outCtx, err := Create(inCtx, Async, coders...)
```
### fan in
several context objects can be grouped together to create a new context
```go
outCtx := FanIn(inContexts...)
```

# kubernetes version
this pkg works with `go mod`. the dependencies have been updated for kubernetes version `kubernetes-1.11.6`. to apply
a new version use following commands:
```bash
go get k8s.io/apimachinery@kubernetes-1.11.6
go get k8s.io/client-go@kubernetes-1.11.6
go get k8s.io/api@kubernetes-1.11.6
go get k8s.io/klog@kubernetes-1.11.6
go mod tidy
```