# Pulumi-Go-K8s-ESO-Demo
Pulumi Golang Kubernetes External Secrets Operator Demo

### Prerequisits
- [DigitalOcean account - referal link for $200 credit](https://m.do.co/c/d3723dfedb27)
- CloudFlare account
- FQDN (Domain Name)
- CA issued TLS cert configured for the FQDN
- Python 3.14.0
- Ansible - 2.18.9
- Golang - go1.25.5
- Unix like system capable of running Ansible - WSL2 works this is due to Ansible
- [Pulumi](https://www.pulumi.com/docs/get-started/download-install/) 

### Structure

- esocluster: DigitalOcean K8s Cluster
    - eso-manifests: K8s manifests to sync credentials from HashiCorp Vault
    - fastapi-container: Python FastAPI Rest API example using a synced secret
- vaultvm: HashiCorp Vault VM

### Setup

1. Configure the Pulumi.dev.yaml config files
   - esocluster/Pulumi.dev.yaml
   - vaultvm/Pulumi.dev.yaml 
2. `$ pulumi up`
   - esocluster
   - vaultvm 

3. Configure kubectl
  - `$ pulumi stack output --show-secrets doKubeConfig > ~/.kube/config`

4. Configure credentials in HashiCorp Vault
    - `$ ssh login root@<FQDN>`
    - Reference the unseal key and root token in the file located at `  vaultvm:vaultTokenOutFilePath`
    - Unseal Vault 
      - `$ vault operator unseal`
    - Login with the root token
      - `$ vault login`
    - Enable the secret/ kv path
      - `vault secrets enable -path=secret/ kv`
    - Load the secret into Vault
      - `$ vault kv put secret/pydemo DEMO_SECRET=password`
    - Create a Vault policy that can read the secret - see pydemo-policy.hcl
      - `$ vim pydemo-policy.hcl`
      - `$ cat pydemo-policy.hcl | vault policy write pydemo-policy -`
    - Create the Vault token for the newly created Vault policy
      - `$ vault token create -policy="pydemo-policy"`
5. Create a K8s namespace to sync external secrets into
   - `$ kubectl create namespace example`
6. Create a K8s Secret for the newly created Vault policy token
   - Note: You can also create a token for the root Vault token which will have full access, however, this is not a good practice.
   - `$ kubectl create secret generic pydemo-do-vault-token --from-literal=token=<POLICY TOKEN> -n example`
6. Create a K8s Secret or ConfigMap for the TLS CA Bundle
   - `$ kubectl create secret generic cabdlsec --from-file=ca.crt -n example`
7. Configure and apply the ESO Secret Store manifest under eso-manifests
   - `$ kubectl apply -f pydemo-do-secret-store.yaml`
   - `$ kubectl describe SecretStore pydemo-do-secret-store`
   - Note: SecretStores are namespace specific
8. Configure and apply the ESO External Secret under eso/manifests
   - `$ kubectl apply -f pydemo-do-ext-secret-namespace.yaml`
   - `$ kubectl describe ExternalSecret pydemo-do-ext-secret-namespace`
9. Deploy the Python FastAPI Demo API under fastapi-container
   - `$ kubectl apply -f deployment.yaml`
   - `$ kubectl apply -f service.yaml`
   - `$ kubectl apply -f ingress.yaml` - see known issues
10. Query the REST GET /secret API endpoint
    - `$ curl http://<FQDN>/secret` - through ingress
    - `$ curl http://127.0.0.1:8080/secret` - after exec-ing into the pod
11. Clean up the infrastructure
12. `$ pulumi destroy`
    - esocluster
    - vaultvm 


### Known Issues
1. The Nginx ingress controller may not properly route to the pyesodemo-ingress
   2. Results in 503 errors
2. Pulumi may get hung up on destroying the ESO and / or Nginx Helm chart resources and require those resources to be deleted from the state manually.
   - `$ pulumi stack -u -l`
   - `$ pulumi state delete urn:pulumi:dev::esocluster::kubernetes:helm.sh/v3:Release::external-secrets`