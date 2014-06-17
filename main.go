package main

import "fmt"

import "github.com/flosch/pongo"
import "github.com/hoisie/web"
import "labix.org/v2/mgo"

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
    ctx.SetCookie(web.NewCookie(USER_COOKIE, "", -3600))
    ctx.Redirect(302, "/")
}

// Run the server
func main() {
    mongodb, err := mgo.Dial("localhost")
    if err != nil { panic(err) }
    defer mongodb.Close()
    mongodb.SetMode(mgo.Monotonic, true)
    
    DB_USERS    = mongodb.DB(DB_NAME).C("users")
    DB_SESSIONS = mongodb.DB(DB_NAME).C("sessions")
    
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
