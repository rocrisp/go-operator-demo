go-operator-demo
=======================

Need Knowledge of openshift, kubernetes, and go

Requirements:

Install [Golang](https://golang.org/doc/install)

Install [operator-sdk](https://sdk.operatorframework.io/docs/install-operator-sdk/)

Install [Openshift Container Platform 4.5](https://docs.openshift.com/container-platform/4.5/welcome/index.html)

Run Operator sdk locally
```bash
git clone https://github.com/rocrisp/go-operator-demo.git
cd go-operator-demo
oc new-project sandbox
oc apply -f deploy/crds/cakephp.example.com_cakephps_crd.yaml 
OPERATOR_NAME=cakephp operator-sdk run --local
```
![](gif/runoperatorsdk.gif)

In another terminal
```bash
cd go-operator-demo
oc get pods
oc apply -f deploy/crds/cakephp.example.com_v1alpha1_cakephp_cr.yaml
oc get all
```

![](gif/createoperator.gif)

![](gif/Operator.gif)






