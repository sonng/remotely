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

    createRequest := &godo.DropletCreateRequest {
        Name: dropletName,
        Region: "sfo2",
        Size: "s-1vcpu-1gb",
        Image: godo.DropletCreateImage {
            Slug: "ubuntu-14-04-x64",
        },
    }

    ctx := context.TODO()

    newDroplet, _, err := client.Droplets.Create(ctx, createRequest)

    if err != nil {
        fmt.Printf("Something bad happened: %s\n\n", err)
    }

   fmt.Printf("New Droplet Created: %s\n\n", newDroplet.Name)
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

    createDroplet(client)
}
