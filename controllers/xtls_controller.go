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
	"fmt"
	"strings"
	"time"

	xcagrpc "github.com/x-ca/go-grpc-api/grpc"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	xcav1alpha1 "github.com/kbcx/xca-operator/api/v1alpha1"
	"github.com/kbcx/xca-operator/utils"
)

// XtlsReconciler reconciles a Xtls object
type XtlsReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	XcaClient  xcagrpc.ServiceClient
	XcaContext context.Context
}

//+kubebuilder:rbac:groups=xca.kb.cx,resources=xtls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=xca.kb.cx,resources=xtls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=xca.kb.cx,resources=xtls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Xtls object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *XtlsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. check resource exist
	var obj xcav1alpha1.Xtls
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Xtls obj not found")
		} else {
			logger.Error(err, "unable to fetch obj from kube-apiserver")
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("fetch Xtls obj from kube-apiserver success")

	// 2. check Secret resource
	secretObj := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      strings.Replace(req.Name, "*", "_", 1),
	}
	var secret *v1.Secret
	if err := r.Get(ctx, secretObj, secret); err != nil {
		if errors.IsNotFound(err) {
			// 2.1 not create, do create Xtls obj
			logger.Info("Call x-ca GRPC API to create Xtls ...", "secretObj", secretObj)
			if secret, err = r.NewTLSSecret(ctx, obj, nil); err != nil {
				logger.Error(err, "create new Secret err", "Xtls", obj)
				return ctrl.Result{}, err
			} else {
				logger.Info("create new Secret success", "Secret", secret)
				return ctrl.Result{}, nil
			}
		} else {
			// 2.2 call kube-apiserver err
			logger.Error(err, "fetch Secret err", "secretObj", secretObj)
			return ctrl.Result{}, err
		}
	}

	// 3. tls secret is already exist, check
	logger.Info("get secret success", "secret", secret)
	// 3.1 check secret Type
	if secret.Type != v1.SecretTypeTLS {
		msg := fmt.Sprintf("secret name %s is already exist, but type is %s, expect type is %s",
			secret.Name, secret.Type, v1.SecretTypeTLS)
		err := errors.NewBadRequest(msg)
		logger.Error(err, "found wrong Type secret")
		return ctrl.Result{}, errors.NewConflict(schema.GroupResource{Group: "kubernetes.io", Resource: "tls"}, secret.Name, err)
	}

	cert, err := utils.ParseCert(secret.Data["tls.crt"])
	if err != nil {
		logger.Error(err, "parse cert error")
		return ctrl.Result{}, err
	}
	var certIPs []string
	for _, ip := range cert.IPAddresses {
		certIPs = append(certIPs, ip.String())
	}
	reGenerateCertFlag := false
	if cert.NotAfter.After(time.Now().Add(10 * 24 * time.Hour)) { // 3.2 check secret expire time, <= 10 days, regenerate cert
		logger.Info("cert will expire at %s, re-generate new cert", cert.NotAfter)
		reGenerateCertFlag = true
	} else if utils.IsMatch(obj.Spec.Domains, cert.DNSNames) == false { // 3.3 check cert Domains
		logger.Info("cert Domains is not match, cert Domains is %s, Xtls Domains is %s, re-generate new cert",
			cert.DNSNames, obj.Spec.Domains)
		reGenerateCertFlag = true
	} else if utils.IsMatch(obj.Spec.IPs, certIPs) == false { // 3.3 check cert IPs
		logger.Info("cert IPs is not match, cert IPs is %s, Xtls IPs is %s, re-generate new cert",
			cert.IPAddresses, obj.Spec.IPs)
		reGenerateCertFlag = true
	}
	if reGenerateCertFlag {
		if secret, err = r.NewTLSSecret(ctx, obj, secret); err != nil {
			logger.Error(err, "re-generate new Secret err", "Xtls", obj, "secretName", secret.Name)
			return ctrl.Result{}, err
		}
	}

	obj.Status.Active = true
	obj.Status.LastUptateTime = &metav1.Time{Time: time.Now()}
	if err := r.Status().Update(ctx, &obj); err != nil {
		logger.Error(err, "update Xtls status err")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *XtlsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&xcav1alpha1.Xtls{}).
		Complete(r)
}

// NewTLSSecret create new secret
func (r *XtlsReconciler) NewTLSSecret(ctx context.Context, x xcav1alpha1.Xtls, secret *v1.Secret) (*v1.Secret, error) {
	// 1. call X-CA GRPC API to create TLS
	tlsReq := xcagrpc.TLSRequest{
		CN:      "xiexianbin.cn",
		Domains: []string{"xiexianbin.cn", "*.xiexianbin.cn"},
		IPs:     []string{},
		Days:    10,
		KeyBits: 512,
	}
	tlsResp, err := r.XcaClient.Sign(r.XcaContext, &tlsReq)
	if err != nil {
		return nil, err
	}

	certData := map[string][]byte{
		"tls.crt": []byte(tlsResp.Cert),
		"tls.key": []byte(tlsResp.Key),
	}
	now := metav1.Now().Rfc3339Copy()
	if secret != nil {
		secret.ObjectMeta.Labels["update_at"] = now.String()
		secret.Data = certData
		if err := r.Update(ctx, secret); err != nil {
			return nil, err
		}
		return secret, nil
	} else {
		secret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      x.Name,
				Namespace: x.Namespace,
				Labels: map[string]string{
					"creator":   "xca-operator",
					"create_at": now.String(),
				},
			},
			Type: v1.SecretTypeTLS,
			Data: certData,
		}
		if err := r.Create(ctx, secret); err != nil {
			return nil, err
		}
		return secret, nil
	}
}
