#!/usr/bin/env bash


# helm is installed into the K8S cluster specified
kubectl config current-context


# use --kube-context for different cluster
brew install kubernetes-helm

source <(helm completion bash)

# install help
helm init

helm version

# if running locally
export HELM_HOST=localhost:44134
tiller

# prometheus 2.0 has changed the rules/alerts format a new chart is (maybe)
# under development, but better check

helm install stable/prometheus --namespace monitoring

helm upgrade --values deployments/prometheus/alerts/alerts.yaml \
    --namespace monitoring funny-iguana stable/prometheus

export POD_NAME=$(kubectl get pods --namespace monitoring -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace monitoring port-forward $POD_NAME 9090


helm install stable/prometheus --namespace monitoring --values deployments/prometheus/alerts/alerts.yaml --debug --dry-run




# install from scripts
$ curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
$ chmod 700 get_helm.sh
$ ./get_helm.sh


# upgrade
helm init --upgrade
export TILLER_TAG=v2.0.0-alfa.1        # Or whatever version you want
kubectl --namespace=kube-system set image deployments/tiller-deploy tiller=gcr.io/kubernetes-helm/tiller:$TILLER_TAG

# delete
helm reset

# describe manifest
helm inspect values stable/mariadb


helm install stable/prometheus --name my-release -f values.yaml



# custom flags
#--home string               location of your Helm config. Overrides $HELM_HOME (default "/Users/luigi/.helm")
#      --host string               address of Tiller. Overrides $HELM_HOST
#      --kube-context string       name of the kubeconfig context to use
#      --tiller-namespace string   namespace of Tiller (default "kube-system")
