package main

import "time"
import "regexp"

import "github.com/flosch/pongo"
import "github.com/hoisie/web"
import "github.com/jameskeane/bcrypt"
import "labix.org/v2/mgo/bson"

// This Email RE will only fit 99% of all emails...
const EMAIL_RE = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`

func Signup(ctx *web.Context) string {
    // Check if user is logged in.
    if IfUserRedirect(ctx, "/dash") { return "" }
    
    if ctx.Request.Method == "GET" {
        return Render("signup.html", nil)
    } else {
        newuser := &User{}
        
        // check that email wasn't messed up
        newuser.Email = ctx.Request.FormValue("email")
        if newuser.Email != ctx.Request.FormValue("emailconfirm") {
            return Render("signup.html", &pongo.Context{
                                            "error":"Emails do not match...",
                                            "firstname":ctx.Request.FormValue("firstname"),
                                            "lastname":ctx.Request.FormValue("lastname"),
                                            "email":ctx.Request.FormValue("email"),
                                            "emailconfirm":ctx.Request.FormValue("emailconfirm"),
                                         })
        }
        
        // check that password wasn't messed up
        newuser.Password = ctx.Request.FormValue("password")
        if newuser.Password != ctx.Request.FormValue("passwordconfirm") {
            return Render("signup.html", &pongo.Context{
                                            "error":"Passwords do not match...",
                                            "firstname":ctx.Request.FormValue("firstname"),
                                            "lastname":ctx.Request.FormValue("lastname"),
                                            "email":ctx.Request.FormValue("email"),
                                            "emailconfirm":ctx.Request.FormValue("emailconfirm"),
                                         })
        }
        
        newuser.Firstname = ctx.Request.FormValue("firstname")
        newuser.Lastname = ctx.Request.FormValue("lastname")
        newuser.Since = time.Now()
        
        // check that email is valid
        if re := regexp.MustCompile(EMAIL_RE); re.Match([]byte(newuser.Email)) != true {
            return Render("signup.html", &pongo.Context{
                                            "error":"Email is not valid...",
                                            "firstname":ctx.Request.FormValue("firstname"),
                                            "lastname":ctx.Request.FormValue("lastname"),
                                            "email":ctx.Request.FormValue("email"),
                                            "emailconfirm":ctx.Request.FormValue("emailconfirm"),
                                         })
        }
        
        // check that email is NOT already registered
        uquery := &User{}
        err := DB_USERS.Find(bson.M{"email": newuser.Email}).One(uquery)
        if err == nil {
            return Render("signup.html", &pongo.Context{"error":"Email is already registered!"})
        } else {
            
            // email does not exist
            if err.Error() == "not found" {
                
                // encrypt password
                salt, err := bcrypt.Salt()
                if err != nil { panic(err) }
                hashedpw, err := bcrypt.Hash(newuser.Password, salt)
                if err != nil { panic(err) }
                newuser.Password = hashedpw
                
                // put user into mongo
                err = DB_USERS.Insert(newuser)
                if err != nil { panic(err) }
                
            } else {
                panic(err)
            }
        }
        
        ctx.Redirect(302, "/login?registration=success")
    }
    return ""
}
