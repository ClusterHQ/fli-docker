# fs3todocker

## Usage

```
fs3docker --help
Usage of fs3docker:
  -c string
    	Options to pass to Docker Compose such as 'up -d' (default "up")
  -f string
    	Stateful application manifest file (default "manifest.yml")
  -t string
    	Flocker Hub user token
  -u string
    	Flocker Hub username
  -v string
    	Flocker Hub endpoint
```

### Example

```
$ fs3docker -c "up -d" -f dev-manifest.yml -t cf4add5b3be133f51de4044b9affd79edeca51d3 -u wallnerryan -v http://10.0.0.2:8080
```

## Stateful App Manifest

The stateful app manifest takes a docker-compose file and translates
volumes in the compose file to Flocker Hub snapshots.

Example could be `dev-manifest.yml` below.
```
docker_app:
    - docker-compose-app1.yml

volume_hub:
    endpoint: http://<ip>:<port>

volumes:
   redis-data:
      snapshot: be4b53d2-a8cf-443f-a672-139b281acf8f
      volumeset: e2799be7-cb75-4686-8707-e66083da3260
   artifacts:
      snapshot: 02d474fa-ab81-4bcb-8a61-a04214896b67
      volumeset: e2799be7-cb75-4686-8707-e66083da3260
```

The compose file in this manifest to be used would be

```
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
       - ‘redis-data:/data'
```

In this case, the CLI commands above would perform the necessary `pull` and `create`
commands with fs3 and manipulate the docker-compose file so that when it is brought up
it can be brought up with your snapshots layed out in the manifest.

- `artifacts` would become snapshot : `02d474fa-ab81-4bcb-8a61-a04214896b67`
- `redis-data` would become snapshot: `be4b53d2-a8cf-443f-a672-139b281acf8f`

### Notes

- You may run this from anywhere `docker-compose`, `docker` and `fs3` are installed.
- Snapshots would need to be pushed to volumesets prior to running this.
