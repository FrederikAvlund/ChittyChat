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

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"math"
	"strconv"

	//	"encoding/json"
	//	"flag"
	"fmt"
	//	"io"
	//	"io/ioutil"
	"log"
	//	"math"
	"net"
	"sync"

	//	"time"

	"google.golang.org/grpc"

	//	"github.com/golang/protobuf/proto"

	pb "chittychat"
)

var grpcServer pb.ChittyChatServer

type Connection struct {
	stream pb.ChittyChat_JoinServer
	id     string
	active bool
	err    chan error
}

type Server struct {
	Connection []*Connection
	pb.UnimplementedChittyChatServer
	local_timestamp int64
}

func GetTimestamp(s *Server, i int64) int64 {
	l := float64(s.local_timestamp)
	i_ := float64(i)
	f := math.Max(l, i_) + 1
	return int64(f)
}

func (s *Server) Join(pconn *pb.Connect, stream pb.ChittyChat_JoinServer) error {

	var msg pb.Message
	var ctx context.Context
	incoming_timestamp, _ := strconv.Atoi(msg.GetTimestamp())
	s.local_timestamp = GetTimestamp(s, int64(incoming_timestamp))
	fmt.Println("Event: Connect recieved: ", s.local_timestamp)
	conn := &Connection{
		stream: stream,
		//      id: pconn.User.Id,
		id:     pconn.User.DisplayName,
		active: true,
		err:    make(chan error),
	}
	s.Connection = append(s.Connection, conn)
	log.Println("Join of user", conn.id)
	//msg.Message = conn.id + " joined Chitty-Chat (" + string(s.local_timestamp) + ")"
	fmt.Println(conn.id+" joined Chitty-Chat (", s.local_timestamp, ")")
	//	msg.User.DisplayName =  "???"
	s.Broadcast(ctx, &msg)
	s.local_timestamp++

	return <-conn.err
}

func (s *Server) Broadcast(ctx context.Context, msg *pb.Message) (*pb.Close, error) {

	//fmt.Printf("Kald af Broadcast\n")

	wait := sync.WaitGroup{}
	done := make(chan int)
	s.local_timestamp = GetTimestamp(s, s.local_timestamp)

	for _, conn := range s.Connection {
		log.Println(conn.id)
		wait.Add(1)

		go func(msg *pb.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				//err:=error(nil)
				//            log.Println("Sending message %v to user %w", msg.Id, conn.id)
				fmt.Printf("Sending message: %v to user %v \n", msg.Message, conn.id)
				s.local_timestamp++

				if err != nil {
					log.Fatalf("Error with stream %v. Error: %v", conn.stream, err)
					conn.active = false
					conn.err <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	//maybe timestamp++ here?

	<-done

	return &pb.Close{}, nil
}

func (s *Server) Publish(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	incoming_timestamp, _ := strconv.Atoi(msg.GetTimestamp())
	s.local_timestamp = GetTimestamp(s, int64(incoming_timestamp))

	log.Println("Publish() from", msg.User.DisplayName, ":", msg.Message)
	//msg.Message = msg.User.DisplayName + ": " + msg.Message + " (Lamport time xxx)"
	fmt.Println(msg.User.DisplayName+": "+msg.Message+" (", s.local_timestamp, ")")

	s.Broadcast(ctx, msg)
	return &pb.Close{}, nil
}

func (s *Server) Leave(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	log.Println("Leave() from", msg.User.DisplayName, ":", msg.Message)
	incoming_timestamp, _ := strconv.Atoi(msg.GetTimestamp())
	s.local_timestamp = GetTimestamp(s, int64(incoming_timestamp))

	//msg.Message = msg.User.DisplayName + ": Left Chitty-Chat (Lamport time xxx)"
	fmt.Println(msg.User.DisplayName+": Left Chitty-Chat"+" (", s.local_timestamp, ")")
	s.Broadcast(ctx, msg)

	for _, conn := range s.Connection {
		log.Println("conn.id: " + conn.id + ", msg.User.DisplayName: " + msg.User.DisplayName)
		//Kan det logges uden at skrive i terminalen?
		if conn.id == msg.User.DisplayName {
			conn.active = false
		}
	}
	return &pb.Close{}, nil
}

func main() {
	var connections []*Connection
	var ThisBroadcastServer pb.UnimplementedChittyChatServer

	server := &Server{connections, ThisBroadcastServer, 0}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	log.Println("Starting server at port :8080")

	pb.RegisterChittyChatServer(grpcServer, server)
	grpcServer.Serve(listener)
}
