package controllers

import (
	"context"

	servicesv1alpha1 "github.com/aescanero/openldap-operator/api/v1alpha1"
	"github.com/aescanero/openldap-operator/controllers/variables"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (reconciler *OpenldapReconciler) defineClusterRole(openldap *servicesv1alpha1.Openldap) *v1.ClusterRole {
	labels := map[string]string{variables.LabelKey: variables.LabelValue}

	clusterRole := &v1.ClusterRole{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ClusterRole"},

		ObjectMeta: metav1.ObjectMeta{Name: variables.ClusterRoleName, Namespace: openldap.Namespace, Labels: labels},
		Rules: []v1.PolicyRule{{
			APIGroups: []string{"services.disasterproject.com"},
			Resources: []string{"openldap"},
			Verbs:     []string{"create", "delete", "get", "list", "watch", "patch", "update"},
		},
			{APIGroups: []string{"services.disasterproject.com"},
				Resources: []string{"openldap/status"},
				Verbs:     []string{"create", "delete", "get", "list", "watch", "patch", "update"},
			}},
	}

	ctrl.SetControllerReference(openldap, clusterRole, reconciler.Scheme)
	return clusterRole
}

func (reconciler *OpenldapReconciler) reconcileClusterRole(ctx context.Context, openldap *servicesv1alpha1.Openldap) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	clusterRoleDefinition := reconciler.defineClusterRole(openldap)
	clusterRole := &v1.ClusterRole{}
	err := reconciler.Get(ctx, types.NamespacedName{Name: variables.ClusterRoleName, Namespace: openldap.Namespace}, clusterRole)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ClusterRole resource " + variables.ClusterRoleName + " not found. Creating or re-creating ClusterRole")
			err = reconciler.Create(ctx, clusterRoleDefinition)
			if err != nil {
				log.Info("Failed to create ClusterRole resource. Re-running reconcile.")
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get ClusterRole resource " + variables.ClusterRoleName + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes services are not updated
		log.Info("")
	}
	return ctrl.Result{}, nil
}
