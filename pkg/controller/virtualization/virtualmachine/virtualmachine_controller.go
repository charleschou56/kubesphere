package virtualmachine

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kvapi "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"

	vlzv1alpha1 "kubesphere.io/api/virtualization/v1alpha1"
)

const (
	controllerName        = "virtualmachine-controller"
	successSynced         = "Synced"
	messageResourceSynced = "VirtualMachine synced successfully"
)

// Reconciler reconciles a cnat object
type Reconciler struct {
	client.Client
	Logger                  logr.Logger
	Recorder                record.EventRecorder
	MaxConcurrentReconciles int
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	if r.Logger == nil {
		r.Logger = ctrl.Log.WithName("controllers").WithName(controllerName)
	}
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor(controllerName)
	}
	if r.MaxConcurrentReconciles <= 0 {
		r.MaxConcurrentReconciles = 1
	}
	return ctrl.NewControllerManagedBy(mgr).
		Named(controllerName).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: r.MaxConcurrentReconciles,
		}).
		For(&vlzv1alpha1.VirtualMachine{}).
		Complete(r)
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	klog.Infof("=== Reconciling VirtualMachine %s/%s", req.Namespace, req.Name)

	rootCtx := context.Background()
	vm := &vlzv1alpha1.VirtualMachine{}
	if err := r.Get(rootCtx, req.NamespacedName, vm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	vm_instance := vm.DeepCopy()

	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})

	// retrive default namespace.
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		klog.Infof("=== error in namespace : %v\n", err.Error())
		return ctrl.Result{}, err
	}

	// get the kubevirt client, using which kubevirt resources can be managed.
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		klog.Infof("=== cannot obtain KubeVirt client: %v\n", err)
		return ctrl.Result{}, err
	}

	if IsDeletionCandidate(vm_instance, vlzv1alpha1.VirtualMachineFinalizer) {
		klog.Infof("=== deleting VirtualMachine %s/%s", req.Namespace, req.Name)
		if err := deleteVirtualMachine(virtClient, namespace, vm_instance); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.removeFinalizer(vm_instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	if NeedToAddFinalizer(vm_instance, vlzv1alpha1.VirtualMachineFinalizer) {
		klog.Infof("=== adding finalizer for VirtualMachine %s/%s", req.Namespace, req.Name)
		if err := r.addFinalizer(vm_instance); err != nil {
			return ctrl.Result{}, err
		}

		createVirtualMachine(virtClient, namespace, vm_instance)
	}

	if !reflect.DeepEqual(vm, vm_instance) {
		if err := r.Update(rootCtx, vm_instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	r.Recorder.Event(vm, corev1.EventTypeNormal, successSynced, messageResourceSynced)
	return ctrl.Result{}, nil

}

func (c *Reconciler) addFinalizer(virtualmachine *vlzv1alpha1.VirtualMachine) error {
	clone := virtualmachine.DeepCopy()
	controllerutil.AddFinalizer(clone, vlzv1alpha1.VirtualMachineFinalizer)

	err := c.Update(context.Background(), clone)
	if err != nil {
		klog.V(3).Infof("Error adding  finalizer to virtualmachine %s: %v", virtualmachine.Name, err)
		return err
	}
	klog.V(3).Infof("Added finalizer to virtualmachine %s", virtualmachine.Name)
	return nil
}

func (c *Reconciler) removeFinalizer(virtualmachine *vlzv1alpha1.VirtualMachine) error {
	clone := virtualmachine.DeepCopy()
	controllerutil.RemoveFinalizer(clone, vlzv1alpha1.VirtualMachineFinalizer)
	err := c.Update(context.Background(), clone)
	if err != nil {
		klog.V(3).Infof("Error removing  finalizer from virtualmachine %s: %v", virtualmachine.Name, err)
		return err
	}
	klog.V(3).Infof("Removed protection finalizer from virtualmachine %s", virtualmachine.Name)
	return nil
}

func createVirtualMachine(virtClient kubecli.KubevirtClient, namespace string, vm_instance *vlzv1alpha1.VirtualMachine) error {

	running := true
	vm := &kvapi.VirtualMachine{
		TypeMeta: metav1.TypeMeta{
			Kind: "VirtualMachine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vm_instance.Spec.Name,
			Namespace: namespace,
		},
		Spec: kvapi.VirtualMachineSpec{
			Running: &running,
			Template: &kvapi.VirtualMachineInstanceTemplateSpec{
				Spec: kvapi.VirtualMachineInstanceSpec{
					Domain: kvapi.DomainSpec{
						Resources: kvapi.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(vm_instance.Spec.Memory),
							},
						},
						Devices: kvapi.Devices{
							Disks: []kvapi.Disk{
								{
									Name: "containerdisk",
									DiskDevice: kvapi.DiskDevice{
										Disk: &kvapi.DiskTarget{
											Bus: "virtio",
										},
									},
								},
							},
						},
					},
					Volumes: []kvapi.Volume{
						{
							Name: "containerdisk",
							VolumeSource: kvapi.VolumeSource{
								ContainerDisk: &kvapi.ContainerDiskSource{
									Image: "kubevirt/cirros-container-disk-demo:latest",
								},
							},
						},
					},
				},
			},
		},
	}

	createdVM, err := virtClient.VirtualMachine(namespace).Create(vm)
	if err != nil {
		klog.Infof(err.Error())
		return err
	}

	for {
		vm, err := virtClient.VirtualMachine(createdVM.Namespace).Get(createdVM.Name, &metav1.GetOptions{})
		if err != nil {
			klog.Infof(err.Error())
			return err
		}

		fmt.Printf("VM name: %v\n", vm.Name)
		fmt.Printf("VM Created Status: %v\n", vm.Status.Created)
		fmt.Printf("VM Ready Status: %v\n", vm.Status.Ready)

		if vm.Status.Ready {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func deleteVirtualMachine(virtClient kubecli.KubevirtClient, namespace string, vm_instance *vlzv1alpha1.VirtualMachine) error {
	err := virtClient.VirtualMachine(namespace).Delete(vm_instance.Name, &metav1.DeleteOptions{})
	if err != nil {
		klog.Infof(err.Error())
		return err
	}
	return nil
}

// IsDeletionCandidate checks if object is candidate to be deleted
func IsDeletionCandidate(obj metav1.Object, finalizer string) bool {
	return obj.GetDeletionTimestamp() != nil && ContainsString(obj.GetFinalizers(),
		finalizer, nil)
}

// NeedToAddFinalizer checks if need to add finalizer to object
func NeedToAddFinalizer(obj metav1.Object, finalizer string) bool {
	return obj.GetDeletionTimestamp() == nil && !ContainsString(obj.GetFinalizers(),
		finalizer, nil)
}

// ContainsString checks if a given slice of strings contains the provided string.
// If a modifier func is provided, it is called with the slice item before the comparation.
func ContainsString(slice []string, s string, modifier func(s string) string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
		if modifier != nil && modifier(item) == s {
			return true
		}
	}
	return false
}
