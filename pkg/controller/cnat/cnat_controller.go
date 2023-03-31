package cnat

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cnatv1alpha1 "kubesphere.io/api/cnat/v1alpha1"
)

const (
	controllerName        = "cnat-controller"
	successSynced         = "Synced"
	messageResourceSynced = "At synced successfully"
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
		For(&cnatv1alpha1.At{}).
		Complete(r)
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	klog.Infof("=== Reconciling At")

	rootCtx := context.Background()
	at := &cnatv1alpha1.At{}
	if err := r.Get(rootCtx, req.NamespacedName, at); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Clone because the original object is owned by the lister.
	instance := at.DeepCopy()

	// If no phase set, default to pending (the initial phase):
	if instance.Status.Phase == "" {
		instance.Status.Phase = cnatv1alpha1.PhasePending
	}

	// Now let's make the main case distinction: implementing
	// the state diagram PENDING -> RUNNING -> DONE
	switch instance.Status.Phase {
	case cnatv1alpha1.PhasePending:
		klog.Infof("instance: phase=PENDING")
		// As long as we haven't executed the command yet,  we need to check if it's time already to act:
		klog.Infof("instance: checking schedule %q", instance.Spec.Schedule)
		// Check if it's already time to execute the command with a tolerance of 2 seconds:
		d, err := timeUntilSchedule(instance.Spec.Schedule)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("schedule parsing failed: %v", err))
			// Error reading the schedule - requeue the request:
			return ctrl.Result{}, err
		}
		klog.Infof("instance : schedule parsing done: diff=%v", d)
		if d > 0 {
			// Not yet time to execute the command, wait until the scheduled time
			return ctrl.Result{RequeueAfter: d}, nil
		}

		klog.Infof("instance: it's time! Ready to execute: %s", instance.Spec.Command)
		instance.Status.Phase = cnatv1alpha1.PhaseRunning
	case cnatv1alpha1.PhaseRunning:
		klog.Infof("instance: Phase: RUNNING")

		pod := newPodForCR(instance)

		// Set At instance as the owner and controller
		owner := metav1.NewControllerRef(instance, cnatv1alpha1.SchemeGroupVersion.WithKind("At"))
		pod.ObjectMeta.OwnerReferences = append(pod.ObjectMeta.OwnerReferences, *owner)

		// Try to see if the pod already exists and if not
		// (which we expect) then create a one-shot pod as per spec:
		err := r.Get(rootCtx, req.NamespacedName, pod)
		if err != nil && errors.IsNotFound(err) {
			err = r.Create(rootCtx, pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			klog.Infof("instance: pod launched: name=%s", pod.Name)
			instance.Status.Phase = cnatv1alpha1.PhaseDone
		} else if err != nil {
			// requeue with error
			return ctrl.Result{}, err
		} else {
			// don't requeue because it will happen automatically when the pod status changes
			return ctrl.Result{}, nil
		}
	case cnatv1alpha1.PhaseDone:
		klog.Infof("instance: phase: DONE")
		return ctrl.Result{}, nil
	default:
		klog.Infof("instance: NOP")
		return ctrl.Result{}, nil
	}

	if !reflect.DeepEqual(at, instance) {
		// Update the At instance, setting the status to the respective phase:
		if err := r.Update(rootCtx, instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	r.Recorder.Event(at, corev1.EventTypeNormal, successSynced, messageResourceSynced)
	return ctrl.Result{}, nil
}

// timeUntilSchedule parses the schedule string and returns the time until the schedule.
// When it is overdue, the duration is negative.
func timeUntilSchedule(schedule string) (time.Duration, error) {
	now := time.Now().UTC()
	layout := "2006-01-02T15:04:05Z"
	s, err := time.Parse(layout, schedule)
	if err != nil {
		return time.Duration(0), err
	}
	return s.Sub(now), nil
}

func newPodForCR(cr *cnatv1alpha1.At) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: strings.Split(cr.Spec.Command, " "),
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}
