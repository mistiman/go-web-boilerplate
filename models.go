package main

import "time"

type User struct {
    Email     string
    Password  string
    
    Firstname string
    Lastname  string
    Since     time.Time
}

type Session struct {
    User       string
    Uid        string
    Created    time.Time
    LastAccess time.Time
}
