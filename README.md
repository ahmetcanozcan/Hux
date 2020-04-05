<p align="center">
  <a ">
    <img
      src="http://ahmetcanozcan.github.io/assets/img/huxlogo.png"
      width="300"
    />
  </a>
</p>

[![CircleCI](https://img.shields.io/circleci/build/github/circleci/circleci-docs?style=flat-square)](https://circleci.com/gh/ahmetcanozcan/hux) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ahmetcanozcan/hux?style=flat-square) ![Codacy grade](https://img.shields.io/codacy/grade/2b1934e3704e44069f7a5c6e89afeca0?style=flat-square)

Hux is a channel based event abstraction for web-sockets

## Installation

Use `go get`  to install  Hux

```bash
go get -u github.com/ahmetcanozcan/hux
```

## Usage

Hux provides both server-side and client-side libraries.

In your go file, add these blocks

```go


// blocks of code
func main() {
  // blocks of code
  hux.Initialize() // Initialize hux
  h := hux.GetHub()
  go func(){ // Handle hub
    for {
      select {
      //When a web socket connected, this block executes
      case sck := <-h.SocketConnection: 
        go handleSocket(sck)//Handle socket
      case sck := <-h.SocketDisconnection:
        go handleDisconnection(sck)
      }
  }
  }()
  http.ListenAndServe(":8080", nil)
}

// blocks of code

func handleSocket(sck *hux.Socket) {
  fmt.Println("Socket connected.")
  for {
    select {
    // When the socket sends Hello event,
    // this block will be executed
    case data := <-sck.GetEvent("Hello"):
      fmt.Println(data) // Print data
      sck.Emit("Hello", "Hello There!")//Send response to client
    }
  }
}


```

```html

<!-- Add this block before your script -->
<script src="https://unpkg.com/hux-client@1.0.0/hux.minifiy.js"></script>
<script>
  //Initialize hux 
  var hux = new Hux();
  // When hux connection is established, open event will be invoked.
  hux.on('open', () => {
  console.log('Connection established');
    // Listen hello events from server
    hux.on("Hello", () => console.log("GOT MESSAGE"));
    // Send Hello message to server
    hux.emit("Hello", "Hi");

  })
</script>
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
