package cakephp

import (
	"context"

	routev1 "github.com/openshift/api/route/v1"
	cakephpv1alpha1 "github.com/rocrisp/go-operator-demo/pkg/apis/cakephp/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileCakephp) ensureDeployment(request reconcile.Request,
	instance *cakephpv1alpha1.Cakephp,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

//See if route exist, create one if it doesn't already exist
func (r *ReconcileCakephp) ensureRoute(request reconcile.Request,
	instance *cakephpv1alpha1.Cakephp,
	rs *routev1.Route,
) (*reconcile.Result, error) {
	found := &routev1.Route{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      rs.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the route
		log.Info("Creating a new Route", "Route.Namespace", rs.Namespace, "Route.Name", rs.Name)
		err = r.client.Create(context.TODO(), rs)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Route", "Route.Namespace", rs.Namespace, "Route.Name", rs.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Route")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileCakephp) ensureService(request reconcile.Request,
	instance *cakephpv1alpha1.Cakephp,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileCakephp) ensureSecret(request reconcile.Request,
	instance *cakephpv1alpha1.Cakephp,
	s *corev1.Secret,
) (*reconcile.Result, error) {
	found := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		log.Info("Creating a new secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the secret not existing
		log.Error(err, "Failed to get Secret")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func labels(v *cakephpv1alpha1.Cakephp, tier string) map[string]string {
	return map[string]string{
		"app":         "Demo",
		"demosite_cr": v.Name,
		"tier":        tier,
	}
}
