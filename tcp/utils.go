package tcp

import (
	//"bufio"
	"io"
)

func ReadAll(rd io.Reader) ([]byte, error) {
	//reader := bufio.NewReader(rd)
	message := make([]byte, 0)
	for {
		buf := make([]byte, 8192)
		received, err := rd.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		message = append(message, buf[:received]...)
		if received < 8192 {
			break
		}
	}
	return message, nil
}
