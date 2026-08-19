package main

// to avoid name collision of "isAuthenticated" => wrap it around type "contextKey" (instead of plain "string")
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
