# Kube Controller Runtime Example

Basic example of using the kube controller-runtime package. Example code stolen
pretty much verbatim from [the package's examples][pkg-example].

Broken into a separate repository so I could split component parts into separate
commits to share with colleagues.

## Running

* Have a valid ~/.kube/config file pointing at the context you want to use
* `kubectl apply -f ./deploy.yaml`
* `go run main.go`
* `kubectl scale deploy/hello --replicas=2`
* `kubectl get rs -Lpod-count` to see label added by controller.

[pkg-example]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/builder#example-Builder
