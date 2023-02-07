package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/jittash/go-grpc-crud/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	host      = "localhost"
	psql_port = 5432
	user      = "admin"
	password  = "root"
	dbname    = "go_practice"
)

const (
	port = ":50051"
)

type userServer struct {
	pb.UnimplementedUsersServer
}

func (s *userServer) CreateNewUser(ctx context.Context, request *pb.CreateNewUserRequest) (*pb.User, error) {
	db := connectDb()
	defer db.Close()
	email := request.GetUser().GetEmail()
	phone_number := request.GetUser().GetPhoneNumber()
	first_name := request.GetUser().GetName().GetFirstName()
	last_name := request.GetUser().GetName().GetLastName()
	created_user := &pb.User{Email: email, PhoneNumber: phone_number, Name: &pb.User_Name{FirstName: first_name, LastName: last_name}}
	fmt.Println("created_user", created_user)
	//Insert in DB
	_, err := db.Exec(`INSERT INTO "User" (first_name, last_name, email, phone_number) VALUES ($1,$2,$3,$4)`, first_name, last_name, email, phone_number)
	if err != nil {
		log.Println("Error while inserting user", err)
	}
	//return the created user to the client
	return created_user, nil
}

func (s *userServer) GetAllUsers(ctx context.Context, request *empty.Empty) (*pb.GetAllUsersResponse, error) {
	db := connectDb()
	rows, err := db.Query(`SELECT * FROM "User"`)
	defer rows.Close()
	defer db.Close()
	if err != nil {
		log.Println("Cannot retrieve users")
	}
	res := []*pb.User{}

	for rows.Next() {
		var (
			id, phone_number             int32
			email, first_name, last_name string
		)
		if err = rows.Scan(&id, &first_name, &last_name, &email, &phone_number); err != nil {
			log.Println("Error scanning rows")
		}
		u := pb.User{
			Id: id,
			Name: &pb.User_Name{
				FirstName: first_name, LastName: last_name,
			},
			Email:       email,
			PhoneNumber: phone_number,
		}
		res = append(res, &u)
	}
	return &pb.GetAllUsersResponse{Users: res}, nil
}

func (s *userServer) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	db := connectDb()
	defer db.Close()
	req_id := request.GetId()
	query := `SELECT * FROM "User" WHERE id=$1`
	var (
		id, phone_number             int32
		email, first_name, last_name string
	)
	row := db.QueryRow(query, req_id)
	err := row.Scan(&id, &first_name, &last_name, &email, &phone_number)
	if err != nil {
		log.Println("Error getting user", err)
	}
	res := pb.User{
		Id: id,
		Name: &pb.User_Name{
			FirstName: first_name,
			LastName:  last_name,
		},
		Email:       email,
		PhoneNumber: phone_number,
	}
	return &res, nil
}

func (s *userServer) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.User, error) {
	db := connectDb()
	defer db.Close()
	id := request.GetUser().GetId()
	email := request.GetUser().GetEmail()
	phone_number := request.GetUser().GetPhoneNumber()
	first_name := request.GetUser().GetName().GetFirstName()
	last_name := request.GetUser().GetName().GetLastName()
	query := `UPDATE "User" SET first_name=$2, last_name=$3, email=$4, phone_number=$5 WHERE id=$1`
	_, err := db.Exec(query, id, first_name, last_name, email, phone_number)
	if err != nil {
		log.Println("Error updating user", err)
	}
	return request.GetUser(), nil
}

func (s *userServer) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*empty.Empty, error) {
	db := connectDb()
	defer db.Close()
	id := request.GetId()
	query := `DELETE FROM "User" WHERE id=$1`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting user", err)
	}
	return &empty.Empty{}, nil
}

func connectDb() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, psql_port, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println("Cannot open sql connection", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("Cannot establish the connection", err)
	}
	return db
}

func main() {
	fmt.Println("Welcome to User Management Server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Failed to listen the port", err)
	}
	//Initialize the server
	s := grpc.NewServer()

	//Register the server as new grpc service
	pb.RegisterUsersServer(s, &userServer{})
	log.Println("Server listening at", lis.Addr())

	//Listen the server
	if err := s.Serve(lis); err != nil {
		log.Panicln("Failed to listen the server", err)
	}
}
