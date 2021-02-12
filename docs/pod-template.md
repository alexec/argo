# Pod Template

![alpha](assets/alpha.svg)

> v3.1 and after

A pod templates is similar to a normal container or script template, but allows you to specify multiple containers to
run within the pod.

Because you have multiple containers within a pod, they will be scheduled on the same host. You can use cheap and fast
empty-dir volumes instead of persistent volume claims to share data between steps.

For example, if you have nested DAGs, and the nested DAGs would benefit from a shared workspace.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: containers-
spec:
  entrypoint: main
  templates:
    - name: main
      volumes:
        - name: workpsace
          emptyDir: { }
      pod:
        volumeMounts:
          - mountPath: /workspace
            name: workspace
        graph:
          - name: a
            image: argoproj/argosay:v2
            args: [ echo, hello, /mnt/share/message ]
          - name: b
            image: argoproj/argosay:v2
          - name: main
            image: argoproj/argosay:v2
            args: [ cat, /workpsace/message ]
            dependencies:
              - a
              - b
      outputs:
        parameters:
          - name: message
            valueFrom:
              path: /workpsace/message
```

There are a couple of caveats:

1. You must use the [Emissary Executor](workflow-executors.md#emissary-emissary).

The containers can be arranged as a graph by specifying dependencies, however as this is only suitable for running 10s
of containers.

At a container level, there is no support for:

* Conditions (`when`).
* Looping (`withSequence`, `withParams`).
* Enhanced `depends` logic.
* Retry (the whole pod will succeed or fail).
* Resource reporting.
* N/M progress.

## Inputs and Outputs

As with the container and script templates, inputs and outputs can only be loaded and saved from a container
named `main`.

All pod templates that have artifacts must/should have a container named `main`.

If you want to use base-layer artifacts, `main` must be last to finish, so it must be the root node in the graph.

That is rarely going to be practical.

Instead, have a workspace volume and make sure all artifacts paths are on that volume.