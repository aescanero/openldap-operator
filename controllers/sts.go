package controllers

import (
	"context"

	servicesv1alpha1 "github.com/aescanero/openldap-operator/api/v1alpha1"
	"github.com/aescanero/openldap-operator/controllers/variables"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (reconciler *OpenldapReconciler) defineStatefulSet(openldap *servicesv1alpha1.Openldap) *appsv1.StatefulSet {
	labels := map[string]string{variables.LabelKey: variables.LabelValue}
	LdifBase := openldap.Spec.Base
	if LdifBase == "" {
		LdifBase = variables.LdifBase
	}

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      variables.StatefulSetName,
			Namespace: openldap.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: variables.ServiceNameHeadLess,
			Replicas:    &variables.StatefulSetReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  variables.ContainerName,
							Image: variables.OpenldapImage,
							Ports: []v1.ContainerPort{
								{
									Name:          "ldap",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: variables.TargetPort,
								},
								{
									Name:          "ldaps",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: variables.TargetTLSPort,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{Name: "data", MountPath: "/var/lib/ldap/0"},
								{Name: "conf", MountPath: "/etc/ldap"},
							},
							Env: []v1.EnvVar{{
								Name: variables.SecretName,
								ValueFrom: &v1.EnvVarSource{
									SecretKeyRef: &v1.SecretKeySelector{
										LocalObjectReference: v1.LocalObjectReference{
											Name: variables.SecretName,
										},
										Key: variables.SecretName,
									},
								}},
								{
									Name:  "LDAP_BASE",
									Value: LdifBase,
								},
							},
							ReadinessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									Exec: &v1.ExecAction{Command: []string{"controller", "status"}},
								},
								InitialDelaySeconds: 20,
							},
							LivenessProbe: &v1.Probe{
								ProbeHandler: v1.ProbeHandler{
									Exec: &v1.ExecAction{Command: []string{"controller", "status"}},
								},
								InitialDelaySeconds: 40,
							},
						},
					},
					Volumes: []v1.Volume{
						{Name: "data", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{SizeLimit: &variables.GigaByte}}},
						{Name: "conf", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{SizeLimit: &variables.GigaByte}}},
					},
				},
			},
		},
	}

	ctrl.SetControllerReference(openldap, statefulset, reconciler.Scheme)
	return statefulset
}

func (reconciler *OpenldapReconciler) reconcileStatefulSet(ctx context.Context, openldap *servicesv1alpha1.Openldap) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	StatefulSetDefinition := reconciler.defineStatefulSet(openldap)
	StatefulSet := &appsv1.StatefulSet{}
	err := reconciler.Get(ctx, types.NamespacedName{Name: variables.StatefulSetName, Namespace: openldap.Namespace}, StatefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("StatefulSet resource " + variables.StatefulSetName + " not found. Creating or re-creating StatefulSet")
			err = reconciler.Create(ctx, StatefulSetDefinition)
			if err != nil {
				log.Info("Failed to create StatefulSet resource. Re-running reconcile.")
				return ctrl.Result{}, err
			}
		} else {
			log.Info("Failed to get StatefulSet resource " + variables.StatefulSetName + ". Re-running reconcile.")
			return ctrl.Result{}, err
		}
	} else {
		// Note: For simplication purposes statefulsets are not updated
		log.Info("")
	}
	return ctrl.Result{}, nil
}
