#!/bin/bash

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
  #docker push localhost:5000/microsoft/powershell:ubuntu14.04
  #docker push localhost:5000/microsoft/powershell:ubuntu16.04
  #docker push localhost:5000/microsoft/powershell:latest

  ### Take a snapshot of the Fli Volume
  snapshot_id=`fli snapshot ${volumeset_id}:${volume_name} ${snapshot_name} | tr -d '\r'`
  echo "$(tput setaf 6)Finished snapshotting the data volume ${volume_dir}. Snapshot ID is: ${snapshot_id}"

  ### Sync the VolumeSet with FlockerHub
  fli sync ${volumeset_id}
  echo "$(tput setaf 6)Synchronized VolumeSet (${volumeset_id}) with FlockerHub"

  ### Push the snapshot up to FlockerHub
  fli push ${volumeset_id}:${snapshot_id}
  echo "$(tput setaf 6)Finished pushing snapshot (${snapshot_id}), in VolumeSet (${volumeset_id}) to FlockerHub$(tput setaf 7)"
}

PrepSnapshot
