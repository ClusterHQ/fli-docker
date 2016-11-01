# How To's for Examples

## mysql

Use snapshots in the SAM

Example SAM

```
docker_app: docker-compose.yml

flocker_hub:
    endpoint: https://flockerhub.clusterhq.com
    tokenfile: /root/fh.token

volumes:
    - name: mysql-data
      snapshot: e6bfe755-6423-48cb-bf22-d9e4b799c305
      volumeset: docker-app-example
```
