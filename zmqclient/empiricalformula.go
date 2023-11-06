package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type EmpiricalFormula struct {
	EmpiricalFormula string `json:"EmpiricalFormula"`
}

// Response.
type EmpiricalFormulaOk struct {
	Ok string
}
type EmpiricalFormulaErr struct {
	Err string
}

func Empirical_formula(req string) (string, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(EmpiricalFormula{
		EmpiricalFormula: req,
	}); err != nil {
		return "", err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return "", err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp EmpiricalFormulaOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return "", err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp EmpiricalFormulaErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return "", err
			}

			return "", err

		}

	}

	return "", nil

}
