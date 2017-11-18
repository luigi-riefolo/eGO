#####################
#       minikube    #
#####################
minikube start

minikube stop

minikube dashboard

# run the service
minikube service hello-node

curl $(minikube service hello-minikube --url)


# check catalog
curl http://localhost:5000/v2/_catalog

curl  http://localhost:5000/v2/alfa/tags/list

curl  http://localhost:5000/v2/alfa/manifests/latest

DOCKER_API_VERSION= DOCKER_TLS_VERIFY= DOCKER_CERT_PATH= DOCKER_HOST= docker ps


###################
#       K8s       #
###################

kubectl get events

kubectl config view

kubectl get nodes

kubectl get all

kubectl replace --force -f src/alfa/conf/k8s-service.yaml

# Determine service IP
kubectl get service $SERVICE --output='jsonpath="{.spec.ports[0].nodePort}"'

# start minikube
minikube start --vm-driver=xhyve

# check the contexts/clusters
cat ~/.kube/config

# set minikube context
kubectl config use-context minikube

# create a deployment
kubectl run service-srv --image=localhost:5000/ --port=8080

# show deployments
kubectl get deployments

# show service
kubectl get services

# expose a service
kubectl expose deployment hello-node --type=LoadBalancer

# Check the kubectl configuration
kubectl cluster-info

kubectl get pod

kubectl get pods --all-namespaces


# build and run service
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -v -a -installsuffix cgo -o main server.go
docker build -t alf:v1 .
kubectl create -f k8s-deployment.yaml
kubectl create -f k8s-service.yaml
#kubectl run alfa --image=alfa:v1 --port=8080
#kubectl expose deployment service --type=LoadBalancer
minikube service alfa-gateway --url

kubectl logs alfa
kubectl logs -f alfa-97854d8ff-9ss8w

kubectl delete service alfa
kubectl delete deployment alfa


kubectl create configmap prometheus-server-conf \
    --from-file=kube/prometheus/config-v2.yaml \
    -o yaml \
    --dry-run | kubectl replace -f -

# check rollout status
kubectl rollout status deployment/nginx-deployment

# check replica set
kubectl get rs

# check pods with labels
kubectl get pods --show-labels
kubectl get pods -l app=alfa -o wide


# Note: A Deployment’s rollout is triggered if and only if the Deployment’s pod template (that is, .spec.template) is changed, for example if the labels or container images of the template are updated. Other updates, such as scaling the Deployment, do not trigger a rollout.
# Suppose that we now want to update the nginx Pods to use the nginx:1.9.1 image instead of the nginx:1.7.9 image.
kubectl set image deployment/nginx-deployment nginx=nginx:1.9.1


docker build -t alfa .

docker tag alfa localhost:5000/alfa

docker push localhost:5000/registry-demo

docker ps -a -f status=exited

########################################
#docker run -p 127.0.0.1:5432:5432
#docker run -d --hostname localhost
##########################################


# run the image
#docker run --publish 6060:7000 --name $CONTAINER --rm $CONTAINER \
#docker run \
#    -e SERVICE_NAME="$SERVICE_NAME" -e BUILD_ID="${BUILD_ID}" -e DEPLOYMENT_TYPE="service" \
#    --hostname=${SERVICE_ARTIFACT} \
#    --name $SERVICE_NAME  \
#    --network "alfa_default" \
#    --publish 6060:7000 \
#    --rm ${IMAGE_NAME}

#--dns 8.8.8.8,8.8.4.4
#--env-file file
#--health-cmd string              Command to run to check health
#--health-retries 1
#--health-timeout 5m
#--ip string                      IPv4 address (e.g., 172.30.100.104)
#--ip6 string                     IPv6 address (e.g., 2001:db8::33


#####################
#       Docker      #
#####################

# Remove all containers
docker rm --force $(docker ps -a -q)

# Remove all images
docker rmi --force $(docker images -q)

# from container ping host
ping docker.for.mac.localhost



# AM

# Run project
docker-compose -f docker-compose.am-test.yml up --remove-orphans --force-recreate --abort-on-container-exit

# Modify an image and run it locally
# add to Dockerfile RUN npm link gulp
npm install
npm start
docker build . -t autonomous.car2go.com/am-system-api

docker-compose -f docker-compose.am-test.yml up am-test-amapi nginx



# Swagger
docker run -p 80:8080 -e SWAGGER_JSON=/data/alfa.swagger.json -v $GOPATH/src/github.com/luigi-riefolo/eGO/src/alfa/pb/definitions:/data swaggerapi/swagger-ui





# Prometheus
docker run -p 9090:9090 prom/prometheus

# provide a config
docker run -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
DOCKER_API_VERSION= DOCKER_TLS_VERIFY= DOCKER_CERT_PATH= DOCKER_HOST= docker run -p 9090:9090 prom/prometheus

