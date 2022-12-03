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

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var monitorlog = logf.Log.WithName("monitor-resource")

func (r *Monitor) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-uptime-uptime-com-v1alpha1-monitor,mutating=false,failurePolicy=fail,sideEffects=None,groups=uptime.uptime.com,resources=monitors,verbs=create;update,versions=v1alpha1,name=vmonitor.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Monitor{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Monitor) ValidateCreate() error {
	monitorlog.Info("validate create", "name", r.Name)

	return r.validateMonitorSubtype()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Monitor) ValidateUpdate(old runtime.Object) error {
	monitorlog.Info("validate update", "name", r.Name)

	return r.validateMonitorSubtype()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Monitor) ValidateDelete() error {
	// monitorlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *Monitor) validateMonitorSubtype() error {
	if r.Spec.MonitorType != 4 {
		return nil
	}

	if r.Spec.MonitorSubtype == 0 {
		return fmt.Errorf("monitor subtype is required for port monitoring")
	}

	if r.Spec.MonitorSubtype != 99 {
		return nil
	}

	if r.Spec.MonitorPort == 0 {
		return fmt.Errorf("monitor port is required for custom port monitoring")
	}

	return nil
}
