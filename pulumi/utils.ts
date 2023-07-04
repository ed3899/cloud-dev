import * as fs from "fs";
import * as path from "path";
import * as Joi from "joi";
import * as R from "ramda";

interface jsonIP {
  publicIP: string;
}

export function getPublicIP(path: string) {
  try {
    const jsonData = fs.readFileSync(path, "utf-8");
    const jsonObject = JSON.parse(jsonData) as jsonIP;
    return jsonObject.publicIP;
  } catch (error) {
    console.error(
      `There was an error while reading the file "${path}" which includes your public ip. Your instance will have default broad ingress access. Here is the error: ${
        (error as Error).message
      } `,
    );
    return;
  }
}

type jsonBuild = {
  name: string;
  builder_type: string;
  build_time: number;
  files: null;
  artifact_id: string;
  packer_run_uuid: string;
  custom_data: {
    [key: string]: string;
  };
};

type jsonManifest = {
  builds: jsonBuild[];
  last_run_uuid: string;
};

export function getLastBuiltAMI_FromPackerManifest(): string {
  let data: jsonManifest;
  const manifest = "manifest.json";
  const filePath = path.join(__dirname, "..", "packer", manifest);
  try {
    const fileContent = fs.readFileSync(filePath, "utf-8");
    data = JSON.parse(fileContent);
  } catch (error) {
    throw new Error(
      `Error reading or parsing ${manifest} file: ${
        (error as Error).message
      }. Please make sure to create an AMI first`,
    );
  }

  const matchingBuild = data.builds.find((build) => build.packer_run_uuid === data.last_run_uuid);
  if (!matchingBuild) {
    throw new Error(
      "Matching build not found. Make sure not to modify the manifest Packer outputs once completing the build",
    );
  }

  const { artifact_id } = matchingBuild;
  if (!artifact_id) {
    throw new Error(
      "Empty artifact_id. Make sure not to modify the manifest Packer outputs once completing the build",
    );
  }

  return artifact_id.split(":")[1];
}

export function extractUserIds(input: string): string[] {
  // Extract the content within the square brackets
  const content = input.substring(1, input.length - 1);

  // Split the content using commas as delimiters
  const ids = content.split(",");

  // Remove leading/trailing whitespace and quotes from each ID
  const cleanedIds = ids.map((id) => id.trim().replace(/['"]/g, ""));

  return cleanedIds;
}

export const extractArrayFromString = (input: string): string[] => {
  return R.pipe(
    (str: string) => str.substring(1, str.length - 1),
    R.split(','),
    R.map(R.trim),
    R.map((id) => id.replace(/['"]/g, '')),
    R.uniq
  )(input);
};

interface ReqVars {
  AWS_ACCESS_KEY: string;
  AWS_SECRET_KEY: string;
  AWS_IAM_PROFILE: string;
  AWS_USER_IDS: string;
  AWS_AMI_NAME: string;
  AWS_INSTANCE_TYPE: string;
  AWS_REGION: string;
  AWS_EC2_AMI_NAME_FILTER: string;
  AWS_EC2_AMI_ROOT_DEVICE_TYPE: string;
  AWS_EC2_AMI_VIRTUALIZATION_TYPE: string;
  AWS_EC2_AMI_OWNERS: string;
  AWS_EC2_SSH_USERNAME: string;
  AWS_EC2_INSTANCE_SSH_KEY_NAME: string;
  PULUMI_PERSONAL_ACCESS_TOKEN: string;
}

interface OptVars {
  AWS_EC2_INSTANCE_VOLUME_TYPE: string;
  AWS_EC2_INSTANCE_VOLUME_SIZE: string;
}

interface MergedVars extends Partial<OptVars>, ReqVars {
}


function checkRequiredVariables(envVars: NodeJS.ProcessEnv): MergedVars {
  const reqVars: ReqVars = {
    AWS_ACCESS_KEY: "",
    AWS_SECRET_KEY: "",
    AWS_IAM_PROFILE: "",
    AWS_USER_IDS: "",
    AWS_AMI_NAME: "",
    AWS_INSTANCE_TYPE: "",
    AWS_REGION: "",
    AWS_EC2_AMI_NAME_FILTER: "",
    AWS_EC2_AMI_ROOT_DEVICE_TYPE: "",
    AWS_EC2_AMI_VIRTUALIZATION_TYPE: "",
    AWS_EC2_AMI_OWNERS: "",
    AWS_EC2_SSH_USERNAME: "",
    AWS_EC2_INSTANCE_SSH_KEY_NAME: "",
    PULUMI_PERSONAL_ACCESS_TOKEN: "",
  }

  const optVars: OptVars = {
    AWS_EC2_INSTANCE_VOLUME_TYPE: "",
    AWS_EC2_INSTANCE_VOLUME_SIZE: "",
  }

  const mergedVars: MergedVars = R.mergeLeft(reqVars, optVars)
  
  // ['alpha234234','asdfasdf234'] or ["alpha234234","asdfasdf234"]. Both as strings
  const stringArrayPattern = new RegExp(/^\[['"][a-zA-Z\d]+['"],['"][a-zA-Z\d]+['"]\]$/);

  const schema = Joi.object<MergedVars>({
    AWS_ACCESS_KEY: Joi.string().required().trim(),
    AWS_SECRET_KEY: Joi.string().required().trim(),
    AWS_IAM_PROFILE: Joi.string().trim(),
    AWS_USER_IDS: Joi.string().required().trim().pattern(stringArrayPattern),
    AWS_AMI_NAME: Joi.string().required().trim(),
    AWS_INSTANCE_TYPE: Joi.string().required(),
    AWS_REGION: Joi.string().required().trim(),
    AWS_EC2_AMI_NAME_FILTER: Joi.string().required().trim(),
    AWS_EC2_AMI_ROOT_DEVICE_TYPE: Joi.string().required().trim(),
    AWS_EC2_AMI_VIRTUALIZATION_TYPE: Joi.string().required().trim(),
    AWS_EC2_AMI_OWNERS: Joi.string().required().trim(),
    AWS_EC2_SSH_USERNAME: Joi.string().required().trim(),
    AWS_EC2_INSTANCE_SSH_KEY_NAME: Joi.string().required().trim(),
    PULUMI_PERSONAL_ACCESS_TOKEN: Joi.string().required().trim(),
    AWS_EC2_INSTANCE_VOLUME_TYPE: Joi.string().trim(),
    AWS_EC2_INSTANCE_VOLUME_SIZE: Joi.string().trim(),
  })

  const validationResult = schema.validate(mergedVars);
  if (validationResult.error) {
      console.error("Ensure no spaces around the variables")
      console.error("Arrays must be like ['alpha234234','asdfasdf234'] or [\"alpha234234\",\"asdfasdf234\"] with no spaces neither in each item and between separators")
      throw new Error(validationResult.error.message);
  }

  return validationResult.value;
}

interface CleanedMergedVars extends Partial<OptVars> {
  AWS_ACCESS_KEY: string;
  AWS_SECRET_KEY: string;
  AWS_IAM_PROFILE: string;
  AWS_USER_IDS: string[];
  AWS_AMI_NAME: string;
  AWS_INSTANCE_TYPE: string;
  AWS_REGION: string;
  AWS_EC2_AMI_NAME_FILTER: string;
  AWS_EC2_AMI_ROOT_DEVICE_TYPE: string;
  AWS_EC2_AMI_VIRTUALIZATION_TYPE: string;
  AWS_EC2_AMI_OWNERS: string[];
  AWS_EC2_SSH_USERNAME: string;
  AWS_EC2_INSTANCE_SSH_KEY_NAME: string;
  PULUMI_PERSONAL_ACCESS_TOKEN: string;
}

function cleanEnvVariables(env:MergedVars): CleanedMergedVars {
  const transformedEnv = R.evolve({
    // AWS_ACCESS_KEY: (v:string) => v.trim(),
    // AWS_SECRET_KEY: (v: string) => v.trim(),
    // AWS_IAM_PROFILE: (v: string) => v.trim(),
    AWS_USER_IDS: (v: string) => extractArrayFromString(v),
    // AWS_AMI_NAME: (v: string) => v.trim(),
    // AWS_INSTANCE_TYPE: (v: string) => v.trim(),
    // AWS_REGION: (v: string) => v.trim(),
    // AWS_EC2_AMI_NAME_FILTER: (v: string) => v.trim(),
    // AWS_EC2_AMI_ROOT_DEVICE_TYPE: (v: string) => v.trim(),
    // AWS_EC2_AMI_VIRTUALIZATION_TYPE: (v: string) => v.trim(),
    AWS_EC2_AMI_OWNERS: (v: string) => extractArrayFromString(v),
    // AWS_EC2_SSH_USERNAME: (v:string) => v.trim(),
    // AWS_EC2_INSTANCE_SSH_KEY_NAME: (v:string) => v.trim(),
    // PULUMI_PERSONAL_ACCESS_TOKEN: (v:string) => v.trim(),
    // AWS_EC2_INSTANCE_VOLUME_TYPE: (v:string) => v.trim(),
    // AWS_EC2_INSTANCE_VOLUME_SIZE: (v:string) => v.trim(),
  })<MergedVars>(env)

  return transformedEnv
}

export function getCleanEnvVars(env: NodeJS.ProcessEnv) {
  return R.pipe(
    checkRequiredVariables,
    cleanEnvVariables
  )(env)
}

export interface SSHConfig {
  host: string;
  hostname: string;
  identityFile: string;
  user: string;
  port?: number;
  userKnownHostsFile?: string;
  strictHostKeyChecking?: string;
  passwordAuthentication?: string;
  identitiesOnly?: string;
  logLevel?: string;
}

function validateSSH_Config(config: SSHConfig): SSHConfig | undefined {
  // Validate for absolute paths on either Windows or Unix and home paths on Unix.
  const pathsRegex = new RegExp(
    `^([A-Za-z]:|${path.posix.sep}|${path.win32.sep}|~[A-Za-z][A-Za-z0-9_-]*)`,
  );
  const sshSchema = Joi.object<SSHConfig>({
    host: Joi.string().required(),
    hostname: Joi.string().ip().required(),
    identityFile: Joi.string().regex(pathsRegex).required(),
    user: Joi.string().required(),
    port: Joi.number().integer(),
    userKnownHostsFile: Joi.string().regex(pathsRegex),
    strictHostKeyChecking: Joi.string(),
    passwordAuthentication: Joi.string(),
    identitiesOnly: Joi.string(),
    logLevel: Joi.string(),
  });

  const { error, value } = sshSchema.validate(config);

  if (error) {
    console.error("There was an error validating the ssh config:", error.details);
    return;
  } else {
    return value;
  }
}

type SSH_ConfigString = string;

function generateSSH_Config(config: SSHConfig | undefined): SSH_ConfigString {
  let result = "";
  if (config) {
    result += `Host ${config.host ?? "cloud-dev"}\n`;
    result += `  HostName ${config.hostname ?? "YOUR_EC2_PUBLIC_IP"}\n`;
    result += `  IdentityFile "${config.identityFile ?? "ABSOLUTE_PATH_TO_YOUR_SSH_FILE"}"\n`;
    result += `  User ${config.user ?? "AWS_EC2_INSTANCE_USERNAME"}\n`;
    result += `  User ${config.port ?? "22"}\n`;
    result += `  UserKnownHostsFile ${config.userKnownHostsFile ?? "/dev/null"}\n`;
    result += `  StrictHostKeyChecking ${config.strictHostKeyChecking ?? "no"}\n`;
    result += `  PasswordAuthentication ${config.passwordAuthentication ?? "no"}\n`;
    result += `  IdentitiesOnly ${config!.identitiesOnly ?? "yes"}\n`;
    result += `  LogLevel ${config!.logLevel ?? "FATAL"}\n`;
  }
  return result;
}

function exportSSH_Config(config: SSH_ConfigString) {
  const parentDir = path.join(__dirname, "..");
  const filePath = path.join(parentDir, "ssh_config");

  // Check if the config string is empty or null
  if (!config || config.trim() === "") {
    console.error("No ssh config file created");
    return;
  }

  fs.writeFile(filePath, config, (error) => {
    if (error) {
      console.error(
        `An error occurred while writing to the file "${filePath}". Here's the error:`,
        error,
      );
    } else {
      console.log(`The ssh config was created at ${filePath}`);
    }
  });
}

export function getSSH_KeyPath() {
  const keyFileName = process.env.AWS_EC2_INSTANCE_SSH_KEY_NAME!;
  const keyFilePath = path.join(__dirname, "..", keyFileName);
  return keyFilePath;
}

export const writeSSHConfig = R.pipe(validateSSH_Config, generateSSH_Config, exportSSH_Config);
