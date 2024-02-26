package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type DBGetProducerrefsReq struct {
	DBGetProducerrefs string `json:"DBGetProducerrefs"`
}

// Response.
type DBGetProducerrefsOk struct {
	Ok json.RawMessage
}
type DBGetProducerrefsErr struct {
	Err string
}

func DBGetProducerrefs(requestRawString string) (json.RawMessage, error) {
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

	if message, err = json.Marshal(DBGetProducerrefsReq{
		DBGetProducerrefs: requestRawString,
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

			var resp DBGetProducerrefsOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp DBGetProducerrefsErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
