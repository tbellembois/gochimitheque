package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type CasbinMatchStoreLocationIsInEntityReq struct {
	CasbinMatchStoreLocationIsInEntity []any `json:"CasbinMatchStoreLocationIsInEntity"`
}

// string, int

// Response.
type CasbinMatchStoreLocationIsInEntityOk struct {
	Ok json.RawMessage
}
type CasbinMatchStoreLocationIsInEntityErr struct {
	Err string
}

func CasbinMatchStoreLocationIsInEntity(store_location_id int64, entity_id int64) (json.RawMessage, error) {
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
	req = append(req, store_location_id)
	req = append(req, entity_id)

	if message, err = json.Marshal(CasbinMatchStoreLocationIsInEntityReq{
		CasbinMatchStoreLocationIsInEntity: req,
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

			var resp CasbinMatchStoreLocationIsInEntityOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp CasbinMatchStoreLocationIsInEntityErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
