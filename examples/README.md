# How To's for Examples

## redis-moby

First, create some snapshots to use with a SAM

```
$ dpcli set volumehub "http://<FlockerHub_URL>"
$ dpcli set tokenfile /root/vhut.txt 
$ dpcli create volumeset -d "a volumeset for fli-docker"  docker-app-example
$ dpcli create volume -v docker-app-example first_volume

$ touch /chq/<volume>/one.txt
$ dpcli create snapshot -V <volume> -b fli-docker-branch snapshotOf_first_volume

$ touch /chq/<volume>/two.txt
$ dpcli create snapshot -V <volume> -b fli-docker-branch -a snap=2 snapshotOf_first_volume_2

$ touch /chq/<volume>/three.txt
$ dpcli create snapshot -V <volume> -b fli-docker-branch -a snap=3 snapshotOf_first_volume_3

$ dpcli sync volumeset docker-app-example
$ dpcli push volumeset docker-app-example
```

Use these three snapshots in the SAM

Example SAM

```
docker_app: docker-compose-app1.yml

flocker_hub:
    endpoint: http://<ip|dnsname>:<port>
    tokenfile: /root/vhut.txt

volumes:
    - name: redis-data
      snapshot: snapshotOf_first_volume
      volumeset: docker-app-example
    - name: artifacts
      snapshot: snapshotOf_first_volume_2
      volumeset: docker-app-example
    - name: /my/path
      snapshot: snapshotOf_first_volume_3
      volumeset: docker-app-example
```
