package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fcproto/prototype/pkg/api"
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/ipfs/go-log/v2"
)

const Timeout = time.Second * 10

type Service struct {
	endpoint string
	log      *log.ZapEventLogger
	client   *http.Client

	bufferMutex sync.Mutex
	bufferPos   int
	buffer      []*api.SensorData
}

func NewService(endpoint string, bufferSize int) *Service {
	return &Service{
		endpoint: endpoint,
		log:      logger.New("service"),
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

func (s *Service) SubmitSensorData(data *api.SensorData) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()
	if s.buffer[s.bufferPos] != nil {
		s.log.Warn("buffer is full, overwriting old data")
	}
	s.buffer[s.bufferPos] = data
	s.bufferPos = s.incPos(s.bufferPos)
}

func (s *Service) SyncUp() error {
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

func (s *Service) String() string {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()
	return fmt.Sprintf("{pos=%d, buffer=%v}", s.bufferPos, s.buffer)
}
