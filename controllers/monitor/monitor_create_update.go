package monitor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hussnain612/uptime-robot-operator/api/v1alpha1"
	util "github.com/hussnain612/uptime-robot-operator/controllers/utils"
	"github.com/hussnain612/uptime-robot-operator/pkg/constants"
	"github.com/hussnain612/uptime-robot-operator/pkg/models"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *MonitorReconciler) handleCreate(ctx context.Context, req ctrl.Request, monitor *v1alpha1.Monitor, apiKey string) (ctrl.Result, error) {
	log := r.Log.WithValues("monitor", req.NamespacedName)
	log.Info("Creating/Updating monitor: " + monitor.ObjectMeta.Name)

	patchBase := client.MergeFrom(monitor.DeepCopy())

	var path string
	res := models.UptimeRobotMonitorResponse{}
	parameters := map[string]string{}

	// Check if monitor is already added to uptime-robot
	if monitor.Status.MonitorID != "" {
		// Update monitor
		path = "editMonitor"

		parameters["id"] = monitor.Status.MonitorID
		parameters = r.getMonitorParameters(*monitor, parameters)
	} else {
		// Add monitor to uptime-robot
		path = "newMonitor"

		parameters = r.getMonitorParameters(*monitor, parameters)
	}

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
		log.Info(fmt.Sprintf("Monitor '%s' added", monitor.Name))
		monitor.Status.MonitorID = fmt.Sprint(res.Monitor.ID)
	} else {
		return reconcilerUtil.RequeueWithError(fmt.Errorf(res.Error.Message))
	}

	err = r.patchMonitor(ctx, log, monitor, patchBase)
	if err != nil {
		return reconcilerUtil.RequeueWithError(err)
	}

	return reconcilerUtil.RequeueAfter(constants.MonitorDefaultRequeueTime)
}
