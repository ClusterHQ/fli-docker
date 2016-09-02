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
