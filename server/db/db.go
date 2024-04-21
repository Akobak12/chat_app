package db

import (
	"database/sql"
	// "fmt"
	// "log"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "")
	if err != nil {
		return nil, err
	}
	// make email unique

	// db.Exec("DROP TABLE IF EXISTS public.users, public.rooms, public.messages, public.room_members, public.invites, public.files;")

	// tables := []string{
	// 	"public.friends",
	// 	"public.bans",
	// 	"public.role_members",
	// 	"public.roles",
	// 	"public.friend_requests",
	// 	"public.files",
	// 	"public.invites",
	// 	"public.room_members",
	// 	"public.messages",
	// 	"public.channels",
	// 	"public.rooms",
	// 	"public.users",
	// }

	// // Drop tables in reverse order to handle dependencies
	// for i := len(tables) - 1; i >= 0; i-- {
	// 	table := tables[i]
	// 	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
	// 	if err != nil {
	// 		log.Fatalf("Error dropping table %s: %v\n", table, err)
	// 	}
	// 	fmt.Printf("Dropped table %s\n", table)
	// }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.users (
    "id" bigserial PRIMARY KEY,
    "username" text NOT NULL UNIQUE,
    "email" text NOT NULL UNIQUE,
	"display_name" text NOT NULL DEFAULT '',
	"bio" text NOT NULL DEFAULT '',
    "password" text NOT NULL,
	"theme_id" int NOT NULL DEFAULT 0
); CREATE TABLE IF NOT EXISTS public.rooms (
	"id" bigserial PRIMARY KEY,
	"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"creator_id" bigint NOT NULL REFERENCES public.users(id),
	"active" boolean NOT NULL DEFAULT TRUE,
	"friendship_room" boolean NOT NULL DEFAULT FALSE,
	"name" text NOT NULL );
	CREATE TABLE IF NOT EXISTS public.channels (
		"id" bigserial PRIMARY KEY,
		"room_id" bigint NOT NULL REFERENCES public.rooms(id),
		"name" text NOT NULL,
		"is_voice" boolean NOT NULL DEFAULT FALSE
	
); CREATE TABLE IF NOT EXISTS public.messages (
	"id" bigserial PRIMARY KEY,
	"content" text NOT NULL,
	"sender_id" bigint NOT NULL REFERENCES public.users(id),
	"channel_id" bigint NOT NULL REFERENCES public.channels(id),
	"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE TABLE IF NOT EXISTS public.room_members (
	"join_index" bigserial NOT NULL PRIMARY KEY,
	"user_id" bigint NOT NULL REFERENCES public.users(id),
	"room_id" bigint NOT NULL REFERENCES public.rooms(id),
	"joined_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE TABLE IF NOT EXISTS public.invites (
	"id" bigserial PRIMARY KEY,
	"room_id" bigint NOT NULL REFERENCES public.rooms(id),
	"inviter_id" bigint NOT NULL REFERENCES public.users(id)
); CREATE TABLE IF NOT EXISTS public.files (
	"hash" bigint NOT NULL PRIMARY KEY
); CREATE TABLE IF NOT EXISTS public.friend_requests (
	"id" bigserial PRIMARY KEY,
	"sender_id" bigint NOT NULL REFERENCES public.users(id),
	"receiver_id" bigint NOT NULL REFERENCES public.users(id),
	"created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
); CREATE TABLE IF NOT EXISTS public.roles (
	"id" bigserial PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	"room_id" bigint NOT NULL REFERENCES public.rooms(id),
	"can_modify_roles" boolean NOT NULL DEFAULT FALSE,
	"can_kick" boolean NOT NULL DEFAULT FALSE,
	"can_ban" boolean NOT NULL DEFAULT FALSE
); CREATE TABLE IF NOT EXISTS public.role_members (
	"role_id" bigint NOT NULL REFERENCES public.roles(id),
	"user_id" bigint NOT NULL REFERENCES public.users(id)
); CREATE TABLE IF NOT EXISTS public.bans (
	"id" bigserial PRIMARY KEY,
	"room_id" bigint NOT NULL REFERENCES public.rooms(id),
	"user_id" bigint NOT NULL REFERENCES public.users(id)
); CREATE TABLE IF NOT EXISTS public.friends (
	"id" bigserial PRIMARY KEY,
	"user_id" bigint NOT NULL REFERENCES public.users(id),
	"friend_id" bigint NOT NULL REFERENCES public.users(id)
);`)

	// drop table users;

	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (db *Database) Close() {
	db.db.Close()
}

func (db *Database) GetDB() *sql.DB {
	return db.db
}
