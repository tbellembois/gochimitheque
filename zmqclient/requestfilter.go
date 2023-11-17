package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type RequestFilterReq struct {
	RequestFilter string `json:"RequestFilter"`
}

// Response.
type RequestFilterOk struct {
	Ok RequestFilter
}
type RequestFilterErr struct {
	Err string
}

func RequestFilterFromRawString(req string) (RequestFilter, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(RequestFilterReq{
		RequestFilter: req,
	}); err != nil {
		return RequestFilter{}, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return RequestFilter{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp RequestFilterOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return RequestFilter{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp RequestFilterErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return RequestFilter{}, err
			}

			return RequestFilter{}, err

		}

	}

	return RequestFilter{}, nil

}
