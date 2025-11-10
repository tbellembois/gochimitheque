package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type CasbinMatchPersonIsInPersonEntityReq struct {
	CasbinMatchPersonIsInPersonEntity []any `json:"CasbinMatchPersonIsInPersonEntity"`
}

// string, int

// Response.
type CasbinMatchPersonIsInPersonEntityOk struct {
	Ok json.RawMessage
}
type CasbinMatchPersonIsInPersonEntityErr struct {
	Err string
}

func CasbinMatchPersonIsInPersonEntity(person_id int64, other_person_id int64) (json.RawMessage, error) {
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

	req := make([]any, 0)
	req = append(req, person_id)
	req = append(req, other_person_id)

	if message, err = json.Marshal(CasbinMatchPersonIsInPersonEntityReq{
		CasbinMatchPersonIsInPersonEntity: req,
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

			var resp CasbinMatchPersonIsInPersonEntityOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp CasbinMatchPersonIsInPersonEntityErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
