module gitlab.w6d.io/w6d/library/init-tools/tools/setsvc

go 1.13

replace (
	gitlab.w6d.io/w6d/library/log => gitlab.w6d.io/w6d/library/log.git v1.2.26
	k8s.io/api => k8s.io/api v0.16.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.8
	k8s.io/client-go => k8s.io/client-go v0.16.8
)

require (
	github.com/cakturk/go-netstat v0.0.0-20200220111822-e5b49efee7a5
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	gitlab.w6d.io/w6d/library/log v1.2.26
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	k8s.io/api v0.18.4 // indirect
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
)
