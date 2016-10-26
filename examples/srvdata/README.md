## How to try if you have no data.

First, setup Fli. See https://clusterhq.com/docs for how to do this.

Create the volumeset and snapshot this example expects
```
$ fli init srvdata-vs
$ fli create srvdata-vs tmpdata
$ touch /chq/pathto/volume/hello.txt
$ touch /chq/pathto/volume/hello2.txt
$ vi /chq/pathto/volume/hello2.txt
$ fli snapshot srvdata-vs:tmpdata tmp-data-snap
$ fli sync srvdata-vs
$ fli push srvdata-vs:tmp-data-snap
```

On the node with `fli-docker`,`docker-compose`, and `docker`

```
$ cd fli-docker/examples/srvdata
$ fli-docker run -f fli-manifest.yml -c -t /root/your.token
INFO[0005] [0/1] [srvdata]: Starting                    
INFO[0005] Building srvdata...                          
INFO[0006] [1/1] [srvdata]: Started 
```

See that data being served from your app
```
$ docker ps 
$ curl localhost:<PORT_OF_CONTAINER>
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN"><html>
<title>Directory listing for /</title>
<body>
<h2>Directory listing for /</h2>
<hr>
<ul>
<li><a href="hello.txt">hello.txt</a>
<li><a href="hello2.txt">hello2.txt</a>
</ul>
<hr>
</body>
</html>
```

Add more data to new volume, view it.

```
$ docker inspect -f "{{.Mounts}}" flicompose_srvdata_1
$ touch /chq/pathto/volume/helloagain.txt
$ curl localhost:<PORT_OF_CONTAINER>
$ <!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN"><html>
<title>Directory listing for /</title>
<body>
<h2>Directory listing for /</h2>
<hr>
<ul>
<li><a href="hello.txt">hello.txt</a>
<li><a href="hello2.txt">hello2.txt</a>
<li><a href="helloagain.txt">helloagain.txt</a>
</ul>
<hr>
</body>
</html>
```

Snapshot and push data back to FlockerHub
```
$ fli-docker snapshot -push
MESSAGE: Snapshotting and Pushing volumes to FlockerHub...
MESSAGE: Snapshotting and Pushing fli-ade3fdfb-1607-4d40-b501-fe2206969b00 from Volumeset srvdata-vs
```

