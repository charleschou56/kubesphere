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
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Clone because the original object is owned by the lister.
	vm_instance := vm.DeepCopy()

	// If no phase set, default to pending (the initial phase):
	if vm_instance.Status.Phase == "" {
		vm_instance.Status.Phase = vlzv1alpha1.PhasePending
	}

	switch vm_instance.Status.Phase {
	case vlzv1alpha1.PhasePending:
		klog.Infof("instance : phase=PENDING")

		vm_instance.Status.Phase = vlzv1alpha1.PhaseRunning
	case vlzv1alpha1.PhaseRunning:
		klog.Infof("instance : phase=RUNNING")

		clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})

		// retrive default namespace.
		namespace, _, err := clientConfig.Namespace()
		if err != nil {
			klog.Infof("error in namespace : %v\n", err.Error())
			break
		}

		// get the kubevirt client, using which kubevirt resources can be managed.
		virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
		if err != nil {
			klog.Infof("cannot obtain KubeVirt client: %v\n", err)
			break
		}

		createVirtualMachine(virtClient, namespace, vm_instance)

		vm_instance.Status.Phase = vlzv1alpha1.PhaseDone
	case vlzv1alpha1.PhaseDone:
		klog.Infof("instance : phase=DONE")

		vm_instance.Status.Phase = vlzv1alpha1.PhaseDone
	default:

	}

	if !reflect.DeepEqual(vm, vm_instance) {
		// Update the At instance, setting the status to the respective phase:
		if err := r.Update(rootCtx, vm_instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	r.Recorder.Event(vm, corev1.EventTypeNormal, successSynced, messageResourceSynced)
	return ctrl.Result{}, nil

}

func createVirtualMachine(virtClient kubecli.KubevirtClient, namespace string, vm_instance *vlzv1alpha1.VirtualMachine) error {

	running := true
	vm := &kvapi.VirtualMachine{
		TypeMeta: k8smetav1.TypeMeta{
			Kind: "VirtualMachine",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
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
		vm, err := virtClient.VirtualMachine(createdVM.Namespace).Get(createdVM.Name, &k8smetav1.GetOptions{})
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
