/*
Copyright 2022 xiexianbin.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XtlsSpec defines the desired state of Xtls
type XtlsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Xtls. Edit xtls_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	CN      string   `json:"cn"`
	Domains []string `json:"domains"`
	IPs     []string `json:"ips,omitempty"`

	// +kubebuilder:validation:Max=10000
	// +kubebuilder:validation:Min=1
	Days    int64 `json:"days,omitempty"`
	KeyBits int64 `json:"keyBits,omitempty"`
}

// XtlsStatus defines the observed state of Xtls
type XtlsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Active          bool         `json:"active,omitempty"`
	LastRequestTime *metav1.Time `json:"lastRequestTime,omitempty"`
	LastUpdateTime  *metav1.Time `json:"lastUpdateTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Xtls is the Schema for the xtls API
type Xtls struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XtlsSpec   `json:"spec,omitempty"`
	Status XtlsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// XtlsList contains a list of Xtls
type XtlsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Xtls `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Xtls{}, &XtlsList{})
}
