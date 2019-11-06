package main

import (
    "context"
    "github.com/digitalocean/godo"
    "github.com/joho/godotenv"
    "golang.org/x/oauth2"
    "os"
    "log"
    "fmt"
)

const (
    tag = "coder"
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

func createDroplet(client *godo.Client) {
    dropletName := "vscoding"
    tags := []string{tag}

    createRequest := &godo.DropletCreateRequest {
        Name: dropletName,
        Region: "sfo2",
        Size: "s-1vcpu-1gb",
        Image: godo.DropletCreateImage {
            Slug: "ubuntu-14-04-x64",
        },
        Tags: tags,
        PrivateNetworking: true,
    }

    ctx := context.TODO()

    newDroplet, _, err := client.Droplets.Create(ctx, createRequest)

    if err != nil {
        fmt.Printf("Something bad happened: %s\n\n", err)
    }

   fmt.Printf("New Droplet Created: %s\n\n", newDroplet.Name)

   ip, ipError := newDroplet.PrivateIPv4()
   if ipError == nil {
       fmt.Printf("It's Private IP is: %s\n", ip)
   } else {
       fmt.Printf("Something bad happened: %s\n\n", ipError)
   }
}

func deleteDroplet(client *godo.Client) {
    ctx := context.TODO()

    _, err := client.Droplets.DeleteByTag(ctx, tag)

    if err != nil {
        fmt.Printf("Something bad happened: %s\n\n", err)
    } else {
        fmt.Printf("Droplet has been deleted")
    }
}

func doesExist(client *godo.Client) bool {
    ctx := context.TODO()

    opt := &godo.ListOptions{}

    droplets, _, _ := client.Droplets.ListByTag(ctx, tag, opt)

    return len(droplets) > 0
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

    if doesExist(client) {
        deleteDroplet(client)
    } else {
        createDroplet(client)
    }
}
