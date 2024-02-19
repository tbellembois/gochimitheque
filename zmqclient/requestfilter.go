package zmqclient

import (
	"encoding/json"
	"errors"

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
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return RequestFilter{}, err
	}
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return RequestFilter{}, err
	}

	var (
		message []byte
	)

	if message, err = json.Marshal(RequestFilterReq{
		RequestFilter: req,
	}); err != nil {
		return RequestFilter{}, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return RequestFilter{}, err
	}

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

			return RequestFilter{}, errors.New(resp.Err)

		}

	}

	return RequestFilter{}, nil

}
