package server

import (
	"database/sql"
	"kabootar/controllers"
	"kabootar/repositories"
	"kabootar/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config          *viper.Viper
	router          *gin.Engine
	usersController *controllers.UsersController
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(c *gin.Context) {
	msgRepository := repositories.NewMsgRepository(dbhandler)
	usersRepository := repositories.NewUsersRepository(dbhandler)
	usersService := services.NewUsersService(usersRepository)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		ok, newErr := usersService.AuthorizeUser(msg.Token)
		if err != nil || newErr != nil || !ok {
			log.Printf("Error reading JSON: %v", err)
			delete(clients, ws)
			break
		}
		msgRepository.SaveMessage(msg.Message, msg.Username)

		broadcast <- msg
	}
}

func PullMsg(ctx *gin.Context) {
	msgRepository := repositories.NewMsgRepository(dbhandler)
	usersRepository := repositories.NewUsersRepository(dbhandler)
	usersService := services.NewUsersService(usersRepository)

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := usersService.AuthorizeUser(accessToken)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	response, responseErr := msgRepository.PullMessage()
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing JSON: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// func (hs *HttpServer) chatRoomHandler(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Error upgrading to WebSocket:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		messageType, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			break
// 		}

// 		if err := conn.WriteMessage(messageType, message); err != nil {
// 			log.Println("Error writing message:", err)
// 			break
// 		}
// 	}
// }

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	usersRepository := repositories.NewUsersRepository(dbHandler)
	usersService := services.NewUsersService(usersRepository)
	usersController := controllers.NewUsersController(usersService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://kabootarme.vercel.app", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin", "Accept", "X-Requested-With", "token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/login", usersController.Login)
	router.POST("/logout", usersController.Logout)
	router.GET("/getmsg", PullMsg)
	router.POST("/create", usersController.CreateUser)

	hs := HttpServer{
		config:          config,
		router:          router,
		usersController: usersController,
	}
	go handleMessages()
	router.GET("/ws", handleConnections)

	return hs
}

func (hs HttpServer) Start() {
	//err := hs.router.Run(hs.config.GetString("http.server_address"))
	err := hs.router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
