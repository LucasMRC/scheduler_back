package database

//
// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"
//
// 	"os"
//
// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/jomei/notionapi"
// )
//
// type NotionAPI struct {
// 	Client *notionapi.Client
// }
//
// func (n NotionAPI) GetUser(alias string) (User, error) {
// 	DB_USERS := os.Getenv("DB_USERS")
//
// 	response, err := n.Client.Database.Query(
// 		context.Background(),
// 		notionapi.DatabaseID(DB_USERS),
// 		&notionapi.DatabaseQueryRequest{
// 			Filter: notionapi.PropertyFilter{
// 				Property: "username",
// 				RichText: &notionapi.TextFilterCondition{
// 					Equals: strings.ToLower(alias),
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println("⚠️ Error querying page: ", err.Error())
// 		return User{}, err
// 	}
// 	var value string
// 	for _, page := range response.Results {
// 		name := page.Properties["username"].(*notionapi.TitleProperty)
// 		value = name.Title[0].PlainText
// 		if value == alias {
// 			user := User{}
// 			password, _ := page.Properties["password"].(*notionapi.RichTextProperty)
// 			account, _ := page.Properties["account"].(*notionapi.PeopleProperty)
// 			user.Alias = value
// 			user.Hash = password.RichText[0].PlainText
// 			user.Email = account.People[0].Person.Email
// 			return user, nil
// 		}
// 	}
// 	return User{}, errors.New("User not found")
// }
//
// func (n NotionAPI) SaveUser(user User) error {
// 	DB_USERS := os.Getenv("DB_USERS")
// 	userId := os.Getenv(strings.ToUpper(user.Alias) + "_ID")
//
// 	_, err := n.Client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
// 		Parent: notionapi.Parent{
// 			DatabaseID: notionapi.DatabaseID(DB_USERS),
// 		},
// 		Properties: notionapi.Properties{
// 			"username": notionapi.TitleProperty{
// 				Title: []notionapi.RichText{
// 					{
// 						Type: "text",
// 						Text: &notionapi.Text{
// 							Content: user.Alias,
// 						},
// 					},
// 				},
// 			},
// 			"password": notionapi.RichTextProperty{
// 				RichText: []notionapi.RichText{
// 					{
// 						Type: "text",
// 						Text: &notionapi.Text{
// 							Content: user.Hash,
// 						},
// 					},
// 				},
// 			},
// 			"account": notionapi.PeopleProperty{
// 				People: []notionapi.User{
// 					{
// 						Person: &notionapi.Person{
// 							Email: user.Email,
// 						},
// 						ID: notionapi.UserID(userId),
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println("⚠️ Error creating page: ", err.Error())
// 		return err
// 	}
// 	fmt.Println("✅ User added to Notion successfully")
// 	return nil
// }
//
// func (n NotionAPI) CreateTask(task Task) error {
// 	DB_TASKS := os.Getenv("DB_TASKS")
// 	date, err := time.Parse(time.DateOnly, task.DueDate)
// 	if err != nil {
// 		fmt.Println("⚠️ Error parsing date: ", err.Error())
// 		return err
// 	}
// 	user, err := n.GetUser(task.AssignedTo)
// 	if err != nil {
// 		fmt.Println("⚠️ Error getting user: ", err.Error())
// 		return err
// 	}
// 	userId := os.Getenv(strings.ToUpper(task.AssignedTo) + "_ID")
// 	dateStart := notionapi.Date(date)
// 	_, err = n.Client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
// 		Parent: notionapi.Parent{
// 			DatabaseID: notionapi.DatabaseID(DB_TASKS),
// 		},
// 		Properties: notionapi.Properties{
// 			"id": notionapi.NumberProperty{
// 				Number: float64(task.Id),
// 			},
// 			"title": notionapi.TitleProperty{
// 				Title: []notionapi.RichText{
// 					{
// 						Type: "text",
// 						Text: &notionapi.Text{
// 							Content: task.Title,
// 						},
// 					},
// 				},
// 			},
// 			"description": notionapi.RichTextProperty{
// 				RichText: []notionapi.RichText{
// 					{
// 						Type: "text",
// 						Text: &notionapi.Text{
// 							Content: task.Description,
// 						},
// 					},
// 				},
// 			},
// 			"due date": notionapi.DateProperty{
// 				Date: &notionapi.DateObject{
// 					Start: &dateStart,
// 				},
// 			},
// 			"recurring": notionapi.SelectProperty{
// 				Select: notionapi.Option{
// 					Name: task.Recurring,
// 				},
// 			},
// 			"assigned to": notionapi.PeopleProperty{
// 				People: []notionapi.User{
// 					{
// 						Person: &notionapi.Person{
// 							Email: user.Email,
// 						},
// 						ID: notionapi.UserID(userId),
// 					},
// 				},
// 			},
// 			"status": notionapi.SelectProperty{
// 				Select: notionapi.Option{
// 					Name: task.Status,
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println("⚠️ Error creating task: ", err.Error())
// 		return err
// 	}
// 	fmt.Println("✅ Task added to Notion successfully")
// 	return nil
// }
//
// func (n NotionAPI) GetTasks(username string) ([]Task, error) {
// 	DB_TASKS := os.Getenv("DB_TASKS")
// 	userId := os.Getenv(strings.ToUpper(username) + "_ID")
// 	response, err := n.Client.Database.Query(
// 		context.Background(),
// 		notionapi.DatabaseID(DB_TASKS),
// 		&notionapi.DatabaseQueryRequest{
// 			Filter: notionapi.PropertyFilter{
// 				Property: "assigned to",
// 				People: &notionapi.PeopleFilterCondition{
// 					Contains: userId,
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println("⚠️ Error querying page: ", err.Error())
// 		return []Task{}, err
// 	}
// 	taskList := make([]Task, 0)
// 	for _, page := range response.Results {
// 		task := Task{}
// 		title := page.Properties["title"].(*notionapi.TitleProperty)
// 		description, _ := page.Properties["description"].(*notionapi.RichTextProperty)
// 		dueDate, _ := page.Properties["due date"].(*notionapi.DateProperty)
// 		recurring, _ := page.Properties["recurring"].(*notionapi.SelectProperty)
// 		assignedTo, _ := page.Properties["assigned to"].(*notionapi.PeopleProperty)
// 		status, _ := page.Properties["status"].(*notionapi.SelectProperty)
// 		id, _ := page.Properties["id"].(*notionapi.NumberProperty)
// 		if len(title.Title) > 0 {
// 			task.Title = title.Title[0].PlainText
// 		}
// 		if len(description.RichText) > 0 {
// 			task.Description = description.RichText[0].PlainText
// 		}
// 		task.DueDate = dueDate.Date.Start.String()
// 		task.Recurring = recurring.Select.Name
// 		if len(assignedTo.People) > 0 {
// 			task.AssignedTo = assignedTo.People[0].Person.Email
// 		}
// 		task.Status = status.Select.Name
// 		task.Id = int(id.Number)
// 		taskList = append(taskList, task)
// 	}
// 	return taskList, nil
// }
//
// func (n NotionAPI) UpdateTask(taskId int, patch Task) (Task, error) {
// 	DB_TASKS := os.Getenv("DB_TASKS")
// 	taskIdInt := float64(taskId)
// 	response, err := n.Client.Database.Query(
// 		context.Background(),
// 		notionapi.DatabaseID(DB_TASKS),
// 		&notionapi.DatabaseQueryRequest{
// 			Filter: notionapi.PropertyFilter{
// 				Property: "id",
// 				Number: &notionapi.NumberFilterCondition{
// 					Equals: &taskIdInt,
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println("⚠️ Error querying page: ", err.Error())
// 		return Task{}, err
// 	}
// 	task := Task{}
// 	for _, page := range response.Results {
// 		title := page.Properties["title"].(*notionapi.TitleProperty)
// 		description, _ := page.Properties["description"].(*notionapi.RichTextProperty)
// 		dueDate, _ := page.Properties["due date"].(*notionapi.DateProperty)
// 		recurring, _ := page.Properties["recurring"].(*notionapi.SelectProperty)
// 		assignedTo, _ := page.Properties["assigned to"].(*notionapi.PeopleProperty)
// 		status, _ := page.Properties["status"].(*notionapi.SelectProperty)
// 		if patch.Title == "" {
// 			task.Title = title.Title[0].PlainText
// 		} else {
// 			task.Title = patch.Title
// 		}
// 		if patch.Description == "" {
// 			task.Description = description.RichText[0].PlainText
// 		} else {
// 			task.Description = patch.Description
// 		}
// 		if patch.DueDate == "" {
// 			task.DueDate = dueDate.Date.Start.String()
// 		} else {
// 			patchDate, err := time.Parse(time.DateOnly, patch.DueDate)
// 			if err != nil {
// 				fmt.Println("⚠️ Error parsing date: ", err.Error())
// 				return Task{}, err
// 			}
// 			task.DueDate = patchDate.String()
// 		}
// 		if patch.Recurring == "" {
// 			task.Recurring = recurring.Select.Name
// 		} else {
// 			task.Recurring = patch.Recurring
// 		}
// 		if patch.AssignedTo == "" {
// 			task.AssignedTo = assignedTo.People[0].Person.Email
// 		} else {
// 			task.AssignedTo = patch.AssignedTo
// 		}
// 		if patch.Status == "" {
// 			task.Status = status.Select.Name
// 		} else {
// 			task.Status = patch.Status
// 		}
// 	}
// 	date, err := time.Parse(time.DateOnly, task.DueDate)
// 	dateStart := notionapi.Date(date)
// 	_, err = n.Client.Page.Update(
// 		context.Background(),
// 		notionapi.PageID(response.Results[0].ID),
// 		&notionapi.PageUpdateRequest{
// 			Properties: notionapi.Properties{
// 				"title": notionapi.TitleProperty{
// 					Title: []notionapi.RichText{
// 						{
// 							Type: "text",
// 							Text: &notionapi.Text{
// 								Content: task.Title,
// 							},
// 						},
// 					},
// 				},
// 				"description": notionapi.RichTextProperty{
// 					RichText: []notionapi.RichText{
// 						{
// 							Type: "text",
// 							Text: &notionapi.Text{
// 								Content: task.Description,
// 							},
// 						},
// 					},
// 				},
// 				"due date": notionapi.DateProperty{
// 					Date: &notionapi.DateObject{
// 						Start: &dateStart,
// 					},
// 				},
// 				"recurring": notionapi.SelectProperty{
// 					Select: notionapi.Option{
// 						Name: task.Recurring,
// 					},
// 				},
// 				"status": notionapi.SelectProperty{
// 					Select: notionapi.Option{
// 						Name: task.Status,
// 					},
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println("⚠️ Error updating page: ", err.Error())
// 		return Task{}, err
// 	}
// 	fmt.Println("✅ Task updated in Notion successfully")
// 	return patch, nil
// }
//
// func (n NotionAPI) DeleteTask(taskId int) error {
// 	DB_TASKS := os.Getenv("DB_TASKS")
// 	taskIdInt := float64(taskId)
// 	response, err := n.Client.Database.Query(
// 		context.Background(),
// 		notionapi.DatabaseID(DB_TASKS),
// 		&notionapi.DatabaseQueryRequest{
// 			Filter: notionapi.PropertyFilter{
// 				Property: "id",
// 				Number: &notionapi.NumberFilterCondition{
// 					Equals: &taskIdInt,
// 				},
// 			},
// 		},
// 	)
// 	if err != nil {
// 		fmt.Println("⚠️ Error querying page: ", err.Error())
// 		return err
// 	}
// 	_, err = n.Client.Page.Update(
// 		context.Background(),
// 		notionapi.PageID(response.Results[0].ID),
// 		&notionapi.PageUpdateRequest{
// 			Archived: true,
// 			Properties: notionapi.Properties{
// 				"status": notionapi.SelectProperty{
// 					Select: notionapi.Option{
// 						Name: "Deleted",
// 					},
// 				},
// 			},
// 		},
// 	)
//
// 	if err != nil {
// 		fmt.Println("⚠️ Error deleting page: ", err.Error())
// 		return err
// 	}
// 	fmt.Println("✅ Task deleted from Notion successfully")
// 	return nil
// }
//
// func (n NotionAPI) GetUsers() ([]User, error) {
// 	DB_USERS := os.Getenv("DB_USERS")
// 	response, err := n.Client.Database.Query(
// 		context.Background(),
// 		notionapi.DatabaseID(DB_USERS),
// 		&notionapi.DatabaseQueryRequest{},
// 	)
// 	if err != nil {
// 		fmt.Println("⚠️ Error querying page: ", err.Error())
// 		return []User{}, err
// 	}
// 	userList := make([]User, 0)
// 	for _, page := range response.Results {
// 		user := User{}
// 		username := page.Properties["username"].(*notionapi.TitleProperty)
// 		password, _ := page.Properties["password"].(*notionapi.RichTextProperty)
// 		account, _ := page.Properties["account"].(*notionapi.PeopleProperty)
// 		user.Alias = username.Title[0].PlainText
// 		user.Hash = password.RichText[0].PlainText
// 		user.Email = account.People[0].Person.Email
// 		userList = append(userList, user)
// 	}
// 	return userList, nil
// }
