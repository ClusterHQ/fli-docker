#!/bin/bash

set -e

function fli () {
  local zpool_name='denverimaging'
  docker run --rm -it --privileged -v /etc/hosts:/etc/hosts -v /root:/root -v /${zpool_name}:/${zpool_name}:shared -v /var/log/fli:/var/log/fli -v /lib/modules:/lib/modules clusterhq/fli "$@"
}

### Set a unique VOlumeSet and Volume name
volumeset_name='flidocker-example'
volume_name='dockerregistry'
snapshot_name='docker-images-1'
flidocker_path='/tmp/fli-docker'
flockerhub_token_path='/root/vhut.txt'

if test ! -e $flockerhub_token_path
then
  echo "FlockerHub token file does not exist: ${flockerhub_token_path}"
  exit 10
fi 

function PrepSnapshot () {
  ### Create a VolumeSet and Volume
  volumeset_id=`fli init ${volumeset_name} | tr -d '\r'`
  volume_dir=`fli create ${volumeset_id} ${volume_name} | tr -d '\r'`
  echo "$(tput setaf 6)Created VolumeSet ${volumeset_id} and Volume ${volume_dir}$(tput setaf 7)"

  ### Start a private, temporary registry instance
  ### This will simply be used to stage some data onto a Fli Volume, so we can snapshot it
  docker run --detach -v ${volume_dir}:/var/lib/registry -p 5000:5000 --name registry-temp registry

  ### Download Docker images locally
  docker pull microsoft/powershell:centos7
  docker pull microsoft/powershell:ubuntu14.04
  docker pull microsoft/powershell:ubuntu16.04
  docker pull microsoft/powershell:latest

  docker tag microsoft/powershell:centos7 localhost:5000/microsoft/powershell:centos7
  docker tag microsoft/powershell:ubuntu14.04 localhost:5000/microsoft/powershell:ubuntu14.04
  docker tag microsoft/powershell:ubuntu16.04 localhost:5000/microsoft/powershell:ubuntu16.04
  docker tag microsoft/powershell:latest localhost:5000/microsoft/powershell:latest

  ### Push images up to the private registry instance
  docker push localhost:5000/microsoft/powershell:centos7
  docker push localhost:5000/microsoft/powershell:ubuntu14.04
  #docker push localhost:5000/microsoft/powershell:ubuntu16.04
  #docker push localhost:5000/microsoft/powershell:latest

  ### Stop and remove the container
  docker rm -f registry-temp

  ### Take a snapshot of the Fli Volume
  snapshot_id=`fli snapshot ${volumeset_id}:${volume_name} ${snapshot_name} | tr -d '\r'`
  echo "Finished snapshotting the data volume ${volume_dir}. Snapshot ID is: ${snapshot_id}"
  
  ### Sync the VolumeSet with FlockerHub
  fli sync ${volumeset_id}
  echo "Synchronized VolumeSet (${volumeset_id}) with FlockerHub"
  
  ### Push the snapshot up to FlockerHub
  fli push ${volumeset_id}:${snapshot_id}
  echo "Finished pushing snapshot (${snapshot_id}), in VolumeSet (${volumeset_id}) to FlockerHub"
}

PrepSnapshot

### Invoke Fli-Docker
echo "Invoking Fli-Docker ..."
${flidocker_path} run -e https://ui.dev.voluminous.io -f fli-manifest.yml -c -t /root/vhut.txt
