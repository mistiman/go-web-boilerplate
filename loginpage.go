package main

// import "fmt"
import "math/rand"
import "time"

import "github.com/flosch/pongo"
import "github.com/hoisie/web"
import "github.com/jameskeane/bcrypt"
import "labix.org/v2/mgo/bson"

func RandomString(n int) string {
    letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~!@#$%^&*_-+?")
    b := make([]rune, n)
    
    rand.Seed(time.Since(time.Unix(0,0)).Nanoseconds())
    
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func Login(ctx *web.Context) string {
    // Check if user is logged in.
    if IfUserRedirect(ctx, "/dash") { return "" }
    
    if ctx.Request.Method == "GET" {
        if ctx.Request.FormValue("registration") == "success" {
            return Render("login.html", &pongo.Context{"notification":"Congratulations, your sign up was a success!"})
        } else {
            return Render("login.html", nil)
        }
    } else {
        useremail := ctx.Request.FormValue("email")
        userpassw := ctx.Request.FormValue("password")
        
        user := &User{}
        err := DB_USERS.Find(bson.M{"email": useremail}).One(user)
        if err != nil {
            if err.Error() == "not found" {
                return Render("login.html", &pongo.Context{"error":"Email or password incorrect. Please check your login information."})
            } else {
                panic(err)
            }
        }
        
        if !bcrypt.Match(userpassw, user.Password) {
            return Render("login.html", &pongo.Context{"error":"Email or password incorrect. Please check your login information."})
        }
        
        crla := time.Now()
        suid := RandomString(256)
        DB_SESSIONS.Insert(&Session{User: user.Email, Uid: suid, Created: crla, LastAccess: crla})
        
        ctx.SetCookie(web.NewCookie(USER_COOKIE, suid, 3600))
        
        ctx.Redirect(302, "/dash")
    }
    
    return Render("login.html", nil)
}