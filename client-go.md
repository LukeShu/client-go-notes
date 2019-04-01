If you're like me, you looked at the Kubernetes HTTP API, and thought
"this is simple enough."  Then you looked at the `k8s.io/client-go`
library to speak it, and thought "where did all this complexity come
from‽"

It may help to know that  Kubernetes was originally written in Java.
In parts of the library API design, it shows that the author was
thinking in Java.

# Generics / Code-Generation

You probably think that Go doesn't have generics.  Go is only missing
generics if you're afraid of code-generation; given a "generic"
template, you can generate an instance of it for each type you want it
to contain.  The Kubernetes authors are not afraid of code-gen.
There's a `k8s.io/code-generator` set of tools to give you generic
goodness.

Given the input types in the packages in
`k8s.io/api/API_GROUP/VERSION`, the code-generator creates the
following generics:

  - client
    * `k8s.io/client-go/kubernetes/scheme` ???
    * `K8s.io/client-go/kubernetes/typed/API_GROUP/VERSION` contains a
      client for RPC operations for the `API_GROUP/VERSION` GV.
    * `k8s.io/client-go/kubernetes` contains a *client-set* of all of
      the clients in `k8s.io/client-go/kubernetes/typed/...`.
  - lister
    * `k8s.io/client-go/listers/API_GROUP/VERSION` contains typed
      listers for the types in the `API_GROUP/VERSION` GV.
  - informer
    * `k8s.io/client-go/informers` contains an untyped informer
      (`GenericInformer`) and a convenience shared InformerFactory
      (`SharedInformerFactory`) that can bu used to create different
      typed informers.
    * `k8s.io/client-go/informers/API_GROUP/VERSION` contains typed
      informers for the types in the `API_GROUP/VERSION` GV.

For "untyped" things, rather than using `interface{}`, the Kubernetes
libraries use `k8s.io/apimachinery/pkg/runtime.Object`.

The code generator also creates a
`k8s.io/api/API_GROUP/VERSION/zz_generated.deepcopy.go` file for each
GV.

`code-generator` also can generate something called `Defaulter`s, but
I'm not sure what that is.

# Clients / Client-Sets

Each API Group/Version gets its own generated *client*; if you want
typing.  If you're OK without typing, there's
`k8s.io/client-go/dynamic`, which is a non-typed/dynamically-typed
client.

All of the generated typed individual clients (at
`k8s.io/client-go/kubernetes/typed/...` are bundled together in to a
*client-set* that makes it "easy" to use any of the individual
clients.  Upon further reflection, I'm not sure it actually makes
things any easier than using individual clients directly, except that
it maybe makes things more easily discoverable (especially from IDE
auto-complete).

The usefulness of client-sets is significantly damped if you ever need
to deal with a API group that's `k8s.io/api/...` /
`k8s.io/client-go/kubernetes/typed/...`.  That could be something like
a CRD, or it could even just be one of the more exotic parts of the
official k8s API, like `apiregistration.k8s.io/v1` /
`k8s.io/kube-aggregator/pkg/apis/apiregistration/v1` /
`k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset` (see
[apigroups.org](./apigroups.org) for a listing of the package names
for all of the API groups in the official k8s API).  You need a
different client-set for each place where the GV lives, and there's no
good tooling around combining client-sets.  So if you're having to
wrangle multiple client-sets, why not just cut-out some complexity and
wrangle multiple clients directly?

# Layering

Here's how things are layered:

    +-------------------------------------------------------------------------------------------------------------------+
    |                                        k8s.io/client-go/kubernetes.Clientset                                      |
    +--------------------------------------------------------+-+--------------------------------------------------------+
    | k8s.io/client-go/kubernetes/typed/core/v1.CoreV1Client | | k8s.io/client-go/kubernetes/typed/apps/v1.AppsV1Client |
    +--------------------------------------------------------+ +--------------------------------------------------------+
    |             k8s.io/client-go/rest.Interface            | |            k8s.io/client-go/rest.Interface             |
    |          (*k8s.io/client-go/rest.RESTClient)           | |          (*k8s.io/client-go/rest.RESTClient)           |
    |                 +--------------------------------------+ |                 +--------------------------------------+
    |                 |         net/http.RoundTripper        | |                 |         net/http.RoundTripper        |
    |                 | (k8s.io/client-go/rest.TransportFor) | |                 | (k8s.io/client-go/rest.TransportFor  |
    +-----------------+--------------------------------------+ +-----------------+------------------------------------- +
    |             k8s.io/client-go/rest.Config               | |             k8s.io/client-go/rest.Config               |
    +--------------------------------------------------------+-+--------------------------------------------------------+
    |                                            k8s.io/client-go/rest.Config                                           |
    +-------------------------------------------------------------------------------------------------------------------+

Each client gets its own `rest.RESTClient` that each has its own
`rest.Config` that is slightly modified for API-Group-specific things
from the shared base `rest.Config`

See https://github.com/datawire/k8scli/blob/master/client.go for the
low-levels of how that base `rest.Config` is created; it closely
mimics what `kubectl` does, but is much easier to read (and isn't
spread out across a dozen different files!).  The linked code uses
`k8s.io/client-go` to load the kubeconfig, and obtain a the
`rest.Config`, but then adapts it for use with
`github.com/ericchiang/k8s` instead of `k8s.io/client-go`; so you can
ignore everything after it obtains the restconfig.  For comparison,
here's what the `datawire/k8scli` layering looks like:

    +--------------------------------------+
    |    github.com/ericchiang/k8s.Client  |
    +--------------------------------------+
    |         net/http.RoundTripper        |
    | (k8s.io/client-go/rest.TransportFor) |
    +--------------------------------------+
    |    k8s.io/client-go/rest.Config      |
    +--------------------------------------+

# Misc

TODO: Discuss the relationship between clients and listers and
informers.

Other things under `k8s.io/client-go`
 - `discovery`
 - `discovery/cached`
 - `dynamic`
 - `dynamic/dynamicinformer`
 - `dynamic/dynamiclister`
 - `pkg/apis/clientauthentication`
 - `pkg/apis/clientauthentication/install`
 - `pkg/apis/clientauthentication/v1alpha1`
 - `pkg/apis/clientauthentication/v1alpha1`
 - `pkg/version`
 - `plugin/pkg/client/auth`
 - `plugin/pkg/client/auth/azure`
 - `plugin/pkg/client/auth/exec`
 - `plugin/pkg/client/auth/gcp`
 - `plugin/pkg/client/auth/oidc`
 - `rest`
 - `restmapper`
 - `scale`
 - `scale/scheme`
 - `scale/scheme/appsint`
 - `scale/scheme/appsv1beta1`
 - `scale/scheme/appsv1beta2`
 - `scale/scheme/autoscalingv1`
 - `scale/scheme/extensionsint`
 - `scale/scheme/extensionsv1beta1`
 - `testing`
 - `tools/auth`
 - `tools/cache`
 - `tools/clientcmd`
 - `tools/clientcmd/api`
 - `tools/clientcmd/api/latest`
 - `tools/clientcmd/api/v1`
 - `tools/leaderelection`
 - `tools/leaderelection/resourcelock`
 - `tools/metrics`
 - `tools/pager`
 - `tools/portforward`
 - `tools/record`
 - `tools/reference`
 - `tools/remotecommand`
 - `tools/watch`
 - `transport`
 - `transport/spdy`
 - `util/buffer`
 - `util/cert`
 - `util/certificate`
 - `util/certificate/csr`
 - `util/connrotation`
 - `util/exec`
 - `util/flowcontrol`
 - `util/homedir`
 - `util/jsonpath`
 - `util/retry`
 - `util/testing`
 - `util/workqueue`
