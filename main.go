package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"gorestdemo/pkg"
	"log"
	"net/http"
)

func main() {
	fmt.Println(123)
	u := UserResource{map[string]User{}}
	restful.DefaultContainer.Add(u.WebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: pkg.EnrichSwaggerObject}

	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/data/self/swagger/ui"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}


type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type UserResource struct {
	users map[string]User
}

func (u UserResource) createUser(request *restful.Request, response *restful.Response) {
	user := User{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&user)
	if err == nil {
		u.users[user.ID] = user
		response.WriteHeaderAndEntity(http.StatusCreated, user)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (u UserResource) findAllUsers(request *restful.Request, response *restful.Response) {
	list := []User{}
	for _, each := range u.users {
		list = append(list, each)
	}
	response.WriteEntity(list)
}

func (u UserResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/users").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	tags := []string{"users"}

	ws.Route(ws.GET("/").To(u.findAllUsers).
		Doc("get all users").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]User{}).
		Returns(200, "OK", User{}))
	ws.Route(ws.PUT("").To(u.createUser).
		Doc("Create a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{}))
	return ws
}

func EnrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "UserService",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				Name:  "john",
				Email: "john@doe.rp",
				URL:   "http://johndoe.org",
			},
			License: &spec.License{
				Name: "MIT",
				URL:  "http://mit.org",
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users",
	}}}
}
