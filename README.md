# templ-quickstart

## Introduction

templ-quickstart provides a quick and easy way to scaffold an Go http server. The tech stack included in this repo includes Go, HTMX, Templ, and Tailwind.

## Core Technologies

As mentioned above, this project depends on some awesome technologies. Let me start by giving credit where credit is due:

- [Go](https://go.dev/) - Version 1.22.0 or greater required
- [Templ](https://templ.guide/)
- [Air](https://github.com/cosmtrek/air)
- [Htmx](https://htmx.org/)
- [Tailwindcss](https://tailwindcss.com/)

## Installation

### Clone the Repository

```bash
git clone https://github.com/phillip-england/templ-quickstart <target-directory>
```

```bash
cd <target-directory>
```

### Install Dependencies

```bash
go mod tidy
```

### Create a .env file and include a PORT variable

```bash
touch .env; 
```

```bash
echo "PORT=8080" > .env
```

## Build Steps and Serving

This project requires a build step. The following are commands needed to build your html and css output.

### Templ HTML Generation

With templ installed and the binary somewhere on your PATH, run the following to generate your HTML components and templates (remove --watch to simply build and not hot reload)

```bash
templ generate --watch
```

### CSS File Generation

With the [Tailwind Binary](https://tailwindcss.com/blog/standalone-cli) installed and moved somewhere on your PATH, run the following to generate your CSS output for your tailwind classes (remove --watch to simply build and not hot reload)

```bash
tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
```

### Serving with Air

With the [Air Binary](https://github.com/cosmtrek/air) installed and moved somewhere on your PATH, run the following to serve and hot reload the application:

```bash
air
```

To configure air, you can modify .air.toml in the root of the project. (it will be auto-generated after the first time you run air in your repo)

## Project Overview

This project has a few core concepts to help you get going, let's start with ./main.go

### Main - ./main.go

This is our applications entry-point and does a few things:

1. Here, we load in our .env file and then we initialize our mux server.

```go
_ = godotenv.Load()
mux := http.NewServeMux()
```

2. We define a few basic routes for our server. I will go into these routes in more depth later. In short, these routes enable you to use static files in your project, to use a favicon.ico, and sets up a view found at "/".

```go
mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
mux.HandleFunc("GET /static/", view.ServeStaticFiles)
mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
    middleware.Chain(w, r, view.Home)
})
```

Please take note of this line here as it will be important in the next section when we discuss middleware:

```go
middleware.Chain(w, r, view.Home)
```

3. We serve our application on the PORT defined at ./.env

```go
fmt.Println(fmt.Sprintf("server is running on port %s", os.Getenv("PORT")))
err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
if err != nil {
    fmt.Println(err)
}
```
### Middleware - ./internal/middleware/middleware.go

Custom middleware can be implemented with ease in this project. Lets first start with our middleware chain.

This function enables you to tack on middleware at the end of a handler instead of having to deeply-nest middleware components (which is what you would usually expect).

```go

type CustomContext struct {
	context.Context
	StartTime time.Time
}

type CustomHandler func(ctx *CustomContext, w http.ResponseWriter, r *http.Request)

type CustomMiddleware func(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error

func Chain(w http.ResponseWriter, r *http.Request, handler CustomHandler, middleware ...CustomMiddleware) {
	customContext := &CustomContext{
		Context:   context.Background(),
		StartTime: time.Now(),
	}
	for _, mw := range middleware {
		err := mw(customContext, w, r)
		if err != nil {
			return
		}
	}
	handler(customContext, w, r)
	Log(customContext, w, r)
}
```

You'll notice we are using a few custom types here. In short, this function works by initializing a custom context, iterating through our middleware, and then finally calling our handler and logger. The custom context is passed through each middleware, enabling you to store and access context values throughout the chain. If a middleware returns an error, the chain will stop executing. This enables you to allow your middleware to write responses early and avoid calling the handler in case of an error.

### Creating Custom Middleware

Let's say you want to create custom middleware. Here is how to do so:

1. If this middleware requires some context, add the context value to the CustomContext type.

```go
type CustomContext struct {
    context.Context
    StartTime time.Time
    NewContextValue string
}
```

2. Define your new middleware functions (remember middleware must match the CustomMiddleware type definition).

```go
// this middleware will be placed early in the chain
func EarlyMiddleware(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	ctx.NewContextValue = "I was set early in the chain" // set your new context value
	return nil
}

// this middleware will be place late in the chain
func LateMiddleware(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	fmt.Println(ctx.NewContextValue) // outputs "I was set early in the chain"
	return nil
}
```

3. Include the middleware in your Chain func in your routes.

```go
// modified version of ./main.go
mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
    middleware.Chain(w, r, view.Home, middleware.EarlyMiddleware, middleware.LateMiddleware)
})
```

That's it! Easily create custom middleware without the need to deeply nest your routes.

### Views - ./internal/view/view.go

Our views are straightforward and rely on templ to generate html content. Here is an example of the Home view found at ./internal/view/view.go

```go
func Home(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // catches 404s, only needed in the '/' route for entire app
		http.NotFound(w, r)
		return
	
	}
	template.Home("Templ Quickstart").Render(ctx, w)
}
```

### Templates - ./internal/template/template.templ

Our templates are included in this file. Here is the Base template discussed in the previous section. This function simply takes in a title and an array of templ.Component. For more info on templ syntax, please visit [Templ.guide](templ.guide)

To put very simple, Base is a 'base-level template' that can take in children. Then, we reuse base in our home template. Please note the sytax for passing children to @Base. Normally you'd expect to pass children as parameters, but with templ, you place children inside brackets.

```html
templ Base(title string) {
    <html>
        <head>
            <meta charset="UTF-8"></meta>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"></meta>
            <script src="https://unpkg.com/htmx.org@1.9.11"></script>
            <link rel="stylesheet" href="/static/css/output.css"></link>
            <title>{title}</title>
        </head>
            @component.Banner()
        <body>
            <main class='p-6 grid gap-4'>
                { children... }
            </main>
        </body>
    </html>
}

templ Home(title string) {
    @Base(title) {
        @component.TextAndTitle("I'm a Component!", "I am included as a content item in the Base Template!")
	    @component.TextAndTitle("I'm another Component!", "I am also included in the Base Template!")
    }
}
```

Also note, htmx and your tailwind output are included in the head of this template:

```html
<script src="https://unpkg.com/htmx.org@1.9.11"></script>
<link rel="stylesheet" href="/static/css/output.css"></link>
```

### Components - ./internal/component/component.templ

Comonents are very similar to templates. Here is an example of the TextAndTitle component used in ./internal/view/view.go

```html
templ TextAndTitle(title string, text string) {
    <div>
        <h1 class='text-lg font-bold'>{title}</h1>
        <p class='text-sm'>{text}</p>
    </div>
}
```