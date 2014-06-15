package main

import "fmt"

import "github.com/flosch/pongo"
import "github.com/hoisie/web"
import "labix.org/v2/mgo"

// setup mongodb
var DB_USERS *mgo.Collection
var DB_SESSIONS *mgo.Collection

var ERRSOURCE string

func home(ctx *web.Context) string {
    // Check if user is logged in.
    user, err := CurrentUser(ctx)
    if err != nil || user == nil {
        fmt.Printf("ERR: %s, %s\n", err, ERRSOURCE)
        return Render("home.html", nil)
    }
    return Render("home.html", &pongo.Context{"user": user})
}

func dashboard(ctx *web.Context) string {
    // Check if user is logged in.
    user, err := CurrentUser(ctx)
    if err != nil || user == nil {
        fmt.Printf("ERR: %s, %s\n", err, ERRSOURCE)
        ctx.Redirect(302, "/login")
        return ""
    }
    
    return Render("dashboard.html", &pongo.Context{"user": user})
}

func logout(ctx *web.Context) {
    ctx.SetCookie(web.NewCookie("application-login", "", -3600))
    ctx.Redirect(302, "/")
}

// Run the server
func main() {
    mongodb, err := mgo.Dial("localhost")
    if err != nil { panic(err) }
    defer mongodb.Close()
    mongodb.SetMode(mgo.Monotonic, true)
    
    DB_USERS    = mongodb.DB("test").C("users")
    DB_SESSIONS = mongodb.DB("test").C("sessions")
    
    web.Config = &web.ServerConfig{
        RecoverPanic: true,
    }
    
    web.Get("/", home)
    
    web.Get("/logout", logout)
    
    web.Get("/login", Login)
    web.Post("/login", Login)
    
    web.Get("/signup", Signup)
    web.Post("/signup", Signup)
    
    web.Get("/dash", dashboard)
    
    // web.Run(":80")
    web.Run(":8080")
}
