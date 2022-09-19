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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MonitorSpec defines the desired state of Monitor
type MonitorSpec struct {
	// Name of uptime monitor
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// URL to monitor
	URL string `json:"url"`

	// +kubebuilder:validation:Required
	// The uptimerobot monitor type (http, keyword, ping, port, heartbeat)
	// +kubebuilder:validation:Enum=http;keyword;ping;port;heartbeat
	MonitorType string `json:"monitorType,omitempty"`

	// The uptimerobot monitor subtype for port monitoring
	MonitorSubtype int `json:"monitorSubtype,omitempty"`
}

// MonitorStatus defines the observed state of Monitor
type MonitorStatus struct {
	// ID of uptime monitor
	MonitorID string `json:"monitorID"`

	// Status conditions for tenant
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Monitor is the Schema for the monitors API
type Monitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitorSpec   `json:"spec,omitempty"`
	Status MonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MonitorList contains a list of Monitor
type MonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitor{}, &MonitorList{})
}
