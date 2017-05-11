package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader {
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

type App struct {
	Engine *echo.Echo
}

func serverWs(hub *Hub, context echo.Context) error {
	conn, err := upgrader.Upgrade(context.Response(), context.Request(), nil)

	if err != nil {
		return err
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}	
	client.hub.register <- client
	go client.openWriter()
	client.openReader()

	return nil
}

func (a *App) Run() {
	a.Engine.Logger.Fatal(a.Engine.Start(":33333"))
}

func NewApp() *App {
	engine := echo.New()
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())

	hub := NewHub()
	go hub.run()

	engine.Static("/", "./client")
	
	engine.GET("/ws", func(context echo.Context) error {
		return serverWs(hub, context)
	})
	
	app := &App {
		Engine: engine,
	}

	return app
}
