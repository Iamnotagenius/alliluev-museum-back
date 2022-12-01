package main

import (
	"strings"

	"github.com/graphql-go/graphql"
)

type (
	Exhibit struct {
		ID          int
		Name        string
		Description string
		Pictures    []string
		Room        int
	}

	Room struct {
		ID       int
		Name     string
		Pictures []string
	}
)

var (
	exhibitType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Exhibit",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return exhibit.Name, nil
					}
					return nil, nil
				},
			},
			"description": &graphql.Field{
				Type:        graphql.NewList(graphql.NewNonNull(graphql.String)),
				Description: "Long description separated to paragraphs",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return strings.FieldsFunc(
							exhibit.Description,
							func(r rune) bool {
								return r == '\n'
							}), nil
					}
					return []interface{}{}, nil
				},
			},
			"pictures": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "Paths to photos of this exhibit",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return exhibit.Pictures, nil
					}
					return []interface{}{}, nil
				},
			},
		},
	})

	roomType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Room",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if room, ok := p.Source.(Room); ok {
						return room.Name, nil
					}
					return nil, nil
				},
			},
			"pictures": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "Paths to photos of this room",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if room, ok := p.Source.(Room); ok {
						return room.Pictures, nil
					}
					return []interface{}{}, nil
				},
			},
			"exhibits": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(exhibitType))),
				Description: "Exhibits located in this room",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db, _ := p.Context.Value("db").(MuseumDB)
					if room, ok := p.Source.(Room); ok {
						return db.ExhibitsByRoomID(room.ID)
					}
					return []interface{}{}, nil
				},
			},
		},
	})

	rootType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Query",
		Description: "Get a singular item from a database",
		Fields: graphql.Fields{
			"exhibit": &graphql.Field{
				Type: exhibitType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "Retreive an exhibit from database",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db, _ := p.Context.Value("db").(MuseumDB)
					return db.ExhibitByID(p.Args["id"].(int))
				},
			},
			"room": &graphql.Field{
				Type: roomType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "Retreive a room from database",
						Type:        graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db, _ := p.Context.Value("db").(MuseumDB)
					return db.RoomByID(p.Args["id"].(int))
				},
			},
			"exhibits": &graphql.Field{
				Type:        graphql.NewList(exhibitType),
				Description: "Query all exhibits from database",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db, _ := p.Context.Value("db").(MuseumDB)
					return db.GetAllExhibits()
				},
			},
			"rooms": &graphql.Field{
				Type:        graphql.NewList(roomType),
				Description: "Query all rooms from database",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db, _ := p.Context.Value("db").(MuseumDB)
					return db.GetAllRooms()
				},
			},
		},
	})
)

func init_schema() (graphql.Schema, error) {
	exhibitType.AddFieldConfig(
		"room", &graphql.Field{
			Type:        roomType,
			Description: "Room in which this exhibit is located",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				db, _ := p.Context.Value("db").(MuseumDB)
				if exhibit, ok := p.Source.(Exhibit); ok {
					return db.RoomByID(exhibit.Room)
				}
				return Room{}, nil
			},
		},
	)

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootType,
	})
}
