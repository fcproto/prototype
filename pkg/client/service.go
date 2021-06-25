package client

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fcproto/prototype/pkg/api"
	"github.com/sirupsen/logrus"
)

const Timeout = time.Second * 3

type Service struct {
	endpoint string
	log      *logrus.Logger
	client   *http.Client
	ClientID string
	NearCars []*api.SensorData

	bufferMutex sync.Mutex
	bufferPos   int
	buffer      []*api.SensorData
}

func NewService(log *logrus.Logger, endpoint string, bufferSize int) (*Service, error) {
	svc := &Service{
		endpoint: endpoint,
		log:      log,
		client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   Timeout,
				MaxIdleConns:          10,
				IdleConnTimeout:       Timeout,
				ResponseHeaderTimeout: Timeout,
			},
			Timeout: Timeout,
		},
		bufferPos: 0,
		buffer:    make([]*api.SensorData, bufferSize),
	}

	if err := svc.loadOrGenerateId(); err != nil {
		return nil, err
	}
	if err := svc.readFromFile(); err != nil {
		return nil, err
	}

	return svc, nil
}

func createRandomId() []byte {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		panic(err)
	}
	return id
}

func (s *Service) loadOrGenerateId() error {
	clientIdFile := "client.id"
	if clientIdFileEnv := os.Getenv("CLIENT_ID_FILE"); clientIdFileEnv != "" {
		clientIdFile = clientIdFileEnv
	}

	rawClientId, err := os.ReadFile(clientIdFile)
	if err != nil {
		s.log.Warn(err)
		rawClientId = createRandomId()
		if err := os.WriteFile(clientIdFile, rawClientId, 0644); err != nil {
			return err
		}
	}
	s.ClientID = hex.EncodeToString(rawClientId)
	return nil
}

func (s *Service) incPos(pos int) int {
	pos++
	if pos == len(s.buffer) {
		pos = 0
	}
	return pos
}

func (s *Service) decPos(pos int) int {
	pos--
	if pos < 0 {
		pos += len(s.buffer)
	}
	return pos
}

func (s *Service) SubmitSensorData(data *api.SensorData) error {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()
	if s.buffer[s.bufferPos] != nil {
		s.log.Warn("buffer is full, overwriting old data")
	}
	data.ClientID = s.ClientID
	s.buffer[s.bufferPos] = data
	s.bufferPos = s.incPos(s.bufferPos)
	if err := s.writeToFile(); err != nil {
		return err
	}
	return nil
}

func (s *Service) syncUp() error {
	return s.GetSensorData(func(data []*api.SensorData) error {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(data)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", s.endpoint, &buf)
		if err != nil {
			return err
		}
		res, err := s.client.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("invalid status code %d", res.StatusCode)
		}
		return nil
	})
}

func (s *Service) syncDown() error {
	req, err := http.NewRequest("GET", s.endpoint+"/near/"+s.ClientID, nil)
	if err != nil {
		return err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("invalid status code %d", res.StatusCode)
	}

	var nearCars []*api.SensorData
	err = json.NewDecoder(res.Body).Decode(&nearCars)
	if err != nil {
		return err
	}
	s.NearCars = nearCars
	return nil
}

func (s *Service) Sync() error {
	if err := s.syncUp(); err != nil {
		return err
	}
	if err := s.syncDown(); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetSensorData(fn func([]*api.SensorData) error) error {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	data := make([]*api.SensorData, 0)
	readPos := s.decPos(s.bufferPos)
	for i := 0; i < len(s.buffer); i++ {
		if el := s.buffer[readPos]; el != nil {
			data = append(data, el)
		}
		readPos = s.decPos(readPos)
	}
	err := fn(data)
	if err != nil {
		return err
	}

	// txn was ok, free buffer
	for i := 0; i < len(s.buffer); i++ {
		s.buffer[i] = nil
	}
	s.bufferPos = 0
	return nil
}

func (s *Service) writeToFile() error {
	data, err := json.Marshal(s.buffer)
	if err != nil {
		return err
	}
	err = os.WriteFile("service-data.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) readFromFile() error {
	data, err := os.ReadFile("service-data.json")
	if err != nil {
		s.log.Warn(err)
		return nil
	}

	var buffer []*api.SensorData
	if err := json.Unmarshal(data, &buffer); err != nil {
		s.log.Warn(err)
		return nil
	}

	if len(buffer) != len(s.buffer) {
		s.log.Warn("ignoring invalid backup file")
		return nil
	}

	for _, entry := range buffer {
		if entry == nil {
			break
		}
		s.buffer[s.bufferPos] = entry
		s.bufferPos = s.incPos(s.bufferPos)
	}

	return nil
}

func (s *Service) String() string {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()
	return fmt.Sprintf("{id=%s, pos=%d, buffer=%v}", s.ClientID, s.bufferPos, s.buffer)
}
