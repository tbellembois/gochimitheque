package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type CasbinMatchStorageIsInEntityReq struct {
	CasbinMatchStorageIsInEntity []any `json:"CasbinMatchStorageIsInEntity"`
}

// string, int

// Response.
type CasbinMatchStorageIsInEntityOk struct {
	Ok json.RawMessage
}
type CasbinMatchStorageIsInEntityErr struct {
	Err string
}

func CasbinMatchStorageIsInEntity(storage_id int64, entity_id int64) (json.RawMessage, error) {
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
	req = append(req, storage_id)
	req = append(req, entity_id)

	if message, err = json.Marshal(CasbinMatchStorageIsInEntityReq{
		CasbinMatchStorageIsInEntity: req,
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

			var resp CasbinMatchStorageIsInEntityOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp CasbinMatchStorageIsInEntityErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
