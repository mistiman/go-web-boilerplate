package main

// setup mongodb
var DB_USERS *mgo.Collection
var DB_SESSIONS *mgo.Collection

// used for error reporting
var ERRSOURCE string

// The location where the HTML Templates are kept.
const TEMPLATE_PATH string = "templates"

// USER_COOKIE is the name of the cookie used for
// user sessions.
const USER_COOKIE string = "application-login"

// This Email RE will only fit 99% of all emails...
const EMAIL_RE = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
