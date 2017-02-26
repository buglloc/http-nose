package httpclient

import (
	"io"
	"strings"
	"bytes"
	"strconv"
	"encoding/json"
)

type Response struct {
	Request
	Status  int
}

const (
	RESPONSE_STATE_LINE         = 0
	RESPONSE_STATE_SKIP_TO_BODY = 1
	RESPONSE_STATE_BODY         = 2
	RESPONSE_STATE_END          = 3
)

func ParseResponse(response []byte) (*Response, error) {
	resp := &Response{}

	state := RESPONSE_STATE_LINE
	reader := bytes.NewBuffer(response)
	//body_len := 0
	parse := true
	for parse {
		switch state {
		case RESPONSE_STATE_LINE:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				state = RESPONSE_STATE_END
				break
			}
			if err != nil {
				return nil, err
			}

			line = strings.TrimRight(line, "\r\n")
			splitted := strings.SplitN(line, " ", 3)
			if len(splitted) >= 2 {
				resp.Status, err = strconv.Atoi(splitted[1])
				if err != nil {
					return nil, err
				}
			}

			state = RESPONSE_STATE_SKIP_TO_BODY
		case RESPONSE_STATE_SKIP_TO_BODY:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				state = RESPONSE_STATE_END
				break
			}
			if err != nil {
				return nil, err
			}
			//if strings.Index(line, "Content-Length:") == 0 {
			//	body_len, _ = strconv.Atoi(strings.Trim(line[15:], " \r\n"))
			//}
			if line == "\r\n" || line == "\n" {
				state = RESPONSE_STATE_BODY
			}
		case RESPONSE_STATE_BODY:
			var buf bytes.Buffer
			io.Copy(&buf, reader)
			//wait := true
			//total := 0
			//for wait {
			//	written, _ := io.Copy(&buf, reader)
			//	total += int(written)
			//	fmt.Println(total, body_len)
			//	wait = body_len != 0 && total < body_len
			//}
			json.Unmarshal(buf.Bytes(), &resp.Request)
			state = RESPONSE_STATE_END
		case RESPONSE_STATE_END:
			parse = false
		}
	}
	return resp, nil
}
