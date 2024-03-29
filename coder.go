package main

import (
    "context"
    "github.com/digitalocean/godo"
    "crypto/rand"
    "encoding/base64"
    "math"
    "github.com/joho/godotenv"
    "time"
    "strconv"
    "golang.org/x/oauth2"
    "os"
    "log"
    "fmt"
)

type TokenSource struct {
    AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
    token := &oauth2.Token {
        AccessToken: t.AccessToken,
    }
    return token, nil
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func createDroplet(client *godo.Client, password string) *godo.Droplet {
    tag := getFlag("REMOTELY_TAG")
    tags := []string{tag}

    echoString := fmt.Sprintf("  - echo %s > /etc/code-server/pass\n", password)

    userData := "#cloud-config\n\n" +
    "runcmd:\n" +
    echoString +
    "  - reboot\n"

    createRequest := &godo.DropletCreateRequest {
        Name: getFlag("REMOTELY_INSTANCE_NAME"),
        Region: getFlag("REMOTELY_REGION"),
        Size: getFlag("REMOTELY_INSTANCE_SIZE"),
        Image: godo.DropletCreateImage {
            Slug: getFlag("REMOTELY_IMAGE_NAME"),
        },
        Tags: tags,
        PrivateNetworking: true,
        UserData: userData,
    }

    ctx := context.TODO()

    newDroplet, _, err := client.Droplets.Create(ctx, createRequest)

    if err != nil {
        fmt.Printf("Something bad happened: %s\n\n", err)
    }

    fmt.Printf("New Droplet Created: %s\n\n", newDroplet.Name)
    return newDroplet
}

func deleteDroplet(client *godo.Client) {
    ctx := context.TODO()
    tag := getFlag("REMOTELY_TAG")

    _, err := client.Droplets.DeleteByTag(ctx, tag)

    if err != nil {
        printError(err)
    } else {
        fmt.Printf("Droplet has been deleted\n")
    }
}

func getDroplet(client *godo.Client) (*godo.Droplet, bool) {
    ctx := context.TODO()
    tag := getFlag("REMOTELY_TAG")
    opt := &godo.ListOptions{}

    droplets, _, _ := client.Droplets.ListByTag(ctx, tag, opt)

    if len(droplets) > 0 {
        return &droplets[0], true
    } else {
        return nil, false
    }
}

func getBlockStorage(client *godo.Client) (*godo.Volume, bool) {
    ctx := context.TODO()

    fmt.Printf("Retrieving volume...\n")
    volume := retrieveStorage(client, getFlag("REMOTELY_STORAGE_NAME"))

    if volume != nil {
        return volume, true
    }

    fmt.Printf("Creating volume...\n")
    tag := getFlag("REMOTELY_TAG")
    tags := []string{tag}

    createRequest := &godo.VolumeCreateRequest {
        Region: getFlag("REMOTELY_REGION"),
        Name: getFlag("REMOTELY_STORAGE_NAME"),
        Description: "Storage for coder",
        Tags: tags,
        SizeGigaBytes: getIntFlag("REMOTELY_STORAGE_SIZE"),
    }

    newVolume, _, volErr := client.Storage.CreateVolume(ctx, createRequest)

    if volErr == nil {
        return newVolume, true
    }

    printError(volErr)
    return nil, false
}

func retrieveStorage(client *godo.Client, storageName string) *godo.Volume {
    ctx := context.TODO()

    params := &godo.ListVolumeParams {
        Name: storageName,
    }

    volumes, _, err := client.Storage.ListVolumes(ctx, params)

    if err != nil {
        printError(err)
        return nil
    }

    if len(volumes) > 0 {
        return &volumes[0]
    } else {
        return nil
    }
}

func getFlag(name string) string {
    flag, exists := os.LookupEnv(name)

    if exists {
        return flag
    } else {
        fmt.Printf("Please define %s in your .env file.\n\nExiting now.\n", name)
        os.Exit(1)
        return ""
    }
}

func getIntFlag(name string) int64 {
    flag := getFlag(name)

    val, err := strconv.ParseInt(flag, 10, 64)

    if err == nil {
        return val
    } else {
        fmt.Printf("Please define %s in your .env file as an integer.\n\nExiting now.\n", name)
        os.Exit(1)
        return 0
    }
}

func printError(err error) {
    fmt.Printf("Something went wrong: %s\n\n", err)
}

func attachBlockStorage(client *godo.Client, volumeID string, dropletID int) bool {
    ctx := context.TODO()

    _, _, err := client.StorageActions.Attach(ctx, volumeID, dropletID)

    if err != nil {
        fmt.Printf("Something bad happened with attaching the storage: %s\n\n", err)
        return false
    } else {
        fmt.Printf("Everything is good to go!\n\n")
        return true
    }
}

func randomBase64String(l int) string {
    buff := make([]byte, int(math.Round(float64(l)/float64(1.3333333))))
    rand.Read(buff)
    str := base64.RawURLEncoding.EncodeToString(buff)
    return str[:l]
}

func main() {
    pat, exists := os.LookupEnv("DO_TOKEN")

    if !exists {
        fmt.Printf("Can't find token")
        return
    }

    tokenSource := &TokenSource {
        AccessToken: pat,
    }

    oauthClient := oauth2.NewClient(context.Background(), tokenSource)
    client := godo.NewClient(oauthClient)
    _, dExists := getDroplet(client)

    if dExists {
        deleteDroplet(client)
    } else {
        password := randomBase64String(12)
        droplet := createDroplet(client, password)
        volume, volumeExists := getBlockStorage(client)

        if volumeExists {
            retries := 0
            for retries < 3 {
                fmt.Printf("Waiting for droplet to finish initialising...\n")
                time.Sleep(20 * time.Second)
                fmt.Printf("Attempting to attach block storage to droplet.\n")
                success := attachBlockStorage(client, volume.ID, droplet.ID)

                if success {
                    break
                } else {
                    retries += 1
                }
            }

            if retries == 3 {
                fmt.Printf("Please check your digital ocean account. Something went wrong.")
            } else {
                droplet, _ := getDroplet(client)
                ip, err := droplet.PublicIPv4()
                if err == nil {
                    fmt.Printf("Public IP is: %s\n", ip)
                    fmt.Printf("Password: %s\n", password)
                } else {
                    printError(err)
                }
            }
        } else {
            fmt.Printf("Something went bad with the volumes")
        }
    }
}
