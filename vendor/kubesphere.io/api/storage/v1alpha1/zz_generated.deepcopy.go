//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CapabilityFeatures) DeepCopyInto(out *CapabilityFeatures) {
	*out = *in
	out.Volume = in.Volume
	out.Snapshot = in.Snapshot
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CapabilityFeatures.
func (in *CapabilityFeatures) DeepCopy() *CapabilityFeatures {
	if in == nil {
		return nil
	}
	out := new(CapabilityFeatures)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PluginInfo) DeepCopyInto(out *PluginInfo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PluginInfo.
func (in *PluginInfo) DeepCopy() *PluginInfo {
	if in == nil {
		return nil
	}
	out := new(PluginInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisionerCapability) DeepCopyInto(out *ProvisionerCapability) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisionerCapability.
func (in *ProvisionerCapability) DeepCopy() *ProvisionerCapability {
	if in == nil {
		return nil
	}
	out := new(ProvisionerCapability)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ProvisionerCapability) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisionerCapabilityList) DeepCopyInto(out *ProvisionerCapabilityList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ProvisionerCapability, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisionerCapabilityList.
func (in *ProvisionerCapabilityList) DeepCopy() *ProvisionerCapabilityList {
	if in == nil {
		return nil
	}
	out := new(ProvisionerCapabilityList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ProvisionerCapabilityList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisionerCapabilitySpec) DeepCopyInto(out *ProvisionerCapabilitySpec) {
	*out = *in
	out.PluginInfo = in.PluginInfo
	out.Features = in.Features
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisionerCapabilitySpec.
func (in *ProvisionerCapabilitySpec) DeepCopy() *ProvisionerCapabilitySpec {
	if in == nil {
		return nil
	}
	out := new(ProvisionerCapabilitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SnapshotFeature) DeepCopyInto(out *SnapshotFeature) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SnapshotFeature.
func (in *SnapshotFeature) DeepCopy() *SnapshotFeature {
	if in == nil {
		return nil
	}
	out := new(SnapshotFeature)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageClassCapability) DeepCopyInto(out *StorageClassCapability) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageClassCapability.
func (in *StorageClassCapability) DeepCopy() *StorageClassCapability {
	if in == nil {
		return nil
	}
	out := new(StorageClassCapability)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageClassCapability) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageClassCapabilityList) DeepCopyInto(out *StorageClassCapabilityList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StorageClassCapability, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageClassCapabilityList.
func (in *StorageClassCapabilityList) DeepCopy() *StorageClassCapabilityList {
	if in == nil {
		return nil
	}
	out := new(StorageClassCapabilityList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageClassCapabilityList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageClassCapabilitySpec) DeepCopyInto(out *StorageClassCapabilitySpec) {
	*out = *in
	out.Features = in.Features
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageClassCapabilitySpec.
func (in *StorageClassCapabilitySpec) DeepCopy() *StorageClassCapabilitySpec {
	if in == nil {
		return nil
	}
	out := new(StorageClassCapabilitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeFeature) DeepCopyInto(out *VolumeFeature) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeFeature.
func (in *VolumeFeature) DeepCopy() *VolumeFeature {
	if in == nil {
		return nil
	}
	out := new(VolumeFeature)
	in.DeepCopyInto(out)
	return out
}
