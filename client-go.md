If you're like me, you looked at the Kubernetes HTTP API, and thought
"this is simple enough."  Then you looked at the `k8s.io/client-go`
library to speak it, and though "where did all this complexity come
from‽"

It may help to know that  Kubernetes was originally written in Java.
In parts of the library API design, it shows that the author was
thinking in Java.

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
    * `k8s.io/client-go/kubernetes/typed/API_GROUP/VERSION` contains a
      client for RPC operations for the `API_GROUP/VERSION` GV.
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
