/*
Copyright 2022.

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

package monitor

import (
	"context"

	util "github.com/hussnain612/uptime-robot-operator/controllers/utils"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"

	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	"github.com/hussnain612/uptime-robot-operator/api/v1alpha1"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	MonitorFinalizer string = "uptimerobot-controller/monitor"
	ApiURL           string = "https://api.uptimerobot.com/v2/"
)

//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;create;update;
//+kubebuilder:rbac:groups=uptime.uptime.com,resources=monitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptime.uptime.com,resources=monitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptime.uptime.com,resources=monitors/finalizers,verbs=update

func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("monitor", req.NamespacedName)
	log.Info("Monitor reconciler called")

	// Fetch the monitor instance
	monitor := &v1alpha1.Monitor{}

	err := r.Get(ctx, req.NamespacedName, monitor)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcilerUtil.RequeueWithError(err)
	}

	apiKey, err := util.GetApiKeyFromSecret(ctx, r.Client, log)
	if err != nil {
		return reconcilerUtil.RequeueWithError(err)
	}
	// Resource is marked for deletion
	if monitor.DeletionTimestamp != nil {
		log.Info("Deletion timestamp found for monitor " + req.Name)
		if finalizerUtil.HasFinalizer(monitor, MonitorFinalizer) {
			return r.handleDelete(ctx, req, monitor, apiKey)
		}
		// Finalizer doesn't exist so clean up is already done
		return reconcilerUtil.DoNotRequeue()
	}

	// Add finalizer if it doesn't exist
	if !finalizerUtil.HasFinalizer(monitor, MonitorFinalizer) {
		log.Info("Adding finalizer for instance " + req.Name)
		patchBase := client.MergeFrom(monitor.DeepCopy())

		finalizerUtil.AddFinalizer(monitor, MonitorFinalizer)

		err := r.Client.Patch(ctx, monitor, patchBase)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, monitor, err, false)
		}
		return ctrl.Result{}, nil
	}

	return r.handleCreate(ctx, req, monitor, apiKey)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Monitor{}).
		Complete(r)
}
