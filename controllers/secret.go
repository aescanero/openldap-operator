package controllers

import (
	"context"

	"github.com/aescanero/openldap-controller/utils"
	servicesv1alpha1 "github.com/aescanero/openldap-operator/api/v1alpha1"
	"github.com/aescanero/openldap-operator/controllers/variables"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (reconciler *OpenldapReconciler) defineSecret(openldap *servicesv1alpha1.Openldap) *v1.Secret {
	labels := map[string]string{variables.LabelKey: variables.LabelValue}
	stringData := map[string]string{
		"adminPassword": utils.Random(10),
	}
	secret := &v1.Secret{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: metav1.ObjectMeta{Name: variables.SecretName, Namespace: openldap.Namespace, Labels: labels},
		Immutable:  new(bool),
		Data:       map[string][]byte{},
		StringData: stringData,
		Type:       "Opaque",
	}

	ctrl.SetControllerReference(openldap, secret, reconciler.Scheme)
	return secret
}

func (reconciler *OpenldapReconciler) reconcileSecret(ctx context.Context, openldap *servicesv1alpha1.Openldap) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	SecretDefinition := reconciler.defineSecret(openldap)
	Secret := &v1.Secret{}
	err := reconciler.Get(ctx, types.NamespacedName{Name: variables.SecretName, Namespace: openldap.Namespace}, Secret)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Secret resource " + variables.SecretName + " not found. Creating or re-creating Secret")
			err = reconciler.Create(ctx, SecretDefinition)
			if err != nil {
				log.Info("Failed to create Secret resource. Re-running reconcile.")
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get Secret resource " + variables.SecretName + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes services are not updated
		log.Info("")
	}
	return ctrl.Result{}, nil
}
