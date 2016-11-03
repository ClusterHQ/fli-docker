# How To's for Examples

## redis-moby

Use three snapshots in the SAM

Example SAM

```
docker_app: docker-compose-app1.yml

flocker_hub:
    endpoint: https://data.flockerhub.clusterhq.com
    tokenfile: /root/fhut.txt

volumes:
    - name: redis-data
      snapshot: example-snapshot-1
      volumeset: docker-app-example
    - name: artifacts
      snapshot: example-snapshot-2
      volumeset: docker-app-example
    - name: /my/path
      snapshot: example-snapshot-3
      volumeset: docker-app-example
```

### Creating the volumesets for the example above.

```
$ fli init docker-app-example
c52d1b2f-c234-407e-94f3-89ac2f6f00a5

$ fli create docker-app-example vol1
/chq/c52d1b2f-c234-407e-94f3-89ac2f6f00a5/3f8eab4d-5ad3-4376-9783-16b407aaa396

$ fli create docker-app-example vol2
/chq/c52d1b2f-c234-407e-94f3-89ac2f6f00a5/f437a687-ddfa-48e6-8ce3-15d14fa9b588

$ fli create docker-app-example vol3
/chq/c52d1b2f-c234-407e-94f3-89ac2f6f00a5/7553a586-3b42-4edc-8aac-1bcae4358e13

$ touch /chq/c52d1b2f-c234-407e-94f3-89ac2f6f00a5/3f8eab4d-5ad3-4376-9783-16b407aaa396/hi1

$ touch /chq/c52d1b2f-c234-407e-94f3-89ac2f6f00a5/f437a687-ddfa-48e6-8ce3-15d14fa9b588/hi2

$ touch /chq/c52d1b2f-c234-407e-94f3-89ac2f6f00a5/7553a586-3b42-4edc-8aac-1bcae4358e13/hi3

$ fli snapshot docker-app-example:vol1 example-snapshot-1
ed3e5c0f-667c-47a3-956b-64e4e801e6c4

$ fli snapshot docker-app-example:vol2 example-snapshot-2
c76f9870-30a3-486a-bf53-b4b1a42b9219

$ fli snapshot docker-app-example:vol3 example-snapshot-3
613d8dd4-42d9-47fc-8623-d9ddbb5ec12e

$ fli sync docker-app-example

$ fli push docker-app-example
```