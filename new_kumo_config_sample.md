```.yaml
Cloud: aws

AWS:
  AccessKeyId: xxxxxxxx
  SecretAccessKey: xxxxxxxxxxxxxxx
  IamProfile: xxxxxxx
  UserIds:
    - "xxxxxxxxxxxx"
  Region: us-west-2
  EC2:
    Instance:
      Type: t2.micro
    Volume:
      Type: gp2 #optional
      Size: 8 #optional

AMI:
  Base:
    Filter: ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20230516
    RootDeviceType: ebs
    VirtualizationType: hvm
    Owners:
      - "099720109477"
    User: ubuntu
  Name: "test"
  User: dev
  Home: home
  Password: test12345
  Tools:
    - always
  IdToBeUsed: "" #optionalx

Git:
  Username: xxxxxxxxx
  Email: xxxxxxxxxxxxx

GitHub:
  PersonalAccessTokenClassic: xxxxxxxxxxxxx #optional

```