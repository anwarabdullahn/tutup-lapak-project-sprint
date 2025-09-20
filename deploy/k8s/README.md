# Kubernetes Deployment Notes

## Core Patterns
1. Each microservice has a Deployment + Service.
2. Database migrations run via separate Job (or Helm hook) before Deployments roll out.
3. Secrets store DATABASE_URL, JWT_SECRET, INTERNAL_SECRET.
4. Service-to-service calls use internal Cluster DNS names: `<service>.<namespace>.svc.cluster.local`.

## Migration Strategies
Choose one:
- Dedicated Job per service (recommended): Apply Job manifests before Deployments. Jobs are idempotent if migrations are.
- Helm hooks: Annotate Job with `helm.sh/hook: pre-install,pre-upgrade` so it runs automatically.
- InitContainer: Runs migrations before main container starts (can block rolling updates if long).

## Helm Chart Structure (optional)
```
charts/
  tutup-lapak/
    Chart.yaml
    values.yaml
    templates/
      namespace.yaml
      secrets.yaml
      auth-deployment.yaml
      product-deployment.yaml
      backend-infra-deployment.yaml
      _helpers.tpl
      jobs/
        auth-migration-job.yaml  # annotated with: helm.sh/hook: pre-install,pre-upgrade
```
Example hook annotations inside Job metadata:
```
annotations:
  helm.sh/hook: pre-install,pre-upgrade
  helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded
```

## Applying
```
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/secrets-config.yaml
# Run migrations first
kubectl apply -f deploy/k8s/auth-migration-job.yaml
# Then services
kubectl apply -f deploy/k8s/auth-deployment.yaml
kubectl apply -f deploy/k8s/product-deployment.yaml
kubectl apply -f deploy/k8s/backend-infra-deployment.yaml
```

Monitor Jobs:
```
kubectl get jobs -n tutup-lapak
kubectl logs job/<job-name> -n tutup-lapak
```

Rollback by creating down migration Job or using migrate CLI externally.

## Next Steps
- Add remaining services (profile, purchase) similarly.
- Add Postgres (StatefulSet) or reference managed cloud DB.
- Introduce Ingress (Nginx / ALB) exposing backend-infra.
- Add HorizontalPodAutoscaler for variable load.
- Enforce NetworkPolicies limiting DB access to services.
- Add PodDisruptionBudgets to maintain availability.
- Integrate CI pipeline: build images, push, helm upgrade --install.
- Observability: Prometheus scraping, Grafana dashboards, centralized logging.
- Security: image scanning, non-root (already), read-only root FS, secrets encryption.
- Autoscaling signals: add metrics endpoint.
