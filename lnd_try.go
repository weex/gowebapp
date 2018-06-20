/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"log"
//	"os"
    "os/user"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
    macaroon "gopkg.in/macaroon.v2"
    pb "github.com/weex/gowebapp/gateways/lnd"
)

const (
	address     = "localhost:10009"
	certFile    = "tls.cert"
    macFile     = "admin.macaroon"
	defaultName = "world"
)

func main() {
    // get user's home
    var homeDir string

	user, err := user.Current()
	if err == nil {
		homeDir = user.HomeDir
	} else {
		homeDir = os.Getenv("HOME")
	}

    // get macaroon
    macBytes, err := ioutil.ReadFile(homeDir + macFile)
    if err != nil {
		fatal(err)
	}
	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		fatal(err)
	}


    // Set up a connection to the server.
	creds, err := credentials.NewClientTLSFromFile(homeDir + certFile, "")
	if err != nil {
		log.Fatalf("could not get creds: %v", err)
	}

    credMac := macaroons.NewMacaroonCredential(mac)

    opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
        grpc.WithPerRPCCredentials(credMac)
	}

	conn, err := grpc.Dial(address, opts...)
	//conn, err := grpc.Dial(address, grpc.WithSecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLightningClient(conn)

	// Contact the server and print out its response.
	//name := defaultName
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
//	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	r, err := c.ListPeers(ctx, &pb.ListPeersRequest{})
	if err != nil {
		log.Fatalf("could not get peers: %v", err)
	}
	for _, peer := range r.Peers {
		log.Printf("Peer. %s", peer.Address)
	}
}
