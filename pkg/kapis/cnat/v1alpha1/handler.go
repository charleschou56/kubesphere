package v1alpha1

import (
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/api/errors"
	"kubesphere.io/kubesphere/pkg/api"
	"kubesphere.io/kubesphere/pkg/client/informers/externalversions"
	cnatlister "kubesphere.io/kubesphere/pkg/client/listers/cnat/v1alpha1"
)

type handler struct {
	atLister cnatlister.AtLister
}

func newHandler(ksInformers externalversions.SharedInformerFactory) *handler {
	return &handler{
		atLister: ksInformers.Cnat().V1alpha1().Ats().Lister(),
	}
}

func (h *handler) HelloCnat(request *restful.Request, response *restful.Response) {

	at, err := h.atLister.Ats("default").Get("example-at")

	if err != nil {
		if errors.IsNotFound(err) {
			api.HandleNotFound(response, request, err)
			return
		} else {
			api.HandleInternalError(response, request, err)
			return
		}
	}

	instance := at.DeepCopy()

	response.WriteAsJson(CnatResponse{
		Schedule: "At spec schedule = " + instance.Spec.Schedule,
		Command:  "At spec command = " + instance.Spec.Command,
	})
}

type CnatResponse struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
}
