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

func (reconciler *OpenldapReconciler) defineClusterRoleBinding(openldap *servicesv1alpha1.Openldap) *v1.ClusterRoleBinding {
	labels := map[string]string{variables.LabelKey: variables.LabelValue}

	clusterRoleBinding := &v1.ClusterRoleBinding{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "ClusterRoleBinding"},
		ObjectMeta: metav1.ObjectMeta{Name: variables.ClusterRoleBindingName, Namespace: openldap.Namespace, Labels: labels},
		Subjects: []v1.Subject{{
			Kind:      "ServiceAccount",
			Name:      variables.RoleBindingServiceAccount,
			Namespace: openldap.Namespace,
		}},
		RoleRef: v1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     variables.ClusterRoleName,
		},
	}

	ctrl.SetControllerReference(openldap, clusterRoleBinding, reconciler.Scheme)
	return clusterRoleBinding
}

func (reconciler *OpenldapReconciler) reconcileClusterRoleBinding(ctx context.Context, openldap *servicesv1alpha1.Openldap) (ctrl.Result, error) {
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
