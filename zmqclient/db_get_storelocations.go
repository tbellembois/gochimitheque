package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type DBGetStorelocationsReq struct {
	DBGetStorelocations []interface{} `json:"DBGetStorelocations"`
}

// string, int

// Response.
type DBGetStorelocationsOk struct {
	Ok json.RawMessage
}
type DBGetStorelocationsErr struct {
	Err string
}

func DBGetStorelocations(requestRawString string, person_id int) (json.RawMessage, error) {
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return json.RawMessage{}, err
	}
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return json.RawMessage{}, err
	}

	var (
		message []byte
	)

	req := make([]interface{}, 0)
	req = append(req, requestRawString)
	req = append(req, person_id)

	if message, err = json.Marshal(DBGetStorelocationsReq{
		DBGetStorelocations: req,
	}); err != nil {
		return json.RawMessage{}, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return json.RawMessage{}, err
	}

	if msg, err := s.Recv(0); err != nil {
		return json.RawMessage{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp DBGetStorelocationsOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp DBGetStorelocationsErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
