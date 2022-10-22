# automatic-umbrella

proof-of-concept for running kustomize krm functions in kubernetes, potentially from within kubernetes.

Using a GitHub autogenerated name until I can give this a proper "short" name.

# Problem

Kustomize currently can run [containerized KRM functions][] only through `docker run`.
But this does not work well with GitOps tools like ArgoCD, which run on kubernetes, where `docker run` will not work.

# Alternatives

An alternative can be found [in this ArgoCD issue][argocd krm functions] regarding this problem.
My understanding is that relies on running nested containers using a privileged podman container.
This might work, but I would like to research how this can be done without escalated privileges
and be delegate to the platform.

# To Prove

I would like fo research and prove that running KRM Functions as Kuberneted Pods is possible with the official client libraries, specifically client-go.

* Level 1:
    * Passing stdin
    * Capturing stdout/stderr
    * ENV
* Level 2: mounts
* Level 3: networking restrictions

# Proof Concept

A go library that can run KRM functions in Kubernetes, through official client libraries.

A a (go) binary that can be used to replace the `docker` binary Kustomize uses for `docker run`, instead running a Kubernetes Pod using the above go library.

## Level 1

The library will use the standard client library means to schedule the Pod
and use `remotecommand` functionality to provide the `ResourceList` input and capture the result.

## Level 2

tbd.

mapping local fs contents using configmaps?
what about fs structures that are not flat?

using something like `kubectl cp`?
seems to rely on tar, requiring that would be questionable.

could we use our own tar container as an init container?
are init containers supported by remotecommand, can they have stdin?

## Level 3

tbd.

[Network Policies][]?

[containerized KRM functions]: https://kubectl.docs.kubernetes.io/guides/extending_kustomize/containerized_krm_functions/
[Network Policies]: https://kubernetes.io/docs/concepts/services-networking/network-policies/
[argocd krm functions]: https://github.com/argoproj/argo-cd/issues/5553#issuecomment-1135355355
