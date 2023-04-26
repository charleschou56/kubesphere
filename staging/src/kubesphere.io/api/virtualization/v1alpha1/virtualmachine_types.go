package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +genclient:nonNamespaced

const (
	PhasePending = "PENDING"
	PhaseRunning = "RUNNING"
	PhaseDone    = "DONE"
)

const VirtualMachineFinalizer = "finalizers.virtualization.kubesphere.io/virtualmachine"

type Requests struct {
	Storage string `json:"storage,omitempty"`
}

type Resources struct {
	Requests Requests `json:"requests,omitempty"`
}

type DiskVolumeTemplateSpec struct {
	Resources Resources `json:"resources,omitempty"`
}

type DiskVolumeTemplateStatus struct {
}

type DiskVolumeTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DiskVolumeTemplateSpec   `json:"spec,omitempty"`
	Status DiskVolumeTemplateStatus `json:"status,omitempty"`
}

// VirtualMachineSpec defines the desired state of VirtualMachine
type VirtualMachineSpec struct {
	// DiskVolumeTemplate is the name of the DiskVolumeTemplate.
	DiskVolumeTemplates []DiskVolumeTemplate `json:"diskVolumeTemplates,omitempty"`
	// Name is the name of the VirtualMachine.
	Name string `json:"name,omitempty"`
	// Memory is the memory of the VirtualMachine.
	Memory string `json:"memory,omitempty"`
}

// VirtualMachineStatus defines the observed state of VirtualMachine
type VirtualMachineStatus struct {
	Phase string `json:"phase,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualMachine runs a vm at a given name.
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualMachineList contains a list of VirtualMachine
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{})
}
