package cakephp

import (
	"context"

	cakephpv1alpha1 "github.com/demo/pkg/apis/cakephp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const sqlPort = 3306

func mysqlDeploymentName() string {
	return "mysql"
}

func mysqlServiceName() string {
	return "mysql"
}

func mysqlAuthName() string {
	return "mysql-auth"
}

func (r *ReconcileCakephp) mysqlAuthSecret(d *cakephpv1alpha1.Cakephp) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlAuthName(),
			Namespace: d.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "root",
			"password": "cakephp",
		},
	}
	controllerutil.SetControllerReference(d, secret, r.scheme)
	return secret
}

func (r *ReconcileCakephp) mysqlDeployment(d *cakephpv1alpha1.Cakephp) *appsv1.Deployment {
	labels := labels(d, "mysql")
	size := int32(1)

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
			Key:                  "password",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlDeploymentName(),
			Namespace: d.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name: "mysql-data",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					}},
					Containers: []corev1.Container{{
						Image: "mysql:5.7",
						Name:  "mysql-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mysql-data",
							MountPath: "/var/lib/mysql",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "cakephp",
							},
							{
								Name:  "MYSQL_DATABASE",
								Value: "cakephp",
							},
							{
								Name:      "MYSQL_USER",
								ValueFrom: userSecret,
							},
							{
								Name:      "MYSQL_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(d, dep, r.scheme)
	return dep
}

func (r *ReconcileCakephp) mysqlService(d *cakephpv1alpha1.Cakephp) *corev1.Service {
	labels := labels(d, "mysql")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlServiceName(),
			Namespace: d.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       3306,
				TargetPort: intstr.FromInt(sqlPort),
			}},
		},
	}

	controllerutil.SetControllerReference(d, s, r.scheme)
	return s
}

// Returns whether or not the MySQL deployment is running
func (r *ReconcileCakephp) isMysqlUp(d *cakephpv1alpha1.Cakephp) bool {
	deployment := &appsv1.Deployment{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      mysqlDeploymentName(),
		Namespace: d.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment mysql not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
