package v1alpha1

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubesphere.io/kubesphere/pkg/api"
	"kubesphere.io/kubesphere/pkg/apiserver/runtime"
	"kubesphere.io/kubesphere/pkg/client/informers/externalversions"
)

const (
	groupName = "cnat.programming-kubernetes.info"
)

var GroupVersion = schema.GroupVersion{Group: groupName, Version: "v1alpha1"}

func AddToContainer(container *restful.Container, ksInformers externalversions.SharedInformerFactory) error {
	webservice := runtime.NewWebService(GroupVersion)
	handler := newHandler(ksInformers)

	webservice.Route(webservice.GET("/cnat").
		Reads("").
		To(handler.HelloCnat).
		Returns(http.StatusOK, api.StatusOK, CnatResponse{})).
		Doc("Api for cnat")

	container.Add(webservice)

	return nil
}
