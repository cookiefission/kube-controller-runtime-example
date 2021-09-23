package main

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func main() {
	logf.SetLogger(zap.New())

	var log = logf.Log.WithName("controller-example")

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	err = builder.
		ControllerManagedBy(mgr).
		For(&appsv1.ReplicaSet{}).
		Complete(&ReplicaSetReconciler{})
	if err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}

// ReplicaSetReconciler must adhere to the sigs.k8s.io/controller-runtime/pkg/reconcile.Reconciler
// interface and so implement a Reconcile func, see below.
type ReplicaSetReconciler struct {
	client.Client
}

func (r *ReplicaSetReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	rs := &appsv1.ReplicaSet{}
	err := r.Get(ctx, req.NamespacedName, rs)
	if err != nil {
		return reconcile.Result{}, err
	}

	pods := &corev1.PodList{}
	err = r.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
	if err != nil {
		return reconcile.Result{}, err
	}

	rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
	err = r.Update(ctx, rs)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (a *ReplicaSetReconciler) InjectClient(c client.Client) error {
	a.Client = c
	return nil
}
