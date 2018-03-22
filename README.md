# Microservice Architecture

## Status

Approved ;)

## Why?

Pods currently don't have the flexibility to do what's best for the pod. Webstore-2 relies on zapi that was designed for webstore-1. We're constantly fighting with the odd ways of how zapi and webstore-1 did things in the past. For example, when a user enters an address we need to tell zapi about this so we can get an address_id back. This address_id is useless to webstore but needed at checkout. Too much logic is being performed on the frontend to determine the state of a registering user; zapi has all this information.

We could go into zapi and update these routes but run the real risk of `.finally is not a function` that our false positive test failed to find because we didn't have chai configured correctly and the moon phase was in its first quarter phase over Easter Island.

I'm proposing not only an architecture but a language addition. The language in this case being [GoLang](https://golang.org/). I will address the advantages of Go at the end of the proposal.

Before starting any application it's important to plan out the organization of the project. Some applications push everything into one package while others group by type or module. Without a good strategy applied to your team, you'll find code scattered across various packages of your application. We need a better standard for our application design because what we have today doesn't scale.

I suggest a better approach. By following a few simple rules we can decouple our code, make it easier to test, and bring a consistent structure to our project. Before we dive into it, though, here are a few other approaches and their weaknesses


### Approach 1: Monolithic package
Throwing all your code in a single package can actually work very well for small applications. It removes any chance of circular dependencies because, within your application, there are no dependencies.

### Approach 2: Functional type layout
Another approach is to group your code by its functional type. For example, all your handlers go in one package, your controllers go in another, and your models go in yet another.

There are two issues with this approach though. First, your names are atrocious. You end up with type names like `controller.UserController` where you're duplicating your package name in your type's name.

The bigger issue, however, is circular dependencies. Your different functional types may need to reference each other. This only works if you have one-way dependencies but many times your application is not that simple.

### Approach 3: Module type layout
This approach is similar to the above style layout except that we are grouping our code by module instead of by function. For example, you may have a users package and an accounts package.

We find the same issues in this approach. Again, we end up with terrible names like `users.User`. We also have the same issue of circular dependencies if our `accounts.Controller` needs to interact with our `users.Controller` and vis-a-versa.

## A better approach
![bubbles](/bubbles.png)
The package strategy that I've used for previous projects involves 4 simple tenets:

1. Root package is for domain types
2. Group subpackages by dependency
3. Use a shared mock subpackage
4. Main package ties together dependencies

These rules help isolate our packages and define a clear domain language across
the entire application. Let’s look at how each one of these rules works in practice.

### 1. Root package is for domain types
Your application has a logical, high-level language that describes how data and processes interact. This is your domain. Your application domain involves things like customers, accounts, charging credit cards, and handling inventory. It’s the stuff that doesn’t depend on our underlying technology.

I place my domain types in my root package. This package only contains simple data types like a User struct for holding user data or a `UserService` interface for fetching or saving user data.

It may look something like:

```go
package myapp

type User struct {
    ID int
    Name string
    Address Address
}

type UserService interface {
    User(id int) (*User, error)
    Users() ([]*User, error)
    CreateUser(u *User) error
    DeleteUser(id int) error
}
```

This makes your root package extremely simple. You may also include types that perform actions but only if they solely depend on other domain types. For example, you could have a type that polls your `UserService` periodically. However, it should not call out to external services or save to a database. That is an implementation detail.

**The root package should not depend on any other package in your application!**

### 2. Group subpackages by dependency
If your root package is not allowed to have external dependencies then we must push those dependencies to subpackages. In this approach to package layout, subpackages exist as an adapter between your domain and your implementation.

For example, your `UserService` would be backed by PostgreSQL. You can introduce a Postgres subpackage in your application that provides a `postgres.UserService` implementation:

```go
package postgres

import (
    "database/sql"

    "github.com/sir-wiggles/myapp"
    _ "github.com/lib/pq"
)

// UserService represents a PostgreSQL implementation of myapp.UserService.
type UserService struct {
    DB *sql.DB
}

// User returns a user for a given id.
func (s *UserService) User(id int) (*myapp.User, error) {
    var u myapp.User
    row := db.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id)
    if row.Scan(&u.ID, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}

// implement remaining myapp.UserService interface...
```

This isolates our PostgreSQL dependency which simplifies testing and provides an easy way to migrate to another database in the future. It can be used as a pluggable architecture if you decide to support other database implementations.

It also gives you a way to layer implementations. Perhaps you want to have a Redis cache in front of PostgreSQL. You can add a `UserCache` that implements `UserService` which can wrap your PostgreSQL implementation:

```go
package myapp

// UserCache wraps a UserService to provide a redis cache.
type UserCache struct {
 cache CacheService
 service UserService
}

// NewUserCache returns a new read-through cache for service.
func NewUserCache(service UserService) *UserCache {
 return &UserCache{
 cache: NewRedisService(),
 service: service,
 }
}

// User returns a user for a given id.
// Returns the cached instance if available.
func (c *UserCache) User(id int) (*User, error) {

 if u := c.cache.Get(id); u != nil {
 return u, nil
 }

 // Otherwise fetch from the underlying service.
 u, err := c.service.User(id)
 if err != nil {
     return nil, err
 } else if u != nil {
     c.cache[id] = u
 }
 return u, err
}
```

This is a common approach I've seen in the go standard library. The io.Reader is a domain type for reading bytes and its implementations are grouped by dependency — tar.Reader, gzip.Reader, multipart.Reader. These can be layered as well. It’s common to see an os.File wrapped by a bufio.Reader which is wrapped by a gzip.Reader which is wrapped in a tar.Reader.

#### Dependencies between dependencies
Your dependencies don’t live in isolation. You may store User data in PostgreSQL but your financial transaction data exists in a third party service like Stripe. In this case, we wrap our Stripe dependency with a logical domain type — let’s call it TransactionService.

By adding our TransactionService to our UserService we decouple our two dependencies:

```go
type UserService struct {
 TransactionService myapp.TransactionService
}
```

Now our dependencies communicate solely through our common domain language. This means that we could swap out PostgreSQL for MySQL or switch Stripe for another payment processor without affecting other dependencies.

### 3. Use a shared mock subpackage
Because our dependencies are isolated from other dependencies by our domain interfaces, we can use these connection points to inject mock implementations.

There are several mocking libraries I personally prefer to just write them myself. I find many of the mocking tools to be overly complicated.

The mocks are very simple. For example, a mock for the UserService looks like:

```go
package mock

import "github.com/sir-wiggles/myapp"

// UserService represents a mock implementation of myapp.UserService.
type UserService struct {
 UserFn func(id int) (*myapp.User, error)
 UserInvoked bool

 UsersFn func() ([]*myapp.User, error)
 UsersInvoked bool

 // additional function implementations...
}

// User invokes the mock implementation and marks the function as invoked.
func (s *UserService) User(id int) (*myapp.User, error) {
 s.UserInvoked = true
 return s.UserFn(id)
}

// additional functions: Users(), CreateUser(), DeleteUser()
```

This mock lets me inject functions into anything that uses the myapp.UserService interface to validate arguments, return expected data, or inject failures.

Let's say we want to test an HTTP handler
```go

package http_test

import (
    "testing"
    "net/http"
    "net/http/httptest"

    "github.com/sir-wiggles/myapp/mock"
)

func TestHandler(t *testing.T) {
    // Inject our mock into our handler.
    var us mock.UserService
    var h Handler
    h.UserService = &us

    // Mock our User() call.
    us.UserFn = func(id int) (*myapp.User, error) {
        if id != 100 {
            t.Fatalf("unexpected id: %d", id)
        }
        return &myapp.User{ID: 100, Name: "susy"}, nil
    }

    // Invoke the handler.
    w := httptest.NewRecorder()
    r, _ := http.NewRequest("GET", "/users/100", nil)
    h.ServeHTTP(w, r)

    // Validate mock.
    if !us.UserInvoked {
        t.Fatal("expected User() to be invoked")
    }
}
```

Our mock lets us completely isolate our unit test to only the handling of the HTTP protocol.

### 4. Main package ties together dependencies
With all these dependency packages floating around in isolation, you may wonder how they all come together. That’s the job of the main package.

Main package layout
An application may produce multiple binaries so we’ll use the Go convention of placing our main package as a subdirectory of the cmd package. For example, our project may have a myapp server binary but also a myappctl client binary for managing the server from the terminal.
```
.
├── cmd
│   ├── cli
│   │   └── main.go
│   └── server
│       └── main.go
├── pkg
│   └── webstore
│       ├── api
│       │   ├── api.go
│       │   └── api_test.go
│       ├── mock
│       │   ├── cache.go
│       │   ├── phone.go
│       │   └── user.go
│       ├── postgres
│       │   └── postgres.go
│       ├── redis
│       │   └── redis.go
│       ├── twilio
│       │   ├── twilio.go
│       │   └── twilio_test.go
│       ├── cache.go
│       ├── database.go
│       ├── phone.go
│       └── user.go
├── docker-compose.yml
├── main.go
├── new_routes.md
└── README.md
```


#### Injecting dependencies at compile time
The term "dependency injection" has gotten a bad rap. It conjures up thoughts of verbose Spring XML files. However, all the term really means is that we’re going to pass dependencies to our objects instead of requiring that the object build or find the dependency itself.

The main package is what gets to choose which dependencies to inject into which objects. Because the main package simply wires up the pieces, it tends to be fairly small and trivial code:

```go

package main

import (
    "log"
    "os"

    "github.com/sir-wiggles/myapp"
    "github.com/sir-wiggles/myapp/postgres"
    "github.com/sir-wiggles/myapp/http"
)

func main() {
    // Connect to database.
    db, err := postgres.Open(os.Getenv("DB"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create services.
    us := &postgres.UserService{DB: db}

    // Attach to HTTP handler.
    var h http.Handler
    h.UserService = us

    // start http server...
}
```

## Advantages

1. very well suited for applications that require high concurrency right out of the box.
2. very concise
3. compiles down to a single binary for native execution.
4. easy-to-learn language whose hallmark is pragmatism.
5. compiles **very** fast.
6. quickly compiles to a single executable
7. interface is very flexible
8. very small memory footprint
9. excellent standard library, can build a fully functional production ready web server with just the standard library with ease.
10. despite being relatively young, the language is very mature and consistent.
11. Opinionated styling with `gofmt`
12. Frameworks are unnecessary with go. They exist but are highly discouraged in the community as they mostly get in the way.
