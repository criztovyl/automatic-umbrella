in case you use minikube with containerd it might not properly configure insecure registry (here `minikube:5000` from registry addon).

to fix that manually, run the following commands in `minikube ssh`:

```
sudo test -d /etc/containerd/certs.d/minikube:5000 || mkdir $_
sudo tee /etc/containerd/certs.d/minikube:5000/hosts.toml >/dev/null <<EOT
server = "http://minikube:5000"

[host."http://minikube:5000"]
  capabilities = ["pull", "resolve", "push"]
  skip_verify = true
EOT
```
