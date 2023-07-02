import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as dotenv from "dotenv";
// Load environment variables
dotenv.config();
import {
  getPublicIP,
  runChecks,
  getAMI_ID,
  extractUserIds,
  writeSSHConfig,
  getSSH_KeyPath,
} from "./utils";

runChecks();

const PUBLIC_IP = getPublicIP("./publicIP.json");

const vpc = new aws.ec2.Vpc("cloud-dev-vpc", {
  cidrBlock: "10.0.0.0/16",
  enableDnsHostnames: true,
  enableDnsSupport: true,
});

const internetGateway = new aws.ec2.InternetGateway(
  "cloud-dev-internet-gateway",
  {
    vpcId: vpc.id,
  },
  { dependsOn: vpc },
);

// Ensure deployment in a valid AZ
const availableZone = aws.getAvailabilityZones({
  state: "available",
}).then(available => available.names?.[0]);

const subnet = new aws.ec2.Subnet(
  "cloud-dev-subnet",
  {
    vpcId: vpc.id,
    cidrBlock: "10.0.0.0/24",
    mapPublicIpOnLaunch: true,
    availabilityZone: availableZone
  },
  { dependsOn: internetGateway },
);

const routeTable = new aws.ec2.RouteTable(
  "cloud-dev-route-table",
  {
    vpcId: vpc.id,
    routes: [
      {
        cidrBlock: "0.0.0.0/0",
        gatewayId: internetGateway.id,
      },
    ],
  },
  { dependsOn: [vpc, internetGateway] },
);

const routeTableAssociation = new aws.ec2.RouteTableAssociation(
  "cloud-dev-route-table-association",
  {
    routeTableId: routeTable.id,
    subnetId: subnet.id,
  },
  {
    dependsOn: [routeTable, subnet],
  },
);

const securityGroup = new aws.ec2.SecurityGroup(
  "cloud-dev-security-group",
  {
    vpcId: vpc.id,
    ingress: [
      {
        protocol: "tcp",
        fromPort: 22,
        toPort: 22,
        cidrBlocks: [`${PUBLIC_IP ?? "0.0.0.0"}/${PUBLIC_IP ? "32" : "0"}`],
      },
    ],
    egress: [
      {
        protocol: "-1",
        fromPort: 0,
        toPort: 0,
        cidrBlocks: ["0.0.0.0/0"],
      },
    ],
  },
  { dependsOn: [vpc, routeTableAssociation] },
);

const ami = pulumi.output(
  aws.ec2.getAmi({
    mostRecent: true,
    owners: extractUserIds(process.env.AWS_USER_IDS!),
    tags: {
      Environment: "development",
      Builder: "packer",
    },
    filters: [
      {
        name: "image-id",
        values: [getAMI_ID()],
      },
      {
        name: "name",
        values: [process.env.AWS_AMI_NAME!],
      },
      {
        name: "root-device-type",
        values: [process.env.AWS_EC2_AMI_ROOT_DEVICE_TYPE!],
      },
      {
        name: "virtualization-type",
        values: [process.env.AWS_EC2_AMI_VIRTUALIZATION_TYPE!],
      },
    ],
  }),
);

const instance = new aws.ec2.Instance(
  "cloud-dev-ec2-instance",
  {
    instanceType: process.env.AWS_INSTANCE_TYPE!,
    ami: ami.id,
    vpcSecurityGroupIds: [securityGroup.id],
    subnetId: subnet.id,
    availabilityZone: availableZone,
  },
  {
    dependsOn: [vpc, subnet, securityGroup],
  },
);

// Post deployment
export const ec2PublicIp = instance.publicIp;
// Display ssh command
export const sshCommand = pulumi.interpolate`You can now ssh into your instance with this command: 'ssh -i ${process.env.AWS_EC2_INSTANCE_SSH_KEY_NAME} ${process.env.AWS_EC2_INSTANCE_USERNAME}@${ec2PublicIp}'`;
// Write ssh config
ec2PublicIp.apply((ip) => {
  writeSSHConfig({
    host: process.env.AWS_AMI_NAME!,
    hostname: ip,
    identityFile: getSSH_KeyPath(),
    user: process.env.AWS_EC2_INSTANCE_USERNAME!,
  });
});
