# ssm-sync

This tool synchronizes files with SSM.
In case key exists in SSM it will download and overwrite file on disk.
In case key doesn't exist in SSM it will upload file to SSM.

## Setup

Tool uses standard AWS profile mechanism - it should fallback to server role in case no profile / credential variables are provided.

## Example usage

ssm-sync -ssm ${ENVIRONMENT}/ssh/${APPLICATION}/${KEY_NAME} -file ${KEY_PATH}/${KEY_NAME} -env ${ENVIRONMENT}
