package lnd

import (
    "log"
    "io/ioutil"
    "os"
    "os/user"

    //    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    macaroon "gopkg.in/macaroon.v2"
    "github.com/weex/gowebapp/macaroons"
    pb "github.com/weex/gowebapp/gateways/lnd"
)


// this will have lnd specific config
type Config struct {
    dataDir string
}

const (
    address     = "localhost:10009"
    certFile    = ".lnd/tls.cert"
    macFile     = ".lnd/admin.macaroon"
)


func InitLnd(cfg Config) (pb.LightningClient, error) {
    // get user's home
    var homeDir string

    user, err := user.Current()
    if err == nil {
        homeDir = user.HomeDir
    } else {
        homeDir = os.Getenv("HOME")
    }

    // get macaroon
    macBytes, err := ioutil.ReadFile(homeDir + "/" + macFile)
    if err != nil {
        log.Fatalf("could not find macaroon: %v", err)
    }
    mac := &macaroon.Macaroon{}
    if err = mac.UnmarshalBinary(macBytes); err != nil {
        log.Fatalf("could not unmarshall macaroon: %v", err)
    }


    // Set up a connection to lnd.
    creds, err := credentials.NewClientTLSFromFile(homeDir + "/" + certFile, "")
    if err != nil {
        log.Fatalf("could not get creds: %v", err)
    }

    credMac := macaroons.NewMacaroonCredential(mac)

    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(credMac),
    }

    conn, err := grpc.Dial(address, opts...)
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    return pb.NewLightningClient(conn), nil
}
