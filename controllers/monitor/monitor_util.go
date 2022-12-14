package monitor

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/hussnain612/uptime-robot-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
)

func (r *MonitorReconciler) patchMonitor(ctx context.Context, log logr.Logger, monitor *v1alpha1.Monitor, basePatch client.Patch) error {

	// Update status
	monitor.Status.Conditions = []metav1.Condition{
		{
			Type:               "ReconcileSuccess",
			LastTransitionTime: metav1.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), 0, 0, time.Now().Location()),
			Message:            reconcilerUtil.SuccessfulMessage,
			Reason:             reconcilerUtil.SuccessfulReason,
			Status:             metav1.ConditionTrue,
		},
	}

	// Patch status
	err := r.Client.Status().Patch(ctx, monitor, basePatch)
	if err != nil {
		log.Error(err, "failed to patch tenant status")
	}
	return err
}
