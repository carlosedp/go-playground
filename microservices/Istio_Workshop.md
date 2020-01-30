# Istio and Tracing workshop

Create a minikube instance

```bash
minikube start \
    --vm-driver=virtualbox \
    --network-plugin=cni \
    --enable-default-cni \
    --container-runtime=cri-o \
    --bootstrapper=kubeadm \
    --cpus=4 \
    --memory=6144 \
    --cache-images=true

minikube addons enable dashboard
minikube addons enable ingress
minikube addons enable metrics-server

# Create dashboard ingress
IP=$(minikube ip)
cat <<EOF | kubectl apply -f -
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: dashboard
  namespace: kubernetes-dashboard
spec:
  rules:
  - host: dashboard.${IP}.nip.io
    http:
      paths:
      - backend:
          serviceName: kubernetes-dashboard
          servicePort: 80
        path: /
EOF
```

Download and deploy Istio

```bash
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.3.4 sh -

cd istio-1.3.4

for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done

kubectl apply -f install/kubernetes/istio-demo.yaml

kubectl patch deployment istio-tracing -n istio-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"jaeger","image":"docker.io/jaegertracing/all-in-one:1.15"}]}}}}'

kubectl patch deployment kiali -n istio-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"kiali","image":"quay.io/kiali/kiali:v1.9"}]}}}}'

kubectl get svc -n istio-system
kubectl get pods -n istio-system

kubectl label namespace default istio-injection=enabled

# Create namespace for demo apps
kubectl create namespace istio-demo
kns istio-demo
kubectl label namespace istio-demo istio-injection=enabled

# Manual injection
istioctl kube-inject -f <your-app-spec>.yaml | kubectl apply -f -
```

Create ingress for Istio dashboards

```yaml
IP=$(minikube ip)
cat <<EOF | kubectl apply -f -
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: istio-grafana
  namespace: istio-system
spec:
  rules:
  - host: istio-grafana.${IP}.nip.io
    http:
      paths:
      - backend:
          serviceName: grafana
          servicePort: 3000
        path: /
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: istio-kiali
  namespace: istio-system
spec:
  rules:
  - host: istio-kiali.${IP}.nip.io
    http:
      paths:
      - backend:
          serviceName: kiali
          servicePort: 20001
        path: /
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: istio-tracing
  namespace: istio-system
spec:
  rules:
  - host: istio-tracing.${IP}.nip.io
    http:
      paths:
      - backend:
          serviceName: tracing
          servicePort: 80
        path: /
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: istio-prometheus
  namespace: istio-system
spec:
  rules:
  - host: istio-prometheus.${IP}.nip.io
    http:
      paths:
      - backend:
          serviceName: prometheus
          servicePort: 9090
        path: /
---
EOF

k get ingresses --all-namespaces
```

## HTTPbin Sample

Deploy HTTPbin sample.

```bash
kubectl apply -f samples/httpbin/httpbin.yaml

IP=$(minikube ip)
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpbin-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "httpbin.${IP}.nip.io"
EOF

kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpbin
spec:
  hosts:
  - "httpbin.${IP}.nip.io"
  gateways:
  - httpbin-gateway
  http:
  - match:
    - uri:
        prefix: /status
    - uri:
        prefix: /delay
    route:
    - destination:
        port:
          number: 8000
        host: httpbin
EOF

export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export INGRESS_HOST=$(minikube ip)

curl -I -H Host:httpbin.${IP}.nip.io http://$INGRESS_HOST:$INGRESS_PORT/status/200
```

## Bookinfo Sample

```bash
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml

# Test
kubectl exec -it $(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}') -c ratings -- curl productpage:9080/productpage | grep -o "<title>.*</title>"

kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml

kubectl get gateway

# With nodeport
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export INGRESS_HOST=$(minikube ip)

# With `minikube tunnel`

export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')

export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

echo http://${GATEWAY_URL}/productpage

curl -s http://${GATEWAY_URL}/productpage | grep -o "<title>.*</title>"

open http://${GATEWAY_URL}/productpage

kubectl apply -f samples/bookinfo/networking/destination-rule-all.yaml
```

## Uninstall

```bash
# HTTPbin
kubectl delete gateway httpbin-gateway
kubectl delete virtualservice httpbin
kubectl delete --ignore-not-found=true -f samples/httpbin/httpbin.yaml
kubectl delete -f samples/httpbin/httpbin.yaml

#Bookinfo
samples/bookinfo/platform/kube/cleanup.sh

#Istio
kubectl delete ingress.extensions/istio-grafana -n istio-system
kubectl delete ingress.extensions/istio-kiali -n istio-system
kubectl delete ingress.extensions/istio-tracing -n istio-system
kubectl delete ingress.extensions/istio-prometheus -n istio-system
kubectl delete -f install/kubernetes/istio-demo.yaml
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl delete -f $i; done


```