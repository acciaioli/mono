# mono

![CI-CD](https://github.com/acciaioli/mono/workflows/CI-CD/badge.svg)

```
▶ mono
monorepo management cli

Usage:
  mono [command]

Available Commands:
  build       Builds artifact for a service
  checksum    Computes services checksums
  help        Help about any command
  list        Lists all services under the current directory
  push        Pushes a service artifact to the cloud
```

## Install

##### Linux
```
▶ sudo wget https://github.com/acciaioli/mono/releases/latest/download/mono.linux-amd64 -O /usr/local/bin/mono
▶ sudo chmod +x /usr/local/bin/mono
```

## Get Started

### AWS
```
▶ export AWS_REGION=my-aws-region
▶ export AWS_PROFILE=my-aws-profile
▶ export MONO_ARTIFACT_BUCKET=my-aws-s3-bucket
```

### mono usage

```
▶ mono --version
mono version snapshot-juan-dc0649b

▶ mono list
SERVICE                	DIFF	VERSION	CHECKSUM	LOCAL CHECKSUM                           
examples/go-service    	true	0      	-       	5215256745a800ab8909097d230933f42ebe249f	
examples/python-service	true	0      	-       	fc09aaad1dc4cfd8f3165d9f5f9af2f394675664	

▶ mono build --clean

▶ mono build
SERVICE                	ARTIFACT                                                                     
examples/go-service    	.builds/examples/go-service/5215256745a800ab8909097d230933f42ebe249f.zip    	
examples/python-service	.builds/examples/python-service/fc09aaad1dc4cfd8f3165d9f5f9af2f394675664.zip	

▶ mono push
ARTIFACT                                                                    	STATUS    	KEY                                                                    	ERROR 
.builds/examples/python-service/fc09aaad1dc4cfd8f3165d9f5f9af2f394675664.zip	successful	examples/python-service/v1.fc09aaad1dc4cfd8f3165d9f5f9af2f394675664.zip	-    	
.builds/examples/go-service/5215256745a800ab8909097d230933f42ebe249f.zip    	successful	examples/go-service/v1.5215256745a800ab8909097d230933f42ebe249f.zip    	-    	

▶ mono list
SERVICE                	DIFF 	VERSION	CHECKSUM                                	LOCAL CHECKSUM                           
examples/go-service    	false	1      	5215256745a800ab8909097d230933f42ebe249f	5215256745a800ab8909097d230933f42ebe249f	
examples/python-service	false	1      	fc09aaad1dc4cfd8f3165d9f5f9af2f394675664	fc09aaad1dc4cfd8f3165d9f5f9af2f394675664	

▶ echo foo bar > examples/python-service/a-new-file.txt

▶ mono list
SERVICE                	DIFF 	VERSION	CHECKSUM                                	LOCAL CHECKSUM                           
examples/python-service	true 	1      	fc09aaad1dc4cfd8f3165d9f5f9af2f394675664	486bff14a41a94c82b37109447398ce1ab225643	
examples/go-service    	false	1      	5215256745a800ab8909097d230933f42ebe249f	5215256745a800ab8909097d230933f42ebe249f	

▶ mono build --clean

▶ mono build
SERVICE                	ARTIFACT                                                                     
examples/python-service	.builds/examples/python-service/486bff14a41a94c82b37109447398ce1ab225643.zip	

▶ mono push
ARTIFACT                                                                    	STATUS    	KEY                                                                    	ERROR 
.builds/examples/python-service/486bff14a41a94c82b37109447398ce1ab225643.zip	successful	examples/python-service/v2.486bff14a41a94c82b37109447398ce1ab225643.zip	-    	

▶ mono list
SERVICE                	DIFF 	VERSION	CHECKSUM                                	LOCAL CHECKSUM                           
examples/go-service    	false	1      	5215256745a800ab8909097d230933f42ebe249f	5215256745a800ab8909097d230933f42ebe249f	
examples/python-service	false	2      	486bff14a41a94c82b37109447398ce1ab225643	486bff14a41a94c82b37109447398ce1ab225643	
```

## To Do

- enhance `mono list` (prefix/contains filter)
- support for docker images as artifacts
