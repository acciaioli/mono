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
```

### mono commands
```
▶ mono list
SERVICE                 STATUS  CHECKSUM        LOCAL CHECKSUM                          
examples/python-service diff                    29d93f3235df9a83f36c471085ba4f0b105ddec2
examples/go-service     diff                    78a79ed651e3d48b7a1a7346b022673e674abb8a

▶ mono checksum --service=examples/python-service
29d93f3235df9a83f36c471085ba4f0b105ddec2

▶ mono checksum --service=examples/python-service --pushed

# build a service
▶ mono build --service=examples/python-service
SERVICE                 ARTIFACT                                                                    
examples/python-service .builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip

# build sereral services
▶ mono build --service=examples/python-service --service=examples/go-service
SERVICE                 ARTIFACT                                                                    
examples/python-service .builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip
examples/go-service     .builds/examples/go-service/78a79ed651e3d48b7a1a7346b022673e674abb8a.zip    

# build all services which are not
▶ mono build
SERVICE                 ARTIFACT                                                                    
examples/python-service .builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip
examples/go-service     .builds/examples/go-service/78a79ed651e3d48b7a1a7346b022673e674abb8a.zip    
 

# push an artifact
▶ mono push --artifact=.builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip
ARTIFACT                                                                        STATUS          KEY                                                                     ERROR
.builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip    successful      /examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip        

# push all built artifacts
▶ mono push
ARTIFACT                                                                        STATUS          KEY                                                                     ERROR 
.builds/examples/go-service/78a79ed651e3d48b7a1a7346b022673e674abb8a.zip        successful      /examples/go-service/78a79ed651e3d48b7a1a7346b022673e674abb8a.zip            
.builds/examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip    successful      /examples/python-service/29d93f3235df9a83f36c471085ba4f0b105ddec2.zip        

# clean all built artifacts
▶ mono build --clean

# nothing to push now
▶ mono push

▶ mono checksum --service=examples/python-service --pushed
29d93f3235df9a83f36c471085ba4f0b105ddec2

▶ mono list
SERVICE                 STATUS  CHECKSUM                                        LOCAL CHECKSUM                          
examples/python-service ok      29d93f3235df9a83f36c471085ba4f0b105ddec2        29d93f3235df9a83f36c471085ba4f0b105ddec2
examples/go-service     ok      78a79ed651e3d48b7a1a7346b022673e674abb8a        78a79ed651e3d48b7a1a7346b022673e674abb8a

```

## To Do

- enhance `mono list` (prefix, contains, suffix filter)
- support for unit test
- support for docker images as artifacts
