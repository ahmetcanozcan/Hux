<p align="center">
  <a ">
    <img
      src="http://ahmetcanozcan.github.io/assets/img/huxlogo.png"
      width="300"
    />
  </a>
</p>

[![CircleCI](https://img.shields.io/circleci/build/github/circleci/circleci-docs?style=flat-square)](https://circleci.com/gh/ahmetcanozcan/hux) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ahmetcanozcan/hux?style=flat-square) ![Codacy grade](https://img.shields.io/codacy/grade/2b1934e3704e44069f7a5c6e89afeca0?style=flat-square)

Hux is a channel based secure WebSocket manager 

## Installation

Use `go get`  to install  Hux

```bash
go get  github.com/ahmetcanozcan/hux
```

## Usage

Hux provides both server-side and client-side libraries.

Firstly import hux
```go
import (
  // Other libraries
  "github.com/ahmetcanozcan/hux"
)

```

Then create a hub to manage rooms and sockets in main function.
```go
hub :=  hux.NewHub()
```

Now, add a http handler 
```go
http.HandleFunc("/ws/hux", func(w http.ResponseWriter, r *http.Request) {
    // Instantiate a Socket
    socket, _:= hub.InstantiateSocket(w, r)
    for {
      select {
      //Listen a event
      case msg := <-socket.GetEvent("Hello"):
        fmt.Println("GOT:", msg)
      }
    }
})

```


Start listening 
```go
http.ListenAndServe(":8080",nil)
```


more event handler can be added using `case`

```go
case msg := <-socket.GetEvent("Join"):
      fmt.Println("Join:", msg)
      hub.GetRoom(msg).Add(socket)
      hub.GetRoom(msg).Emit("New", "NEW SOCKET CONNECTED.")
```

You can send and receive json


firstly, define a struct for  json data

```go

type Person struct {
  Name string `json:"name"`
}

```

then, handle the event that receive a json

```go
case msg := <- socket.GetEvent("json"):
  var p Person
  msg.ParseJSON(&p)
  socket.Emit("Hello","Wellcome, "+ p.Name )
```


``

On client-side, hux provides a library too.
Firstly add this script block before your  code

```html
<script src="https://unpkg.com/hux-client@1.0.1/hux.minifiy.js"></script>
```

Then, you can write your client-side code like this:
```html
<script>
  //Initialize hux 
  var hux = new Hux();
  // When hux connection is established, open event will be invoked.
  hux.on('open', () => {
  console.log('Connection established');
    // Listen hello events from server
    hux.on("World", () => console.log("GOT MESSAGE"));
    // Send Hello message to server
    hux.emit("Hello", "Hi");
    // JSONs can be sent too
    hux.emit("json",{name:"Test-User"})
  })
</script>
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
