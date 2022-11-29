package main

import (
	"github.com/graphql-go/graphql"
)

type (
	Exhibit struct {
		Slug        string
		Name        string
		Description string
		Pictures    []string
		Room        Room
	}

	Room struct {
		Name        string
		Description string
		Pictures    []string
		Exhibits    []Exhibit
	}
)

var (
	Data map[string]Exhibit = map[string]Exhibit{
		"chairs": {
			Slug:        "chairs",
			Name:        "Венские стулья",
			Description: "Стулья из венеции. Второй параграф",
			Pictures:    []string{"assets/pictures/chairs1.jpeg", "assets/pictures/chairs2.jpeg"},
		},
		"clock": {
			Slug:        "clock",
			Name:        "Часы",
			Description: "Часы. Второй параграф",
			Pictures:    []string{"assets/pictures/clock.jpeg"},
		},
	}

	RoomData map[string]Room = map[string]Room{
		"hall": {
			Name:        "Гостиная",
			Description: "Описание гостиной",
			Pictures:    []string{"assets/pictures/chairs1.jpeg", "assets/pictures/chairs2.jpeg"},
			Exhibits:    []Exhibit{Data["chairs"]},
		},
		"kitchen": {
			Name:        "Кухня",
			Description: "Описание",
			Pictures:    []string{"assets/pictures/clock.jpeg"},
			Exhibits:    []Exhibit{Data["clock"]},
		},
	}
)

func init_schema() (graphql.Schema, error) {
	data := Data["chairs"]
	data.Room = RoomData["hall"]
	Data["chairs"] = data
	data = Data["clock"]
	data.Room = RoomData["kitchen"]
	Data["clock"] = data

	var roomType *graphql.Object

	exhibitType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Exhibit",
		Fields: graphql.Fields{
			"slug": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return exhibit.Slug, nil
					}
					return nil, nil
				},
			},
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
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return exhibit.Description, nil
					}
					return []interface{}{}, nil
				},
			},
			"pictures": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return exhibit.Pictures, nil
					}
					return []interface{}{}, nil
				},
			},
			"room": &graphql.Field{
				Type: roomType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := p.Source.(Exhibit); ok {
						return exhibit.Room, nil
					}
					return Room{}, nil
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
			"description": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if room, ok := p.Source.(Room); ok {
						return room.Description, nil
					}
					return []interface{}{}, nil
				},
			},
			"pictures": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if room, ok := p.Source.(Room); ok {
						return room.Pictures, nil
					}
					return []interface{}{}, nil
				},
			},
			"exhibits": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(exhibitType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if room, ok := p.Source.(Room); ok {
						return room.Exhibits, nil
					}
					return []interface{}{}, nil
				},
			},
		},
	})

	rootType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Query",
		Description: "Get a singular item from a database",
		Fields: graphql.Fields{
			"exhibit": &graphql.Field{
				Type: exhibitType,
				Args: graphql.FieldConfigArgument{
					"slug": &graphql.ArgumentConfig{
						Description: "Returns an exhibit based on slug",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if exhibit, ok := Data[p.Args["slug"].(string)]; ok {
						return exhibit, nil
					}
					return Exhibit{}, nil
				},
			},
			"room": &graphql.Field{
				Type: roomType,
				Args: graphql.FieldConfigArgument{
					"slug": &graphql.ArgumentConfig{
						Description: "Returns a room based on slug",
						Type:        graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if room, ok := RoomData[p.Args["slug"].(string)]; ok {
						return room, nil
					}
					return Room{}, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootType,
	})
}
