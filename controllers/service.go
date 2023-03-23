package controllers

import (
	"context"

	servicesv1alpha1 "github.com/aescanero/openldap-operator/api/v1alpha1"
	"github.com/aescanero/openldap-operator/controllers/variables"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (reconciler *OpenldapReconciler) defineService(openldap *servicesv1alpha1.Openldap) *v1.Service {
	labels := map[string]string{variables.LabelKey: variables.LabelValue}

	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      variables.ServiceName,
			Namespace: openldap.Namespace,
			Labels:    labels},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Ports: []v1.ServicePort{{
				Port:     variables.Port,
				Protocol: "TCP",
				TargetPort: intstr.IntOrString{
					IntVal: variables.TargetPort,
				},
			},
				{
					Port:     variables.TLSPort,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						IntVal: variables.TargetTLSPort,
					},
				}},
			Selector: labels,
		},
	}

	ctrl.SetControllerReference(openldap, service, reconciler.Scheme)
	return service
}

func (reconciler *OpenldapReconciler) reconcileService(ctx context.Context, openldap *servicesv1alpha1.Openldap) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	serviceDefinition := reconciler.defineService(openldap)
	service := &v1.Service{}
	err := reconciler.Get(ctx, types.NamespacedName{Name: variables.ServiceName, Namespace: openldap.Namespace}, service)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Service resource " + variables.ServiceName + " not found. Creating or re-creating Service")
			err = reconciler.Create(ctx, serviceDefinition)
			if err != nil {
				log.Info("Failed to create Service resource. Re-running reconcile.")
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get Service resource " + variables.ServiceName + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes services are not updated
		log.Info("")
	}
	return ctrl.Result{}, nil
}
