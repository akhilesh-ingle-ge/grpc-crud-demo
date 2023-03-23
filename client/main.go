package main

import (
	"log"
	"net/http"
	"strconv"

	pb "grpc-crud/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:8080"

type Movie struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
	Genre string `json:"genre"`
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect")
	}
	defer conn.Close()

	client := pb.NewMovieServiceClient(conn)

	router := gin.Default()

	// Get all movies
	router.GET("/movies", func(ctx *gin.Context) {
		res, err := client.GetMovies(ctx, &pb.GetMoviesRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Movies not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data":    res,
			"message": "Fetched data successfully",
		})
	})

	// Get a movie record
	router.GET("/movie/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.GetMovie(ctx, &pb.GetMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Movie not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data":    res.Movie,
			"message": "Record fetched successfully",
		})
	})

	// Create a movie record
	router.POST("/movie", func(ctx *gin.Context) {
		var movie Movie

		err := ctx.ShouldBindJSON(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Error in movie creation",
			})
			return
		}
		data := pb.Movie{
			Title: movie.Title,
			Genre: movie.Genre,
		}
		res, err := client.CreateMovie(ctx, &pb.CreateMovieRequest{
			Movie: &data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data":    res.Movie,
			"message": "Movie created successfully",
		})
	})

	// Update a movie record
	router.PUT("/movie/:id", func(ctx *gin.Context) {
		var movie Movie
		id := ctx.Param("id")
		idint, _ := strconv.Atoi(id)
		idInt := int32(idint)
		err := ctx.ShouldBindJSON(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Error in getting user data",
			})
			return
		}
		res, err := client.UpdateMovie(ctx, &pb.UpdateMovieRequest{
			Movie: &pb.Movie{
				Id:    idInt,
				Title: movie.Title,
				Genre: movie.Genre,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Movie not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data":    res.Movie,
			"message": "Record updated successfully",
		})
	})

	// Delete a movie record
	router.DELETE("/movie/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.DeleteMovie(ctx, &pb.DeleteMovieRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Movie not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Status": res.Success,
		})
	})

	// Start server
	router.Run(":8000")
}
