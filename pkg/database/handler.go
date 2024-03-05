package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"os"

	"github.com/LucasMRC/lb_back/pkg/tasks"
	"github.com/LucasMRC/lb_back/pkg/users"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jomei/notionapi"
)

var tokens []string
var dbClien *sql.DB

func GetUser(username string) (users.User,  error) {
    DB_TOKEN := os.Getenv("DB_TOKEN")
    DB_USERS := os.Getenv("DB_USERS")
    client := notionapi.NewClient(notionapi.Token(DB_TOKEN))
   
    response, err := client.Database.Query(
        context.Background(),
        notionapi.DatabaseID(DB_USERS),
        &notionapi.DatabaseQueryRequest{
            Filter: notionapi.PropertyFilter{
                Property: "username",
                RichText: &notionapi.TextFilterCondition{
                    Equals: username,
                },
            },
        },
    )
    if err != nil {
        fmt.Println("⚠️ Error querying page: ", err.Error())
        return users.User{}, err
    }
    var value string
    for _, page := range response.Results{
        name := page.Properties["username"].(*notionapi.TitleProperty)
        value = name.Title[0].PlainText
        if value == username {
            user := users.User{}
            password, _ := page.Properties["password"].(*notionapi.RichTextProperty)
            account, _ := page.Properties["account"].(*notionapi.PeopleProperty)
            user.Username = value
            user.Password = password.RichText[0].PlainText
            user.Email = account.People[0].Person.Email
            return user, nil
        }
    }
    return users.User{}, errors.New("User not found")
}

func SaveUser(user users.User) error {
    DB_TOKEN := os.Getenv("DB_TOKEN")
    DB_USERS := os.Getenv("DB_USERS")
    client := notionapi.NewClient(notionapi.Token(DB_TOKEN))
    userId := os.Getenv(strings.ToUpper(user.Username) + "_ID")
   
    _, err := client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
        Parent: notionapi.Parent{
            DatabaseID: notionapi.DatabaseID(DB_USERS),
        },
        Properties: notionapi.Properties{
            "username": notionapi.TitleProperty{
                Title: []notionapi.RichText{
                    {
                        Type: "text",
                        Text: &notionapi.Text{
                            Content: user.Username,
                        },
                    },
                },
            },
            "password": notionapi.RichTextProperty{
                RichText: []notionapi.RichText{
                    {
                        Type: "text",
                        Text: &notionapi.Text{
                            Content: user.Password,
                        },
                    },
                },
            },
            "account": notionapi.PeopleProperty{
                People: []notionapi.User{
                    {
                        Person: &notionapi.Person{
                            Email: user.Email,
                        },
                        ID: notionapi.UserID(userId),
                    },
                },
            },
        },
    })
    if err != nil {
        fmt.Println("⚠️ Error creating page: ", err.Error())
        return err
    }
    return nil
}

func CreateTask(task tasks.Task) error {
    DB_TOKEN := os.Getenv("DB_TOKEN")
    DB_TASKS := os.Getenv("DB_TASKS")
    client := notionapi.NewClient(notionapi.Token(DB_TOKEN))
    date, err := time.Parse(time.DateOnly, task.DueDate)
    if err != nil {
        fmt.Println("⚠️ Error parsing date: ", err.Error())
        return err
    }
    user, err := GetUser(task.AssignedTo)
    if err != nil {
        fmt.Println("⚠️ Error getting user: ", err.Error())
        return err
    }
    userId := os.Getenv(strings.ToUpper(task.AssignedTo) + "_ID")
    dateStart := notionapi.Date(date)
    _, err = client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
        Parent: notionapi.Parent{
            DatabaseID: notionapi.DatabaseID(DB_TASKS),
        },
        Properties: notionapi.Properties{
            "title": notionapi.TitleProperty{
                Title: []notionapi.RichText{
                    {
                        Type: "text",
                        Text: &notionapi.Text{
                            Content: task.Title,
                        },
                    },
                },
            },
            "description": notionapi.RichTextProperty{
                RichText: []notionapi.RichText{
                    {
                        Type: "text",
                        Text: &notionapi.Text{
                            Content: task.Description,
                        },
                    },
                },
            },
            "due date": notionapi.DateProperty{
                Date: &notionapi.DateObject{
                    Start: &dateStart,
                },
            },
            "recurring": notionapi.SelectProperty{
                Select: notionapi.Option{
                    Name: task.Recurring,
                },
            },
            "assigned to": notionapi.PeopleProperty{
                People: []notionapi.User{ 
                    {
                        Person: &notionapi.Person{
                            Email: user.Email,
                        },
                        ID: notionapi.UserID(userId),
                    },
                },
            },
            "status": notionapi.SelectProperty{
                Select: notionapi.Option{
                    Name: task.Status,
                },
            },
        },
    })
    if err != nil {
        fmt.Println("⚠️ Error creating task: ", err.Error())
        return err
    }
    return nil
}

func SaveToken(token string) error {
    tokens = append(tokens, token)
    return nil
}

func ConnectToDB() error {
    db,err := sql.Open("mysql", os.Getenv("DNS"))
    if err != nil {
        fmt.Println("⚠️ Error connecting to database: ", err.Error())
        return err
    }
    fmt.Println("🔌 Connected to database")
    dbClien = db
    return nil
}
