# How To's for Examples

## redis-moby

First, create some snapshots to use with a SAM

```
$ dpcli set volumehub "http://<FlockerHub_URL>"
$ dpcli set tokenfile /root/vhut.txt 
$ dpcli create volumeset -n docker-app-example
$ dpcli create volume -v docker-app-example

$ touch /chq/<volume>/one.txt
$ dpcli create snapshot -V <volume> -b sam-example -a purpose=sam

$ touch /chq/<volume>/two.txt
$ dpcli create snapshot -V <volume> -b sam-example -a purpose=sam,snap=two

$ touch /chq/<volume>/three.txt
$ dpcli create snapshot -V <volume> -b sam-example -a purpose=sam,snap=three

$ dpcli sync volumeset docker-app-example
$ dpcli push volumeset docker-app-example
```

Use these three snapshots in the SAM

> example output

```
$ dpcli show snapshot -v ca7f73e8-3665-4559-9414-36f89bbb80ec
BRANCH      ID                                   ATTRIBUTES
sam-example 11105373-b878-4433-8c8a-af6d684fe506 purpose=sam,snap=three
            7c5c6dcb-8c65-4e68-ba60-262f8d5bf015 purpose=sam,snap=two
            1670c1ff-c8be-4087-8eee-5a8598061a33 purpose=sam
```

Example SAM

```
docker_app: docker-compose-app1.yml

flocker_hub:
    endpoint: http://<ip|dnsname>:<port>
    tokenfile: /root/vhut.txt

volumes:
    - name: redis-data
      snapshot: 11105373-b878-4433-8c8a-af6d684fe506
      volumeset: docker-app-example
    - name: artifacts
      snapshot: 7c5c6dcb-8c65-4e68-ba60-262f8d5bf015
      volumeset: docker-app-example
    - name: /my/path
      snapshot: 1670c1ff-c8be-4087-8eee-5a8598061a33
      volumeset: docker-app-example
```

## Example 2

TODO

## Example 3

TODO