# Kubernetes


## Ingress Nginx
Step 1: Set up nginx 
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.0/deploy/static/provider/baremetal/deploy.yaml
kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx
kubectl get svc -n ingress-nginx
```
Step 2: Create TLS/SSL
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout pizza-tls.key \
  -out pizza-tls.crt \
  -subj "/CN=pizza.local/O=pizza-hub" \
  -addext "subjectAltName=DNS:pizza.local,DNS:*.pizza.local"
```
Step 3: Attach to namespace 
```bash
kubectl create secret tls pizza-tls \
  --cert=pizza-tls.crt \
  --key=pizza-tls.key \
  -n pizza-hub
```
Verify 
```bash
kubectl get secret pizza-tls -n pizza-hub
```
Step 4: Run Ingress Following template 
```bash
kubectl apply -f ./default/ingress.yaml
```