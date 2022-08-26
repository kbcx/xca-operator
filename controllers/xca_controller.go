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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1beta1 "github.com/kbcx/xca-operator/api/v1beta1"
)

// XcaReconciler reconciles a Xca object
type XcaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.xca.k8s.kb.cx,resources=xcas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.xca.k8s.kb.cx,resources=xcas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.xca.k8s.kb.cx,resources=xcas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Xca object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *XcaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx, "namespace", req.Namespace)

	// 1. check resource exist
	var obj appsv1beta1.Xca
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("fetch Xca not found", "name", obj.Name, "namespace", obj.Namespace)
			err = nil
		} else {
			logger.Error(err, "unable to fetch Xca")
		}

		//return ctrl.Result{}, client.IgnoreNotFound(err)
		return ctrl.Result{}, err
	}
	logger.Info("fetch Xca Object", "name", obj.Name, "namespace", obj.Namespace)

	// 2. check resource status
	if obj.Status.Active {
		logger.Info("Xca is ative", "name", obj.Name, "namespace", obj.Namespace)
		return ctrl.Result{}, nil
	}

	// 3. change to Active
	logger.Info("do something to set object to active", obj.Name, "namespace", obj.Namespace)

	obj.Status.Active = true
	obj.Status.LastUpdateTime = &metav1.Time{Time: time.Now()}
	if err := r.Status().Update(ctx, &obj); err != nil {
		logger.Error(err, "update object status err", obj.Name, "namespace", obj.Namespace)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *XcaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1beta1.Xca{}).
		Complete(r)
}