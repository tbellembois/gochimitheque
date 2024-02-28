package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type PubchemGetCompoundByNameReq struct {
	PubchemGetCompoundByName string `json:"PubchemGetCompoundByName"`
}

// Response.
type PubchemGetCompoundByNameOk struct {
	Ok Compounds
}
type PubchemGetCompoundByNameErr struct {
	Err string
}

func PubchemGetCompoundByName(req string) (Compounds, error) {
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return Compounds{}, err
	}
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return Compounds{}, err
	}

	var (
		message []byte
	)

	if message, err = json.Marshal(PubchemGetCompoundByNameReq{
		PubchemGetCompoundByName: req,
	}); err != nil {
		return Compounds{}, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return Compounds{}, err
	}

	if msg, err := s.Recv(0); err != nil {
		return Compounds{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp PubchemGetCompoundByNameOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Compounds{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp PubchemGetCompoundByNameErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Compounds{}, err
			}

			return Compounds{}, errors.New(resp.Err)

		}

	}

	return Compounds{}, nil

}