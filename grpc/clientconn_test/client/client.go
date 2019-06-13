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
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	addr := os.Getenv("ADDR")

	var gRPCFunc func() string
	method := os.Getenv("METHOD")

	if method == "" || method == "ONE_CONNECTION_PER_REQUEST" {
		log.Println("using ONE_CONNECTION_PER_REQUEST...")
		gRPCFunc = func() string {
			// Set up a connection to the server.
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)

			// Contact the server and print out its response.
			name := defaultName
			if len(os.Args) > 1 {
				name = os.Args[1]
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				log.Println("could not greet: ", err)
				return "No Destination"
			}
			return r.Message
		}
	} else if method == "ONE_CONNECTION" {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		var clientPool = sync.Pool{
			New: func() interface{} {
				return pb.NewGreeterClient(conn)
			},
		}

		gRPCFunc = func() string {
			// Set up a connection to the server.

			//conn := conns[rand.Int()%len(conns)]
			//c := pb.NewGreeterClient(conn)
			c := clientPool.Get().(pb.GreeterClient)
			defer clientPool.Put(c)

			// Contact the server and print out its response.
			name := defaultName
			if len(os.Args) > 1 {
				name = os.Args[1]
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				log.Println("could not greet: ", err)
				return "No Destination"
			}
			return r.Message
		}
	} else if method == "CONNECTION_POOL_WITH_EXPANSION" {
		log.Println("using CONNECTION_POOL_WITH_EXPANSION...")
		var conns = make(chan *grpc.ClientConn, 5)
		for i := 0; i < 5; i++ {
			// Set up a connection to the server.
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			conns <- conn
		}

		gRPCFunc = func() string {
			// Set up a connection to the server.
			var conn *grpc.ClientConn
			select {
			case conn = <-conns:
				// success
			default:
				var err error
				conn, err = grpc.Dial(address, grpc.WithInsecure())
				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}
			}

			defer func() {
				select {
				case conns <- conn:
					// return back
				default:
					// recycle
					conn.Close()
				}
			}()

			c := pb.NewGreeterClient(conn)

			// Contact the server and print out its response.
			name := defaultName
			if len(os.Args) > 1 {
				name = os.Args[1]
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				log.Println("could not greet: ", err)
				return "No Destination"
			}
			return r.Message
		}
	}

	http.HandleFunc("/performance", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(gRPCFunc()))
	})

	http.ListenAndServe(addr, nil)
}
