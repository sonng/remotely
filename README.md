# remotely
Scripts and Configurations for a remote development machine


## Configuration

Place a `.env` file in the root of this project. It must contain the following in order to work properly.

```
DO_TOKEN= # Your DigitalOcean API Token
REMOTELY_TAG= # A tag that groups these services together on Digital Ocean
REMOTELY_STORAGE_NAME= # Name of the storage for the droplet
REMOTELY_REGION= # Region where these services will be created
REMOTELY_INSTANCE_SIZE= # Size of the instance
REMOTELY_IMAGE_NAME= # Name of the image that the instance will be created
REMOTELY_INSTANCE_NAME= # Name of the droplet
```

An example of this might look something like;

```
DO_TOKEN=<your digital ocean token>
REMOTELY_TAG=coder
REMOTELY_STORAGE_NAME=coderStorage
REMOTELY_REGION=sfo2REMOTELY_INSTANCE_SIZE=s-1vcpu-1gb
REMOTELY_IMAGE_NAME=ubuntu-14-04-x64
REMOTELY_INSTANCE_NAME=vscoding
```
