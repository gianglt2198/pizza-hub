# APPLICATION DEPLOYMENT

## SETUP ARGOCD APPLICATION
Prequisites:
- Tools:
```bash
# Install ArgoCD CLI
brew install argocd

# Install helm (optional, if using Helm charts)
brew install helm
```
- ArgoCD installed and configured in your Kubernetes cluster.
```bash
# Create ArgoCD namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready (takes ~2-3 minutes)
kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd

# Verify installation
kubectl get pods -n argocd
kubectl get svc -n argocd

# Install ApplicationSets
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/applicationset/v0.4.1/manifests/install.yaml  

# Verify ApplicationSets installed
kubectl get crd applicationsets.argoproj.io

# Get initial admin password
ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d)

echo "ArgoCD Admin Password: $ARGOCD_PASSWORD"

# Port forward to access ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443 &

# Access ArgoCD: https://localhost:8080
# Username: admin
# Password: $ARGOCD_PASSWORD
```
- Your GitHub repository containing the Kubernetes manifests for Pizza Hub.
- Access to the ArgoCD UI or CLI.
- Create an ArgoCD Application manifest (e.g., `pizza-hub-prod.yaml`) in the `infras/apps/argocd/applications` directory with the following content:

```yaml
Developer → Git Push → ArgoCD watches repo → Auto deploy to EKS
     ↓           ↓              ↓                    ↓
  Local     Git Repository   ArgoCD Server      K8s Cluster
  Changes   (Kubernetes     (GitOps Engine)    (Pizza Hub App)
            Manifests)
```


1. Update the `repoURL`, `targetRevision`, and `path` fields in the ArgoCD Application manifest to point to your GitHub repository and the correct path where your Kubernetes manifests are located.
2. Apply the ArgoCD Application manifest to create the application in ArgoCD:
```bash
# Deploy App of Apps
kubectl apply -f infras/apps/argocd/app-of-apps.yaml

# Deploy ApplicationSet  
kubectl apply -f infras/apps/argocd/appsets/pizza-hub-environments.yaml  

# Verify applications created  
kubectl get applications -n argocd  

# Check application status  
argocd app list  

# Sync applications (if not auto-sync)  
argocd app sync pizza-hub-dev  
argocd app sync pizza-hub-prod  
```
3. Monitor the deployment status in the ArgoCD UI or using the CLI. Once the application is synced successfully, your Pizza Hub application should be deployed to the EKS cluster.
```bash
# Port forward to ArgoCD UI  
kubectl port-forward svc/argocd-server -n argocd 8080:443 
# Open https://localhost:8080  
# Login với admin/$ARGOCD_PASSWORD   
```