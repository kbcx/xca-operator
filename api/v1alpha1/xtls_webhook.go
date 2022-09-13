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
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"net/netip"
	"net/url"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var xtlslog = logf.Log.WithName("xtls-resource")

func (r *Xtls) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-xca-kb-cx-v1alpha1-xtls,mutating=true,failurePolicy=fail,sideEffects=None,groups=xca.kb.cx,resources=xtls,verbs=create;update,versions=v1alpha1,name=mxtls.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Xtls{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Xtls) Default() {
	xtlslog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-xca-kb-cx-v1alpha1-xtls,mutating=false,failurePolicy=fail,sideEffects=None,groups=xca.kb.cx,resources=xtls,verbs=create;update,versions=v1alpha1,name=vxtls.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Xtls{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Xtls) ValidateCreate() error {
	xtlslog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.paramsValidate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Xtls) ValidateUpdate(old runtime.Object) error {
	xtlslog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	o := old.DeepCopyObject().(*Xtls)
	if r.Spec.CN != o.Spec.CN {
		return fmt.Errorf("xtls.Spec.CN %s not match witch old %s", r.Spec.CN, o.Spec.CN)
	}
	return r.paramsValidate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Xtls) ValidateDelete() error {
	xtlslog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Xtls) paramsValidate() error {
	if _, err := url.ParseRequestURI(r.Spec.CN); err != nil {
		return fmt.Errorf("xtls.Spec.CN %s is un-vaild", r.Spec.CN)
	}

	if r.Spec.Days < 1 {
		return fmt.Errorf("xtls.Spec.Days %d must greater than 1", r.Spec.Days)
	}

	for _, domain := range r.Spec.Domains {
		if _, err := url.ParseRequestURI(domain); err != nil {
			return fmt.Errorf("domain %s in xtls.Spec.Domains %s is un-vaild", domain, r.Spec.Domains)
		}
	}

	for _, ip := range r.Spec.IPs {
		if _, err := netip.ParseAddr(ip); err != nil {
			return fmt.Errorf("IP %s in xtls.Spec.IPs %s is un-vaild",
				ip, r.Spec.IPs)
		}
	}

	if r.Spec.Days%256 != 0 {
		return fmt.Errorf("xtls.Spec.Days %d must be divisibility by 256, like 256/512/1024/2048",
			r.Spec.Days)
	}

	return nil
}
