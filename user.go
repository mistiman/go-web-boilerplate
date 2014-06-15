package main

import "github.com/hoisie/web"
import "labix.org/v2/mgo/bson"

func IfUserRedirect(ctx *web.Context, path string) bool {
    user, err := CurrentUser(ctx)
    if err == nil && user != nil {
        ctx.Redirect(302, path)
        return true
    }
    return false
}

func CurrentUserOrRedirect(ctx *web.Context, path string) *User {
    user, err := CurrentUser(ctx)
    if err != nil {
        ctx.Redirect(302, path)
    }
    return user
}

func CurrentUser(ctx *web.Context) (*User, error) {
    cookie, err := ctx.Request.Cookie("application-login")
    if err != nil {
        return nil, err
    }
    
    session := &Session{}
    err = DB_SESSIONS.Find(bson.M{"uid": cookie.Value}).One(&session)
    if err != nil {
        ERRSOURCE = "sessions: "+cookie.Value
        return nil, err
    }
    
    user := &User{}
    err = DB_USERS.Find(bson.M{"email": session.User}).One(&user)
    if err != nil {
        ERRSOURCE = "users"
        return nil, err
    }
    return user, nil
}
