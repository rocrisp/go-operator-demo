package cakephp

import (
	"context"
	"time"

	routev1 "github.com/openshift/api/route/v1"
	cakephpv1alpha1 "github.com/rocrisp/go-operator-demo/pkg/apis/cakephp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const demoPort = 8080
const demoServicePort = 8080
const demoImage = "quay.io/rocrisp/cakedemo:v1"

func demoDeploymentName(d *cakephpv1alpha1.Cakephp) string {
	return d.Name + "-deployment"
}

func demoServiceName(d *cakephpv1alpha1.Cakephp) string {
	return d.Name + "-service"
}

func (r *ReconcileCakephp) demoDeployment(d *cakephpv1alpha1.Cakephp) *appsv1.Deployment {
	labels := labels(d, "demo")
	size := d.Spec.Size

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
			Name:      demoDeploymentName(d),
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
					Containers: []corev1.Container{{
						Image: demoImage,
						Name:  "cakephp-demo",
						Ports: []corev1.ContainerPort{{
							ContainerPort: demoPort,
							Name:          "cakephp",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "DATABASE_SERVICE_NAME",
								Value: "mysql",
							},
							{
								Name:  "DATABASE_NAME",
								Value: "cakephp",
							},
							{
								Name:  "FIRST_LASTNAME",
								Value: d.Spec.Title,
							},
							{
								Name:  "MYSQL_SERVICE_HOST",
								Value: mysqlServiceName(),
							},
							{
								Name:  "SESSION_DEFAULTS",
								Value: "database",
							},
							{
								Name:      "DATABASE_USER",
								ValueFrom: userSecret,
							},
							{
								Name:      "DATABASE_PASSWORD",
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

func (r *ReconcileCakephp) demoService(d *cakephpv1alpha1.Cakephp) *corev1.Service {
	labels := labels(d, "demo")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      demoServiceName(d),
			Namespace: d.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       5000,
				TargetPort: intstr.FromInt(8080),
				NodePort:   0,
			}},
		},
	}

	controllerutil.SetControllerReference(d, s, r.scheme)
	return s
}

func (r *ReconcileCakephp) demoRoute(d *cakephpv1alpha1.Cakephp) *routev1.Route {
	labels := labels(d, "demo")

	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      demoServiceName(d),
			Namespace: d.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: demoServiceName(d),
			},
		},
	}
	controllerutil.SetControllerReference(d, route, r.scheme)
	return route
}

func (r *ReconcileCakephp) updateDemoStatus(d *cakephpv1alpha1.Cakephp) error {
	d.Status.PodStatus = demoImage
	err := r.client.Status().Update(context.TODO(), d)
	return err
}

func (r *ReconcileCakephp) handleDemoChanges(d *cakephpv1alpha1.Cakephp) (*reconcile.Result, error) {
	found := &appsv1.Deployment{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      demoDeploymentName(d),
		Namespace: d.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &reconcile.Result{RequeueAfter: 5 * time.Second}, err
	}

	size := d.Spec.Size

	if size != *found.Spec.Replicas {
		found.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return &reconcile.Result{Requeue: true}, nil
	}

	title := d.Spec.Title
	evarIndex, ok := getEnvVarIndexByName("FIRST_LASTNAME", (*found).Spec.Template.Spec.Containers[0].Env) // 0
	if ok {
		existing := (*found).Spec.Template.Spec.Containers[0].Env[evarIndex].Value

		if title != existing {
			(*found).Spec.Template.Spec.Containers[0].Env[2].Value = title
			err = r.client.Update(context.TODO(), found)
			if err != nil {
				log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
				return &reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return &reconcile.Result{Requeue: true}, nil
		}

	}

	return nil, nil
}

func getEnvVarIndexByName(varname string, vars []corev1.EnvVar) (int, bool) {
	for i, val := range vars {
		if val.Name == varname {
			return i, true
		}
	}
	return -1, false
}
