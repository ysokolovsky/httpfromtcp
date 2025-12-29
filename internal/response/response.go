package response

import (
	"fmt"
	"io"
	"strconv"

	"httpfromtcp/internal/headers"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

type StatusCode int

const (
	StatusCode200 StatusCode = 200
	StatusCode400 StatusCode = 400
	StatusCode500 StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusCode200:
		fmt.Fprintf(w.writer, "%s\r\n", "HTTP/1.1 200 OK")
		return nil
	case StatusCode400:
		fmt.Fprintf(w.writer, "%s\r\n", "HTTP/1.1 400 Bad Request")
		return nil
	case StatusCode500:
		fmt.Fprintf(w.writer, "%s\r\n", "HTTP/1.1 500 Internal Server Error")
		return nil
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["content-type"] = "text/plain"
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w.writer, "\r\n")

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	fmt.Fprintf(w.writer, "%s", p)
	return len(p), nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return n, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	fmt.Fprintf(w.writer, "0\r\n")

	for k, v := range h {
		fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w.writer, "\r\n")

	return nil

}
