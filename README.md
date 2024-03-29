# remotely
Scripts and Configurations for a remote development machine

# Coder

## Usage

The script will see if there is an instance running, if there is, it'll destroy it. Otherwise it'll create a new instance.

`go run coder.go`

## Configuration

Place a `.env` file in the root of this project. It must contain the following in order to work properly.

If you want to have visual studio code server, you must use `code-server-18-04` as your instance image name.

```
DO_TOKEN= # Your DigitalOcean API Token
REMOTELY_TAG= # A tag that groups these services together on Digital Ocean
REMOTELY_STORAGE_NAME= # Name of the storage for the droplet
REMOTELY_REGION= # Region where these services will be created
REMOTELY_INSTANCE_SIZE= # Size of the instance
REMOTELY_IMAGE_NAME= # Name of the image that the instance will be created
REMOTELY_INSTANCE_NAME= # Name of the droplet
REMOTELY_STORAGE_SIZE= # Size of your storage
```

An example of this might look something like;

```
DO_TOKEN=<your digital ocean token>
REMOTELY_TAG=coder
REMOTELY_STORAGE_NAME=coderStorage
REMOTELY_REGION=sfo2
REMOTELY_INSTANCE_SIZE=s-1vcpu-1gb
REMOTELY_IMAGE_NAME=code-server-18-04
REMOTELY_INSTANCE_NAME=vscoding
REMOTELY_STORAGE_SIZE=25
```
