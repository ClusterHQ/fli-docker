# flitodock

The `flitodock` utility is designed to simplify the deployment of stateful applications inside Docker containers.

This is achieved through creation of a Flocker Hub Stateful Application Manifest (SAM) file (aka. "manifest"), which essentially acts as a wrapper to a Docker Compose file.
The SAM file is a YAML file that defines data volumes from [ClusterHQ](https://clusterhq.com)'s Flocker Hub,
synchronizes data snapshots locally, and maps them to Docker volumes in the underlying Docker Compose file.

## Usage

To utilize the ClusterHQ `flitodock` utility, examine the following command line arguments.

```
flitodock --help
Usage of flitodock:
  -c --compose-arguments string
    	Options to pass to Docker Compose such as 'up -d' (default "up")
  -f --manifest string
    	Stateful application manifest file (default "manifest.yml")
  -t --token string
    	Flocker Hub user token
  -u --username string
    	Flocker Hub username
  -e --endpoint string
    	Flocker Hub endpoint
```

### Example

In the following example command, we are passing in a base Docker Compose YAML file, referencing a 

```
$ flitodock -c "up -d" -f dev-manifest.yml -t cf4add5b3be133f51de4044b9affd79edeca51d3 -u wallnerryan -e http://10.0.0.2:8080
```

## Stateful Application Manifest (SAM)

The Stateful Application Manifest (SAM) looks similar to a Docker Compose file, with a few key changes.

- The `volume_hub` node references an `endpoint` and a valid `auth_token`
- The volumes are defined by name, and each reference a `snapshot` and `volumeset`

The `flitodock` utility takes a `docker-compose.yml` file as input, and translates
volumes in the Docker Compose file to Flocker Hub snapshots.

An example of a Stateful App Manifest (SAM) YAML file could be `dev-manifest.yml` below. Notice, under the `volumes:` section of the 
manifest, that each named volume references a `volumeset` and a `snapshot`.
You can obtain these identifiers from the Flocker Hub user interface, or the `fs3` command line utility.
Documentation about the Flocker Hub product itself can be found at [ClusterHQ Documentation](https://clusterhq.com).

```yaml
docker_app:
    - docker-compose-app1.yml

volume_hub:
    endpoint: http://<ip|dnsname>:<port>
    auth_token: 021e3d0f-9ad3-49dc-8d0a-dbe96a0477dc

volumes:
   redis-data:
      snapshot: be4b53d2-a8cf-443f-a672-139b281acf8f
      volumeset: e2799be7-cb75-4686-8707-e66083da3260
   artifacts:
      snapshot: 02d474fa-ab81-4bcb-8a61-a04214896b67
      volumeset: e2799be7-cb75-4686-8707-e66083da3260
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
```

In this case, the CLI commands above would perform the necessary `pull` and `create`
commands with fs3 and manipulate the docker-compose file so that when it is brought up
it can be brought up with your snapshots layed out in the manifest.

- `artifacts` would become snapshot : `02d474fa-ab81-4bcb-8a61-a04214896b67`
- `redis-data` would become snapshot: `be4b53d2-a8cf-443f-a672-139b281acf8f`

### Notes

- You may run this from anywhere `docker-compose`, `docker` and `fs3` are installed.
- Snapshots would need to be pushed to volumesets in ClusterHQ Flocker Hub prior to running this.
