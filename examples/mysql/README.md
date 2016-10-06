# How To's for Examples

## mysql

First, create some snapshots to use with a SAM

```
$ dpcli set volumehub "http://<FlockerHub_URL>"
$ dpcli set tokenfile /root/vhut.txt 
$ dpcli create volumeset -d "a volumeset for fli-docker"  docker-app-example
$ dpcli create volume -v docker-app-example first_volume

// snapshot an empty volume for MySQL to use.
$ dpcli create snapshot -V <volume> -b fli-docker-branch snapshotOf_first_volume

$ dpcli sync volumeset docker-app-example
$ dpcli push volumeset docker-app-example
```

Use these three snapshots in the SAM

Example SAM

```
docker_app: docker-compose.yml

flocker_hub:
    endpoint: http://<ip|dnsname>:<port>
    tokenfile: /root/vhut.txt

volumes:
    - name: mysql-data
      snapshot: e6bfe755-6423-48cb-bf22-d9e4b799c305
      volumeset: docker-app-example
```
