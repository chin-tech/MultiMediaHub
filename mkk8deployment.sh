#!/usr/bin/env bash

pushd k8s

NECESSARY_ENV_VARS=(
   "DB_PASSWORD"
   "PORT"
   )

# Check if necessary vars are set ()
for var in "${NECESSARY_ENV_VARS[@]}"; do
   [[ -z ${var} ]] &&  echo "[-] Error: $var is not set!" && exit 1
done

echo '[+] Environment variables are set!'



## Create Deployment Directory for minikube ##
mkdir -p ../deployments
echo "[+] Creating deployments folder and files"

for f in *; do
   envsubst  < $f > ../deployments/$f
done

popd

echo 'Deployment directory created with relevant configuration'
echo 'RUN: minikube start; kubectl apply -f deployments'
