package main

import (
	"log"
	"server/db"
	"server/internal/user"
	"server/internal/ws"
	"server/router"
)

func main() {
	dbConnection, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRepository := user.NewRepository(dbConnection.GetDB())
	userService := user.NewService(userRepository, *dbConnection)
	userHandler := user.NewHandler(userService)

	hub := ws.NewHub(dbConnection)

	rooms, err := dbConnection.GetDB().Query("SELECT id, creator_id FROM public.rooms")
	if err != nil {
		log.Fatalf("could not get rooms: %s", err)
	}

	// todo: refactor later
	for rooms.Next() {
		var id, creatorID uint64
		err = rooms.Scan(&id, &creatorID)
		if err != nil {
			log.Fatalf("could not scan rooms: %s", err)
		}

		hub.Rooms[id] = &ws.Room{
			Id:           id,
			Clients:      make(map[uint64]*ws.Client),
			AllowedUsers: []uint64{creatorID},
		}

		dbUsers, err := dbConnection.GetDB().Query("SELECT user_id FROM public.room_members WHERE room_id = $1", id)
		if err != nil {
			log.Fatalf("could not get users: %s", err)
		}
		for dbUsers.Next() {
			var userID uint64
			err = dbUsers.Scan(&userID)
			if err != nil {
				log.Fatalf("could not scan users: %s", err)
			}
			hub.Rooms[id].AllowedUsers = append(hub.Rooms[id].AllowedUsers, userID)
		}
		if dbUsers.Err() != nil {
			log.Fatalf("could not get users: %s", dbUsers.Err())
		}
	}

	if rooms.Err() != nil {
		log.Fatalf("could not get rooms: %s", rooms.Err())
	}

	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("localhost:3030")

}
