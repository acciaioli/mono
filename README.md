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

### AWS
```
▶ export AWS_REGION=my-aws-region
▶ export AWS_PROFILE=my-aws-profile
▶ export MONO_ARTIFACT_BUCKET=my-aws-s3-bucket

▶ mono list
Service: examples/python-service
Status: ok
Local Checksum: 29d93f3235df9a83f36c471085ba4f0b105ddec2
Pushed Checksum: 29d93f3235df9a83f36c471085ba4f0b105ddec2

Service: examples/go-service
Status: diff
Local Checksum: 2742f0cd2d31fb3d04dd1d14482941026de9f56d
Pushed Checksum: 1ce479560b36cd620b659c7bafe970cc307536dd

▶ mono checksum --service=examples/python-service
29d93f3235df9a83f36c471085ba4f0b105ddec2

▶ mono checksum --service=examples/python-service --pushed


▶ mono build --service=examples/python-service
.tmp/builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip

▶ mono push --artifact=.tmp/builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip  
/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip

▶ mono checksum --service=examples/python-service --pushed
29d93f3235df9a83f36c471085ba4f0b105ddec2
```

## To Do

- enhance `mono list` (prefix, contains, suffix filter)
- `build --all` support (build if current checksum differs from last pushed one)
- `push --all` support (push all artifacts under the builds directory)
- github actions
- support for unit test
- support for docker images as artifacts
