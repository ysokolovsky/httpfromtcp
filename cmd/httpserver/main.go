package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handlerHttpbin(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		handlerVideo(w, req)
		return
	}
	handler200(w, req)
}

func handlerVideo(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCode200)

	data, _ := os.ReadFile("assets/vim.mp4")

	h := response.GetDefaultHeaders(len(data))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(data)
}

func handlerHttpbin(w *response.Writer, r *request.Request) {
	url := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin")
	resp, err := http.Get("https://httpbin.org" + url)
	if err != nil {
		log.Fatalf("%s", err)
	}
	w.WriteStatusLine(response.StatusCode200)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h["Transfer-Encoding"] = "chunked"
	h["Trailer"] = "X-Content-SHA256, X-Content-Length"

	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)
	responseBody := make([]byte, 0)
	for {
		n, rerr := resp.Body.Read(buffer)
		responseBody = append(responseBody, buffer[:n]...)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}

	sum := sha256.Sum256(responseBody)
	hashString := fmt.Sprintf("%x", sum)

	trailerHeaders := headers.NewHeaders()
	trailerHeaders.Set("X-Content-SHA256", hashString)
	trailerHeaders.Set("X-Content-Length", strconv.Itoa((len(responseBody))))

	w.WriteTrailers(trailerHeaders)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCode400)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCode500)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCode200)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}
