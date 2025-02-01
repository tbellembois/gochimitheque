package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type DBCreateUpdateProducerReq struct {
	DBCreateUpdateProducer json.RawMessage `json:"DBCreateUpdateProducer"`
}

// Response.
type DBCreateUpdateProducerOk struct {
	Ok json.RawMessage
}
type DBCreateUpdateProducerErr struct {
	Err string
}

func DBCreateUpdateProducer(ProducertRawString json.RawMessage) (json.RawMessage, error) {
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

	if message, err = json.Marshal(DBCreateUpdateProducerReq{
		DBCreateUpdateProducer: ProducertRawString,
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

			var resp DBCreateUpdateProducerOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp DBCreateUpdateProducerErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
