#!/bin/bash

set -e

function fli () {
  local zpool_name='denverimaging'
  docker run --rm -it --privileged -v /etc/hosts:/etc/hosts -v /root:/root -v /${zpool_name}:/${zpool_name}:shared -v /var/log/fli:/var/log/fli -v /lib/modules:/lib/modules quay.io/clusterhq_prod/fli:c6a5deac3bb68b93341c8accfdde66fd7d13fc1f "$@"
}

### Set a unique VOlumeSet and Volume name
volumeset_name='flidocker-example-registry-6'
volume_name='dockerregistry'
snapshot_name='docker-images'
flidocker_path='/tmp/fli-docker'
flockerhub_token_path='/root/vhut.txt'

if test ! -e $flockerhub_token_path
then
  echo "FlockerHub token file does not exist: ${flockerhub_token_path}"
  exit 10
fi

### Invoke Fli-Docker
echo "$(tput setaf 6)Invoking Fli-Docker ..."
${flidocker_path} run -verbose -f fli-manifest.yml -c 
