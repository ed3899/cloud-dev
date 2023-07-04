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

export const extractArrayFromString = (input: string): string[] => {
  return R.pipe(
    // Remove brackets
    (str: string) => str.substring(1, str.length - 1),
    R.split(","),
    R.map(R.trim),
    // Remove quotes
    R.map((id) => id.replace(/['"]/g, "")),
    R.uniq,
  )(input);
};

interface EnvVars {
  AWS_USER_IDS: string;
  AWS_AMI_NAME: string;
  AWS_INSTANCE_TYPE: string;
  AWS_EC2_AMI_NAME_FILTER: string;
  AWS_EC2_AMI_ROOT_DEVICE_TYPE: string;
  AWS_EC2_AMI_VIRTUALIZATION_TYPE: string;
  AWS_EC2_AMI_OWNERS: string;
  AWS_EC2_SSH_USERNAME: string;
  AWS_EC2_INSTANCE_SSH_KEY_NAME: string;
  AWS_EC2_INSTANCE_VOLUME_TYPE?: string;
  AWS_EC2_INSTANCE_VOLUME_SIZE?: string;
  PACKER_AMI_TO_LAUNCH_FROM?: string;
}

function checkRequiredVariables(env: NodeJS.ProcessEnv): EnvVars {
  const envVars: EnvVars = {
    AWS_USER_IDS: "",
    AWS_AMI_NAME: "",
    AWS_INSTANCE_TYPE: "",
    AWS_EC2_AMI_NAME_FILTER: "",
    AWS_EC2_AMI_ROOT_DEVICE_TYPE: "",
    AWS_EC2_AMI_VIRTUALIZATION_TYPE: "",
    AWS_EC2_AMI_OWNERS: "",
    AWS_EC2_SSH_USERNAME: "",
    AWS_EC2_INSTANCE_SSH_KEY_NAME: "",
    AWS_EC2_INSTANCE_VOLUME_TYPE: undefined,
    AWS_EC2_INSTANCE_VOLUME_SIZE: undefined,
    PACKER_AMI_TO_LAUNCH_FROM: undefined,
  };

  // Iterate and fill from process.env
  const fillEnvVars = (envVars: EnvVars) => {
    R.forEachObjIndexed((_, key) => {
      if (process.env[key]) {
        envVars = R.assoc(key, process.env[key], envVars) as EnvVars;
      }
    }, envVars);

    return envVars;
  };

  const filledEnvVars = fillEnvVars(envVars);

  // ['alpha234234','asdfasdf234'] or ["alpha234234","asdfasdf234"]. Both as strings
  const stringArrayPattern = new RegExp(/^\[['"][a-zA-Z\d]+['"](,\s*['"][a-zA-Z\d]+['"])*\]$/);

  const schema = Joi.object<EnvVars>({
    AWS_USER_IDS: Joi.string().required().trim().pattern(stringArrayPattern),
    AWS_AMI_NAME: Joi.string().required().trim(),
    AWS_INSTANCE_TYPE: Joi.string().required(),
    AWS_EC2_AMI_NAME_FILTER: Joi.string().required().trim(),
    AWS_EC2_AMI_ROOT_DEVICE_TYPE: Joi.string().required().trim(),
    AWS_EC2_AMI_VIRTUALIZATION_TYPE: Joi.string().required().trim(),
    AWS_EC2_AMI_OWNERS: Joi.string().required().trim().pattern(stringArrayPattern),
    AWS_EC2_SSH_USERNAME: Joi.string().required().trim(),
    AWS_EC2_INSTANCE_SSH_KEY_NAME: Joi.string().required().trim(),
    AWS_EC2_INSTANCE_VOLUME_TYPE: Joi.string().trim(),
    AWS_EC2_INSTANCE_VOLUME_SIZE: Joi.string().trim(),
    PACKER_AMI_TO_LAUNCH_FROM: Joi.string().trim(),
  });

  const validationResult = schema.validate(filledEnvVars);
  if (validationResult.error) {
    console.error("Ensure no spaces around the variables");
    console.error(
      "Arrays must be like ['alpha234234','asdfasdf234'] or [\"alpha234234\",\"asdfasdf234\"] with no spaces neither in each item and between separators",
    );
    throw new Error(validationResult.error.message);
  }

  return validationResult.value;
}

interface EvolvedOptVars {
  AWS_EC2_INSTANCE_VOLUME_TYPE: string;
  AWS_EC2_INSTANCE_VOLUME_SIZE: number;
  PACKER_AMI_TO_LAUNCH_FROM: string;
}

interface EvolvedVars extends Partial<EvolvedOptVars> {
  AWS_USER_IDS: string[];
  AWS_AMI_NAME: string;
  AWS_INSTANCE_TYPE: string;
  AWS_EC2_AMI_NAME_FILTER: string;
  AWS_EC2_AMI_ROOT_DEVICE_TYPE: string;
  AWS_EC2_AMI_VIRTUALIZATION_TYPE: string;
  AWS_EC2_AMI_OWNERS: string[];
  AWS_EC2_SSH_USERNAME: string;
  AWS_EC2_INSTANCE_SSH_KEY_NAME: string;
}

function evolveVars(env: EnvVars): EvolvedVars {
  const evolvedVars = R.evolve({
    AWS_USER_IDS: (v: string) => extractArrayFromString(v),
    AWS_EC2_INSTANCE_VOLUME_SIZE: (v: string | undefined) => (v ? parseInt(v, 10) : undefined),
    AWS_EC2_AMI_OWNERS: (v: string) => extractArrayFromString(v),
  })<EnvVars>(env);

  return evolvedVars;
}

export function getCleanEnvVars(env: NodeJS.ProcessEnv) {
  return R.pipe(checkRequiredVariables, evolveVars)(env);
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
