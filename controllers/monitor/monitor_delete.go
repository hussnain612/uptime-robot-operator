package monitor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hussnain612/uptime-robot-operator/api/v1alpha1"
	util "github.com/hussnain612/uptime-robot-operator/controllers/utils"
	"github.com/hussnain612/uptime-robot-operator/pkg/models"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
)

func (r *MonitorReconciler) handleDelete(ctx context.Context, req ctrl.Request, monitor *v1alpha1.Monitor, apiKey string) (ctrl.Result, error) {
	log := r.Log.WithValues("monitor", req.NamespacedName)
	log.Info("Deleting monitor: " + monitor.ObjectMeta.Name)

	// Remove monitor from uptime-robot
	res := models.UptimeMonitorNewMonitorResponse{}
	path := "deleteMonitor"

	parameters := map[string]string{}
	parameters["id"] = monitor.Status.MonitorID

	body, err := util.HandleAPIRequestForUptimeRobot(log, ApiURL, path, apiKey, parameters)
	if err != nil {
		return reconcilerUtil.RequeueWithError(err)
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error(err, "failed to unmarshal new monitor response body")
	}

	if res.Stat == "ok" {
		// Monitor added successfully
		log.Info(fmt.Sprintf("Monitor '%s' deleted", monitor.Name))
	} else {
		return reconcilerUtil.RequeueWithError(fmt.Errorf(res.Error.Message))
	}

	// Delete finalizer
	patchBase := client.MergeFrom(monitor.DeepCopy())
	finalizerUtil.DeleteFinalizer(monitor, MonitorFinalizer)
	log.Info("Finalizer removed for monitor : " + monitor.ObjectMeta.Name)

	// Patch monitor
	err = r.Client.Patch(ctx, monitor, patchBase)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, monitor, err, false)
	}

	return reconcilerUtil.DoNotRequeue()
}
