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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	reflect "reflect"
	sync "sync"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
	versioned "kubesphere.io/kubesphere/pkg/client/clientset/versioned"
	application "kubesphere.io/kubesphere/pkg/client/informers/externalversions/application"
	auditing "kubesphere.io/kubesphere/pkg/client/informers/externalversions/auditing"
	cluster "kubesphere.io/kubesphere/pkg/client/informers/externalversions/cluster"
	devops "kubesphere.io/kubesphere/pkg/client/informers/externalversions/devops"
	gateway "kubesphere.io/kubesphere/pkg/client/informers/externalversions/gateway"
	iam "kubesphere.io/kubesphere/pkg/client/informers/externalversions/iam"
	internalinterfaces "kubesphere.io/kubesphere/pkg/client/informers/externalversions/internalinterfaces"
	network "kubesphere.io/kubesphere/pkg/client/informers/externalversions/network"
	notification "kubesphere.io/kubesphere/pkg/client/informers/externalversions/notification"
	quota "kubesphere.io/kubesphere/pkg/client/informers/externalversions/quota"
	servicemesh "kubesphere.io/kubesphere/pkg/client/informers/externalversions/servicemesh"
	storage "kubesphere.io/kubesphere/pkg/client/informers/externalversions/storage"
	tenant "kubesphere.io/kubesphere/pkg/client/informers/externalversions/tenant"
	types "kubesphere.io/kubesphere/pkg/client/informers/externalversions/types"
	virtualization "kubesphere.io/kubesphere/pkg/client/informers/externalversions/virtualization"
)

// SharedInformerOption defines the functional option type for SharedInformerFactory.
type SharedInformerOption func(*sharedInformerFactory) *sharedInformerFactory

type sharedInformerFactory struct {
	client           versioned.Interface
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	lock             sync.Mutex
	defaultResync    time.Duration
	customResync     map[reflect.Type]time.Duration

	informers map[reflect.Type]cache.SharedIndexInformer
	// startedInformers is used for tracking which informers have been started.
	// This allows Start() to be called multiple times safely.
	startedInformers map[reflect.Type]bool
}

// WithCustomResyncConfig sets a custom resync period for the specified informer types.
func WithCustomResyncConfig(resyncConfig map[v1.Object]time.Duration) SharedInformerOption {
	return func(factory *sharedInformerFactory) *sharedInformerFactory {
		for k, v := range resyncConfig {
			factory.customResync[reflect.TypeOf(k)] = v
		}
		return factory
	}
}

// WithTweakListOptions sets a custom filter on all listers of the configured SharedInformerFactory.
func WithTweakListOptions(tweakListOptions internalinterfaces.TweakListOptionsFunc) SharedInformerOption {
	return func(factory *sharedInformerFactory) *sharedInformerFactory {
		factory.tweakListOptions = tweakListOptions
		return factory
	}
}

// WithNamespace limits the SharedInformerFactory to the specified namespace.
func WithNamespace(namespace string) SharedInformerOption {
	return func(factory *sharedInformerFactory) *sharedInformerFactory {
		factory.namespace = namespace
		return factory
	}
}

// NewSharedInformerFactory constructs a new instance of sharedInformerFactory for all namespaces.
func NewSharedInformerFactory(client versioned.Interface, defaultResync time.Duration) SharedInformerFactory {
	return NewSharedInformerFactoryWithOptions(client, defaultResync)
}

// NewFilteredSharedInformerFactory constructs a new instance of sharedInformerFactory.
// Listers obtained via this SharedInformerFactory will be subject to the same filters
// as specified here.
// Deprecated: Please use NewSharedInformerFactoryWithOptions instead
func NewFilteredSharedInformerFactory(client versioned.Interface, defaultResync time.Duration, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) SharedInformerFactory {
	return NewSharedInformerFactoryWithOptions(client, defaultResync, WithNamespace(namespace), WithTweakListOptions(tweakListOptions))
}

// NewSharedInformerFactoryWithOptions constructs a new instance of a SharedInformerFactory with additional options.
func NewSharedInformerFactoryWithOptions(client versioned.Interface, defaultResync time.Duration, options ...SharedInformerOption) SharedInformerFactory {
	factory := &sharedInformerFactory{
		client:           client,
		namespace:        v1.NamespaceAll,
		defaultResync:    defaultResync,
		informers:        make(map[reflect.Type]cache.SharedIndexInformer),
		startedInformers: make(map[reflect.Type]bool),
		customResync:     make(map[reflect.Type]time.Duration),
	}

	// Apply all options
	for _, opt := range options {
		factory = opt(factory)
	}

	return factory
}

// Start initializes all requested informers.
func (f *sharedInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for informerType, informer := range f.informers {
		if !f.startedInformers[informerType] {
			go informer.Run(stopCh)
			f.startedInformers[informerType] = true
		}
	}
}

// WaitForCacheSync waits for all started informers' cache were synced.
func (f *sharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool {
	informers := func() map[reflect.Type]cache.SharedIndexInformer {
		f.lock.Lock()
		defer f.lock.Unlock()

		informers := map[reflect.Type]cache.SharedIndexInformer{}
		for informerType, informer := range f.informers {
			if f.startedInformers[informerType] {
				informers[informerType] = informer
			}
		}
		return informers
	}()

	res := map[reflect.Type]bool{}
	for informType, informer := range informers {
		res[informType] = cache.WaitForCacheSync(stopCh, informer.HasSynced)
	}
	return res
}

// InternalInformerFor returns the SharedIndexInformer for obj using an internal
// client.
func (f *sharedInformerFactory) InformerFor(obj runtime.Object, newFunc internalinterfaces.NewInformerFunc) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	informerType := reflect.TypeOf(obj)
	informer, exists := f.informers[informerType]
	if exists {
		return informer
	}

	resyncPeriod, exists := f.customResync[informerType]
	if !exists {
		resyncPeriod = f.defaultResync
	}

	informer = newFunc(f.client, resyncPeriod)
	f.informers[informerType] = informer

	return informer
}

// SharedInformerFactory provides shared informers for resources in all known
// API group versions.
type SharedInformerFactory interface {
	internalinterfaces.SharedInformerFactory
	ForResource(resource schema.GroupVersionResource) (GenericInformer, error)
	WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool

	Application() application.Interface
	Auditing() auditing.Interface
	Cluster() cluster.Interface
	Devops() devops.Interface
	Gateway() gateway.Interface
	Iam() iam.Interface
	Network() network.Interface
	Notification() notification.Interface
	Quota() quota.Interface
	Servicemesh() servicemesh.Interface
	Storage() storage.Interface
	Tenant() tenant.Interface
	Types() types.Interface
	Virtualization() virtualization.Interface
}

func (f *sharedInformerFactory) Application() application.Interface {
	return application.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Auditing() auditing.Interface {
	return auditing.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Cluster() cluster.Interface {
	return cluster.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Devops() devops.Interface {
	return devops.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Gateway() gateway.Interface {
	return gateway.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Iam() iam.Interface {
	return iam.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Network() network.Interface {
	return network.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Notification() notification.Interface {
	return notification.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Quota() quota.Interface {
	return quota.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Servicemesh() servicemesh.Interface {
	return servicemesh.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Storage() storage.Interface {
	return storage.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Tenant() tenant.Interface {
	return tenant.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Types() types.Interface {
	return types.New(f, f.namespace, f.tweakListOptions)
}

func (f *sharedInformerFactory) Virtualization() virtualization.Interface {
	return virtualization.New(f, f.namespace, f.tweakListOptions)
}
