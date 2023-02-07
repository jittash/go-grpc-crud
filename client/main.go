package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/jittash/go-grpc-crud/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	fmt.Println("Welcome to User Management Client")
	//Dial the connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println("Could not connect", err)
	}
	defer conn.Close()

	//Create a new client
	c := pb.NewUsersClient(conn)

	//Create new user request

	req := &pb.CreateNewUserRequest{
		User: &pb.User{
			Name: &pb.User_Name{
				FirstName: "Jittash", LastName: "Mandhyannee",
			},
			Email:       "j7m@test.com",
			PhoneNumber: 9,
		},
	}

	res, err := c.CreateNewUser(context.Background(), req)
	if err != nil {
		log.Println("Error in creating new user")
	}
	fmt.Println("Created user is", res)

	//Get all Users request

	/*
		req := &empty.Empty{}
		res, err := c.GetAllUsers(context.Background(), req)
		if err != nil {
			log.Println("Error in getting all users")
		}
		fmt.Println(res)
	*/

	//Get a user

	/*
		req := &pb.GetUserRequest{Id: 4}
		res, err := c.GetUser(context.Background(), req)
		if err != nil {
			log.Println("Error in getting user")
		}
		fmt.Println("Retrieved User is", res)
	*/

	//Update User

	/*
		req := pb.UpdateUserRequest{User: &pb.User{
			Id: 4,
			Name: &pb.User_Name{
				FirstName: "Jay",
				LastName:  "Yadav",
			},
			Email:       "jayyadav@gmail.com",
			PhoneNumber: 8,
		}}
		res, err := c.UpdateUser(context.Background(), &req)
		if err != nil {
			log.Println("Error in update user")
		}
		fmt.Println("Updated User is", res)
	*/

	//Delete User

	/*
		req := &pb.DeleteUserRequest{Id: 4}
		_, err = c.DeleteUser(context.Background(), req)
		if err != nil {
			log.Println("user deletion error")
		} else {
			log.Println("User Deleted")
		}
	*/
}
