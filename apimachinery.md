At the core of the Kubernetes implementation is a run-time type-system
that lives in `k8s.io/apimachinery`.  It would honestly be a pretty
swell library that I think would be popular in the Go ecosystem, if
they put just a little effort in to cleaning up the API a bit.  But
for now, it's a mess because it's internal to Kubernetes.

Some references:
 - https://medium.com/@arschles/go-experience-report-generics-in-kubernetes-25da87430301
   There's insight in the article about the k8s implementation design.
   A lot of the info is now wrong because it's outdated.  Some of the
   info was always wrong.  But still insightful.
 - https://github.com/kubernetes/client-go/issues/193 is a reasonable
   follow-up.

This type-system is most visible in that it is used to handle
different resource versions in the api-server.  However, it is used
for a bunch of things throughout Kubernetes.

# A case study

This type-system is used by `kubectl` to handle parsing the
kubeconfig; the schema ("api") of the file is in
`k8s.io/client-go/tools/clientcmd/api`.

           "one stop shop" for building a client, which
           means that it's a dumping ground for all kinds
        ,- of crap someone was too lazy to find a proper  ,- register.go: api.SchemeGroupVersion = GV{G:"", V:"__internal"}
        |  home for.                                      |- register.go: api.SchemeBuilder      = (addKnownTypes)
        |                                             ,---+- register.go: api.AddToScheme        = func(s *runtime.Scheme) error { ... }
        |                                            /    |- types.go: types for GV{G:"", V:"__internal"}
    ,--' `-------------------------,                /     `- helpers.go: util functions: IsEmptyConfig, MinifyConfig, ShortenConfig, FlattenConfig, FlattenContent, ResolvePath, MakeAbs
    k8s.io/client-go/tools/clientcmd/api        <--'
    k8s.io/client-go/tools/clientcmd/api/v1     <--,
    k8s.io/client-go/tools/clientcmd/api/latest <-, \     ,- register.go: v1.SchemeGroupVersion = GV{G:"", V:"v1"}
                                                  |  \    |- register.go: v1.SchemeBuilder      = (addKnownTypes, addConversionFuncs)
      ,-------------------------------------------'   `---+- register.go: v1.AddToScheme        = (func(s *runtime.Scheme) error)(SchemeBuilder.AttToScheme)
      |                                                   |- types.go: types for GV{G:"", V:"v1"}
      |                                                   `- conversion.go: conversion functions; part of .SchemeBuilder
      | ,- latest.Version         = "v1"
      | |- latest.ExternalVersion = GV{G:"", V:"v1"}
      | |- latest.OldestVersion   = "v1"
      | |- latest.Versions        = []string{"v1"}
      `-+- latest.Scheme          = runtime.NewScheme()
        |    |- api.AddToScheme(Scheme)
        |    `- v1.AddToScheme(Scheme)
        |- latest.yamlSerializer  = json.NewYAMLSerializer(json.DefaultMetaFactory, Scheme, Scheme),
        `- latest.Codec           = versioning.NewDefaultCodecForScheme(
                                        /* scheme        = */ Scheme,
                                        /* encoder       = */ yamlSerializer,
                                        /* decoder       = */ yamlSerializer,
                                        /* encodeVersion = */ GV{G:"", V:"v1"},
                                        /* decodeVersion = */ GV{G:"", V:"__internal"})

Here we see two versions of the GV: `v1` and `__internal`.  They are
mostly the same, except that `__internal` has programming niceties,
like using a map instead of a list of K/V pairs.  IMO,
`GV{"","__internal"}` should live in
`k8s.io/client-go/tools/clientcmd/api/internalversion` (like with
k8s.io/apiextensions-apiserver).

There are schools of thought that it would be better to not have an
internal version for things, and that the external version should be
used internally; so that there are not two things to learn.  That is
not the opinion of the Kubernetes developers.

# Deserializing magic

`k8s.io/apimachinery/pkg/runtime/serializer/json.DefaultMetaFactory`

The magic is in
https://github.com/kubernetes/apimachinery/blob/a1e35b736404f6ced690f255586e2e98ad371244/pkg/runtime/serializer/json/meta.go#L46-L63
It then takes the returned schema.GroupVersionKind and looks it up in
a runtime.Scheme, and creates an instance of the appropriate type
using `scheme.New(GVK)`, which assumes that the appropriate GVK→type
mapping had been added with `scheme.AddKnownTypes(GV, prototypes...)`

If there are multiple versions that can be converted between, you use
k8s.io/apimachinery/pkg/runtime/serializer/versioning.Something to
parse it like normal, but then if the version doesn't match your
preferred version, it knows how to convert it to, because you
registered a converter with `scheme.AddConversionFuncs(func(in *FooV1,
out *FooV2, scope conversion.Scope) { ... })`
