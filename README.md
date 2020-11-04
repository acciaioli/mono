# mono

```
▶ mono
monorepo management cli

Usage:
  mono [command]

Available Commands:
  build       Builds artifact for a service
  checksum    Computes/Fetches a service checksum
  help        Help about any command
  list        Lists all services under the current directory
  push        Pushes a service artifact to the cloud

Flags:
  -h, --help      help for mono
  -v, --version   version for mono

Use "mono [command] --help" for more information about a command.
```


## Get Started

```
▶ export AWS_REGION=my-aws-region
▶ export AWS_PROFILE=my-aws-profile

▶ mono list
examples/go-service
examples/python-service

▶ mono checksum --service=examples/python-service
29d93f3235df9a83f36c471085ba4f0b105ddec2

▶ mono checksum --service=examples/python-service --pushed --bucket=s3://my-aws-s3-bucket


▶ mono build --service=examples/python-service
.tmp/builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip
                                                                                                                                                                                                            7m ⚑  
▶ mono push --artifact=.tmp/builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip --bucket=s3://my-aws-s3-bucket  
/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip
                                                                                                                                                                                                           7m ⚑  
▶ mono checksum --service=examples/python-service --pushed --bucket=s3://my-aws-s3-bucket                                       
29d93f3235df9a83f36c471085ba4f0b105ddec2
```

## To Do

- enhance `mono list`
- `build --updated` support (build if current checksum differs from last pushed one)
- `push --all` support (push all artifacts under the builds directory)
- github actions
- support for unit test
- support for docker images as artifacts