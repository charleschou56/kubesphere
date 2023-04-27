package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +genclient:nonNamespaced

const VirtualMachineFinalizer = "finalizers.virtualization.kubesphere.io/virtualmachine"

type Image struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type Source struct {
	Image Image `json:"image,omitempty"`
}

type Requests struct {
	Storage string `json:"storage,omitempty"`
}

type DVResources struct {
	Requests Requests `json:"requests,omitempty"`
}

type DiskVolumeTemplateSpec struct {
	// Resources represents the minimum resources the volume should have.
	Resources DVResources `json:"resources,omitempty"`
	Source    Source      `json:"source,omitempty"`
}

type DiskVolumeTemplateStatus struct {
}

type DiskVolumeTemplate struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// DiskVolumeSpec is the spec for a DiskVolume resource
	Spec   DiskVolumeTemplateSpec   `json:"spec,omitempty"`
	Status DiskVolumeTemplateStatus `json:"status,omitempty"`
}

type Cpu struct {
	Cores int32 `json:"cores,omitempty"`
}

type MacVtap struct {
}

type Interface struct {
	MacVtap MacVtap `json:"macvtap,omitempty"`
	Name    string  `json:"name,omitempty"`
}

type Devices struct {
	Interfaces []Interface `json:"interfaces,omitempty"`
}

type Limits struct {
	Memory string `json:"memory,omitempty"`
}

type DomainResources struct {
	Limits Limits `json:"limits,omitempty"`
}

type Domain struct {
	Cpu       Cpu             `json:"cpu,omitempty"`
	Devices   Devices         `json:"devices,omitempty"`
	Resources DomainResources `json:"resources,omitempty"`
}

type Multus struct {
	NetworkName string `json:"networkName,omitempty"`
}

type Network struct {
	Multus Multus `json:"multus,omitempty"`
	Name   string `json:"name,omitempty"`
}

type CloudInitNoCloud struct {
	UserDataBase64 string `json:"userDataBase64,omitempty"`
}

type Volume struct {
	CloudInitNoCloud CloudInitNoCloud `json:"cloudInitNoCloud,omitempty"`
	Name             string           `json:"name,omitempty"`
}

type Hardware struct {
	Domain           Domain    `json:"domain,omitempty"`
	EvictionStrategy string    `json:"evictionStrategy,omitempty"`
	Hostname         string    `json:"hostname,omitempty"`
	Networks         []Network `json:"networks,omitempty"`
	Volumes          []Volume  `json:"volumes,omitempty"`
}

// VirtualMachineSpec defines the desired state of VirtualMachine
type VirtualMachineSpec struct {
	// DiskVolumeTemplate is the name of the DiskVolumeTemplate.
	DiskVolumeTemplates []DiskVolumeTemplate `json:"diskVolumeTemplates,omitempty"`
	// DiskVolume is the name of the DiskVolume.
	DiskVolumes []string `json:"diskVolumes,omitempty"`
	// Name is the name of the VirtualMachine.
	Name string `json:"name,omitempty"`
	// Memory is the memory of the VirtualMachine.
	Memory string `json:"memory,omitempty"`
	// Hardware is the hardware of the VirtualMachine.
	Hardware Hardware `json:"hardware,omitempty"`
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
