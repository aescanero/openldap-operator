/*
Copyright 2023.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	servicesv1alpha1 "github.com/aescanero/openldap-operator/api/v1alpha1"
)

// OpenldapReconciler reconciles a Openldap object
type OpenldapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=services.disasterproject.com,resources=openldaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=services.disasterproject.com,resources=openldaps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=services.disasterproject.com,resources=openldaps/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Openldap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *OpenldapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Reconcile started for Openldap CRD")

	//Searching for ldap object
	openldap := &servicesv1alpha1.Openldap{}
	err := r.Get(ctx, req.NamespacedName, openldap)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Openldap resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Info("Failed to get Openldap resource. Re-running reconcile.")
		return ctrl.Result{}, err
	}

	_, err = r.reconcileClusterRole(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.reconcileClusterRoleBinding(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.reconcileSecret(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.reconcileService(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.reconcileServiceHeadLess(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.reconcileStatefulSet(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = r.reconcileService(ctx, openldap)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpenldapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicesv1alpha1.Openldap{}).
		Complete(r)
}
