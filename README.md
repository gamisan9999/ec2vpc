# ec2vpc
get VPC ID from EC2 InstanceID

# install & execute

- run the EC2

```
local $ gox
local $ scp ec2vpc_<os>_<arch> <target EC2>:~
server $ ./ec2vpc_linux_amd64 --i `curl -s http://169.254.169.254/latest/meta-data/instance-id`
vpc-XXXXXXXX
```

- run the local

```
orenomac$ go run main.go --instance-id <insance_id> --region ap-northeast-1
vpc-XXXXXXXX
```
