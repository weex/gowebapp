package lnd

import (
    "log"
    "io/ioutil"
    // "os"
    // "os/user"
    "time"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    macaroon "gopkg.in/macaroon.v2"
    "github.com/weex/gowebapp/macaroons"
    pb "github.com/weex/gowebapp/gateways/lnd"
)


// this will have lnd specific config
type Config struct {
    DataDir string
}

const (
    address     = "localhost:10009"
    certFile    = "tls.cert"
    macFile     = "admin.macaroon"
)

func InitLnd(cfg Config) (*LndLn, error) {
    // get user's home
    // var homeDir string

    // user, err := user.Current()
    // if err == nil {
    //    homeDir = user.HomeDir
    // } else {
    //    homeDir = os.Getenv("HOME")
    // }

    // get macaroon
    macBytes, err := ioutil.ReadFile(cfg.DataDir +  "/" + macFile)
    if err != nil {
        log.Fatalf("could not find macaroon: %v", err)
    }
    mac := &macaroon.Macaroon{}
    if err = mac.UnmarshalBinary(macBytes); err != nil {
        log.Fatalf("could not unmarshall macaroon: %v", err)
    }


    // Set up a connection to lnd.
    creds, err := credentials.NewClientTLSFromFile(cfg.DataDir + "/" + certFile, "")
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
    // defer conn.Close()
    l := &LndLn{client: pb.NewLightningClient(conn)}
    return l, nil
}

type LndLn struct {
    client pb.LightningClient
}

func (l *LndLn) ListPeers() (*pb.ListPeersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

    r, err := l.client.ListPeers(ctx, &pb.ListPeersRequest{})
    if err != nil {
        log.Fatalf("could not get peers: %v", err)
    }
    for _, peer := range r.Peers {
		log.Printf("Peer. %s", peer)
	}
    return r, nil
}

func (l *LndLn) MakeInvoice(sats int) (*pb.AddInvoiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

    invoice := &pb.Invoice{
        Memo:            "this is a test",
        //Receipt:         "",
        //RPreimage:       "",
        Value:           100,
        //DescriptionHash: "",
        //FallbackAddr:    "",
        Expiry:          86400,
        Private:         true,
    }

    r, err := l.client.AddInvoice(ctx, invoice)
    if err != nil {
        log.Fatalf("could not get invoice: %v", err)
    }
    return r, nil
}
