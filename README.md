# fli-docker

The `fli-docker` utility is designed to simplify the deployment of stateful applications inside Docker containers.

This is achieved through creation of a Flocker Hub Stateful Application Manifest (SAM) file (aka. "manifest"), which is used side by side with the Docker Compose file.
The SAM file is a YAML file that defines data volumes from [ClusterHQ](https://clusterhq.com)'s Flocker Hub,
synchronizes data snapshots locally, and maps them to Docker volumes in the Docker Compose file.

## Usage

To utilize the ClusterHQ `fli-docker` utility, examine the following command line arguments.

```
fli-docker --help
Usage of fli-docker:
  -c  [OPTIONAL] if flag is present, fli-docker will start the compose services
  -e string
      [OPTIONAL] Flocker Hub endpoint, optionally set it in the manifest YAML
  -f string
      [OPTIONAL] Stateful application manifest file (default "manifest.yml")
  -p string
      [OPTIONAL] project name for compose if using -c (default "fli-compose")
  -t string
      [OPTIONAL] Flocker Hub user token, optionally set it in the manifest YAML
  -verbose
      [OPTIONAL] verbose logging
```

### Example

You can use the example here in this repository. Follow the below instructions.

#### Install `fli-docker`

- Install `fli` (TODO)
- Install `docker` and `docker-compose`, see [here](https://docs.docker.com/compose/install/)
- Install the fli-docker binary. (TODO)

#### Run the example
```
$ git clone https://github.com/ClusterHQ/fli-docker/

$ cd fli-docker/examples/redis-moby

$ fli-docker -f fli-manifest.yml -c -p myproject
MESSAGE: 2016/10/03 21:55:19 main.go:83: Parsing the fli manifest...
MESSAGE: 2016/10/03 21:55:19 main.go:144: Pulling FlockerHub volumes...
MESSAGE: 2016/10/03 21:55:24 main.go:149: Creating volumes from snapshots...
MESSAGE: 2016/10/03 21:55:24 main.go:161: Mapping new volumes in compose file...
INFO[0005] [0/2] [redis]: Starting                      
INFO[0005] [1/2] [redis]: Started                       
INFO[0005] [1/2] [web]: Starting                        
INFO[0006] [2/2] [web]: Started                         
[root@ip-10-0-172-60 redis-moby]# docker ps
CONTAINER ID        IMAGE                          COMMAND                  CREATED             STATUS              PORTS                NAMES
dd97db0fc133        clusterhq/moby-counter:dcsea   "node index.js"          6 seconds ago       Up 5 seconds        0.0.0.0:80->80/tcp   myproject_web_1
5d535f5e1d55        redis:latest                   "docker-entrypoint.sh"   6 seconds ago       Up 5 seconds        6379/tcp             myproject_redis_1


$ $ docker inspect -f "{{.Mounts}}" myproject_web_1
[{ /chq/aff85bcb-3e2d-44b4-a458-9a4d7f030795 /myapp/artifacts  rw true rprivate}]
$ docker inspect -f "{{.Mounts}}" myproject_redis_1
[{ /chq/4083f9c2-3d8c-475e-bcab-06eefd49f60b /data  rw true rprivate}
 { /chq/33a74a69-ef86-40bb-afe6-8d82a2128dd7 /tmp/path  rw true rprivate}]
```

Optionally you dont have to specify `-c` or  `-p` so you can start the compose app yourself. Without `-c`
`fli-docker` will just modify the docker compose file and let you manage bring the services up.

```
$ fli-docker -f fli-manifest.yml
MESSAGE: 2016/10/03 22:06:21 main.go:83: Parsing the fli manifest...
MESSAGE: 2016/10/03 22:06:21 main.go:144: Pulling FlockerHub volumes...
MESSAGE: 2016/10/03 22:06:26 main.go:149: Creating volumes from snapshots...
MESSAGE: 2016/10/03 22:06:26 main.go:161: Mapping new volumes in compose file...

$ cat docker-compose-app1.yml
version: '2'
services:
  web:
    image: clusterhq/moby-counter:dcsea
    environment:
       - "USE_REDIS_HOST=redis"
    links:
      - redis
    ports:
      - "80:80"
    volumes:
      - /chq/aff85bcb-3e2d-44b4-a458-9a4d7f030795:/myapp/artifacts/
  redis:
    image: redis:latest
    volumes:
       - '/chq/4083f9c2-3d8c-475e-bcab-06eefd49f60b:/data'
       - /chq/33a74a69-ef86-40bb-afe6-8d82a2128dd7:/tmp/path

$ docker-compose -f docker-compose-app1.yml up -d
Pulling web (clusterhq/moby-counter:dcsea)...
dcsea: Pulling from clusterhq/moby-counter
a3ed95caeb02: Pull complete
93a86e942d51: Pull complete
faecfcc1d7ff: Pull complete
ddf3e3db435e: Pull complete
81a4604c1077: Pull complete
f8f4d4eabd85: Pull complete
Digest: sha256:e11ed56f5dad87ddef9865e758067ae6a182c234d9f10d1cf5c2a7d18a811eea
Status: Downloaded newer image for clusterhq/moby-counter:dcsea
Creating redismoby_redis_1
Creating redismoby_web_1
[root@ip-10-0-40-199 redis-moby]# docker-compose -f docker-compose-app1.yml ps
      Name                     Command               State         Ports        
-------------------------------------------------------------------------------
redismoby_redis_1   docker-entrypoint.sh redis ...   Up      6379/tcp           
redismoby_web_1     node index.js                    Up      0.0.0.0:80->80/tcp 

$ docker inspect -f "{{.Mounts}}" redismoby_web_1
[{ /chq/aff85bcb-3e2d-44b4-a458-9a4d7f030795 /myapp/artifacts  rw true rprivate}]
$ docker inspect -f "{{.Mounts}}" redismoby_redis_1
[{ /chq/4083f9c2-3d8c-475e-bcab-06eefd49f60b /data  rw true rprivate}
 { /chq/33a74a69-ef86-40bb-afe6-8d82a2128dd7 /tmp/path  rw true rprivate}]

```

## Stateful Application Manifest (SAM)

The Stateful Application Manifest (SAM) looks similar to a Docker Compose file, with a few key changes.

- The `volume_hub` node references an `endpoint` and a valid `auth_token`
- The volumes are defined by name, and each reference a `snapshot` and `volumeset`

The `fli-docker` utility takes a `docker-compose.yml` file as input, and translates
volumes in the Docker Compose file to Flocker Hub snapshots.

> Important Note: Right now, this only works with "named volumes" (see below)

[Compose File Reference Link](https://docs.docker.com/compose/compose-file/#/volumes-volume-driver)

```
# Named volume
  - datavolume1:/var/lib/mysql
  - 'datavolume2:/var/lib/mysql'
  - /my/path1:/tmp/path1
  - '/my/path2':/tmp/path2
```

An example of a Stateful App Manifest (SAM) YAML file could be `dev-manifest.yml` below. Notice, under the `volumes:` section of the 
manifest, that each named volume references a `volumeset` and a `snapshot`.
You can obtain these identifiers from the Flocker Hub user interface, or the `fli` command line utility.
Documentation about the Flocker Hub product itself can be found at [ClusterHQ Documentation](https://clusterhq.com).

```yaml
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

The Docker Compose file that the SAM file leverages would be:

```yaml
version: '2'
services:
  web:
    image: clusterhq/moby-counter
    environment:
       - "USE_REDIS_HOST=redis"
    links:
      - redis
    ports:
      - "80:80"
    volumes:
      - artifacts:/myapp/artifacts/
  redis:
    image: redis:latest
    volumes:
       - 'redis-data:/data'
       - /my/path:/tmp/path
```

In this case, the CLI commands above would perform the necessary `pull` and `create`
commands with fli and manipulate the docker-compose file so that when it is brought up
it can be brought up with your snapshots layed out in the manifest.

- `artifacts` would become snapshot : `11105373-b878-4433-8c8a-af6d684fe506`
- `redis-data` would become snapshot: `7c5c6dcb-8c65-4e68-ba60-262f8d5bf015`
- `/my/path` would become snapshot: `1670c1ff-c8be-4087-8eee-5a8598061a33`

### Notes

- You may run this from anywhere `docker-compose`, `docker` and `fli` are installed.
- Snapshots would need to be pushed to volumesets in ClusterHQ Flocker Hub using a manifest that references them, otherwise `pull` will fail.
