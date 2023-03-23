package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	pb "grpc-crud/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Movie struct {
	Id    int32  `json:"id" gorm:"primaryKey autoIncrement"`
	Title string `json:"title"`
	Genre string `json:"genre"`
}

const port = ":8080"

var DB *gorm.DB

// Connect to DB
func SetupDB() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
	}
	db.AutoMigrate(&Movie{})
	DB = db
	log.Println("Connected to database..")
}

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Failed to load .env file")
	}
	SetupDB()
}

type movieServer struct {
	pb.MovieServiceServer
}

// Get all movies
func (s *movieServer) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	fmt.Println("Get Movies")
	movies := []*pb.Movie{}
	res := DB.Find(&movies)
	if res.RowsAffected == 0 {
		return nil, errors.New("Movie not found")
	}
	return &pb.GetMoviesResponse{
		Movies: movies,
	}, nil
}

// Get a movie record
func (s *movieServer) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	fmt.Println("Get Movie")
	var movie Movie
	res := DB.First(&movie, req.Id)
	if res.RowsAffected == 0 {
		return nil, errors.New("Movie not found")
	}
	return &pb.GetMovieResponse{
		Movie: &pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

// Create a movie record
func (s *movieServer) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	fmt.Println("Create Movie")
	movie := req.Movie

	data := Movie{
		Title: movie.Title,
		Genre: movie.Genre,
	}
	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("Error in movie creation")
	}
	return &pb.CreateMovieResponse{
		Movie: &pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

// Update a movie record
func (s *movieServer) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	fmt.Println("Update movie record")
	var movie Movie
	reqMovie := req.Movie
	res := DB.Model(&movie).Where("id = ?", reqMovie.Id).Updates(Movie{
		Title: reqMovie.Title,
		Genre: reqMovie.Genre,
	})
	if res.RowsAffected == 0 {
		return nil, errors.New("Movie not found")
	}
	return &pb.UpdateMovieResponse{
		Movie: &pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

// Delete a movie record
func (s *movieServer) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	fmt.Println("Delete Movie record")
	var movie Movie
	res := DB.Delete(&movie, req.Id)
	// res := DB.Where("id = ?", req.Id).Delete(&movie)
	if res.RowsAffected == 0 {
		return nil, errors.New("Movie not found")
	}

	return &pb.DeleteMovieResponse{
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMovieServiceServer(grpcServer, &movieServer{})
	log.Printf("Server listening at: %v", lis.Addr())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
