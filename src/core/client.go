package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const defaultBufferSize = 1024 * 10
const defaultConnRetries = 30

// Client is
type Client struct {
	ID       int
	Name     string
	Strategy ClientStrategy
}

// NewClient creates a Client instance
func NewClient(id int, name string, strategy ClientStrategy) *Client {
	return &Client{
		ID:       id,
		Name:     name,
		Strategy: strategy,
	}
}

// GetRegistration is
func (c *Client) GetRegistration() *Registration {
	return &Registration{c.ID, c.Name}
}

// LegStart is
func (c *Client) LegStart(legStart *LegStart) error {
	return c.Strategy.LegStart(legStart)
}

// LegEnd is
func (c *Client) LegEnd(legEnd *LegEnd) error {
	return c.Strategy.LegEnd(legEnd)
}

// Round is
func (c *Client) Round(round *Round) (*Action, error) {
	return c.Strategy.Round(round)
}

// GameOver is
func (c *Client) GameOver() error {
	return c.Strategy.GameOver()
}

// ClientService is
type ClientService struct {
	Client      Client
	ServerPort  int
	ServerIP    string
	CenterAgent *CenterAgentStub
}

// NewClientService creates a ClientService instance
func NewClientService(client Client, ip string, port int) *ClientService {
	return &ClientService{
		Client:     client,
		ServerIP:   ip,
		ServerPort: port,
	}
}

// Connect is
func (s *ClientService) Connect() error {
	s.CenterAgent = NewCenterAgentStub(s.Client.ID, s.Client.Name, s.ServerIP, s.ServerPort)
	err := s.CenterAgent.Connect()
	if err != nil {
		log.Println(err) // todo
		return err
	}
	reg := s.Client.GetRegistration()
	err = s.CenterAgent.Register(reg)
	if err != nil {
		log.Println(err) // todo
		return err
	}

	go s.Bind()

	return nil
}

// Bind is
func (s *ClientService) Bind() {
	for {
		select {
		case legStart := <-s.CenterAgent.LegStartChan:
			if err := s.Client.LegStart(&legStart); err != nil {
				s.CenterAgent.ErrChan <- err
			}
		case legEnd := <-s.CenterAgent.LegEndChan:
			if err := s.Client.LegEnd(&legEnd); err != nil {
				s.CenterAgent.ErrChan <- err
			}
		case round := <-s.CenterAgent.RoundChan:
			action, err := s.Client.Round(&round)
			if err != nil {
				s.CenterAgent.ErrChan <- err
				break
			}
			if err := s.CenterAgent.Act(action); err != nil {
				s.CenterAgent.ErrChan <- err
			}
		case <-s.CenterAgent.GameOverChan:
			if err := s.Client.GameOver(); err != nil {
				s.CenterAgent.ErrChan <- err
			}
			break
		case err := <-s.CenterAgent.ErrChan:
			fmt.Printf("ERROR: %v", err) // todo
		}
	}
}

// CenterAgentStub is
type CenterAgentStub struct {
	TeamID       int
	TeamName     string
	ServerIP     string
	ServerPort   int
	Conn         net.Conn
	Connected    bool
	DialRetries  int
	BufferSize   int
	LegStartChan chan LegStart
	LegEndChan   chan LegEnd
	RoundChan    chan Round
	GameOverChan chan struct{}
	ErrChan      chan error
	reader       *bufio.Reader
	writer       *bufio.Writer
}

// NewCenterAgentStub is
func NewCenterAgentStub(teamID int, teamName string, serverIP string, serverPort int) *CenterAgentStub {
	return &CenterAgentStub{
		TeamID:     teamID,
		TeamName:   teamName,
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}
}

// Connect is
func (s *CenterAgentStub) Connect() (err error) {
	address := fmt.Sprintf("%s:%d", s.ServerIP, s.ServerPort)
	teamDesc := fmt.Sprintf("%v:%v", s.TeamID, s.TeamName)
	retries := s.DialRetries
	if retries <= 0 {
		retries = defaultConnRetries
	}

	log.Printf("team (%v) client is connecting to game server@%v", teamDesc, address)
	for i := 1; i <= retries; i++ {
		s.Conn, err = net.DialTimeout("tcp4", address, time.Second*1)
		if err != nil {
			fmt.Printf("client dial error, try %vth time, error: %v\n", i, err)
		} else {
			log.Printf("team (%v) client is connected to game server@%v", teamDesc, address)
			s.init()
			s.Connected = true
			return nil
		}
	}
	errMsg := fmt.Sprintf("team (%v) client is connected to game server@%v, error: %v", teamDesc, address, err)
	log.Println(errMsg)
	return errors.New(errMsg)
}

// Disconnect is
func (s *CenterAgentStub) Disconnect() (err error) {
	address := fmt.Sprintf("%s:%d", s.ServerIP, s.ServerPort)
	teamDesc := fmt.Sprintf("%v:%v", s.TeamID, s.TeamName)
	if !s.Connected {
		log.Printf("team (%v) client is not connected to game server@%v, no need to disconnect", teamDesc, address)
		return nil
	}

	if err = s.Conn.Close(); err != nil {
		log.Printf("team (%v) client fails to disconnect server@%v", teamDesc, address)
		return err
	}
	s.Connected = false
	s.uninit()
	log.Printf("team (%v) client is disconnected to game server@%v", teamDesc, address)

	go s.receive()

	return nil
}

func (s *CenterAgentStub) init() {
	if s.BufferSize <= 0 {
		s.BufferSize = defaultBufferSize
	}
	s.reader = bufio.NewReaderSize(s.Conn, s.BufferSize)
	s.writer = bufio.NewWriterSize(s.Conn, s.BufferSize)
	s.LegStartChan = make(chan LegStart, 1)
	s.LegEndChan = make(chan LegEnd, 1)
	s.RoundChan = make(chan Round, 1)
	s.GameOverChan = make(chan struct{}, 1)
	s.ErrChan = make(chan error, 10)
}

// todo
func (s *CenterAgentStub) uninit() {
}

func (s *CenterAgentStub) receive() {
	scanner := bufio.NewScanner(s.reader)
	scanner.Split(splitMsg)
	for scanner.Scan() {
		var msgData json.RawMessage
		msg := &APIMsg{Data: &msgData}
		err := json.Unmarshal(scanner.Bytes(), msg)
		if err != nil {
			log.Printf("fail to unmarshal msg, error: %v", err)
			continue
		}
		fmt.Println("RECV->", string(scanner.Bytes()))

		switch msg.Name {
		case "leg_start":
			legStart := new(LegStart)
			err := json.Unmarshal(msgData, legStart)
			if err != nil {
				s.ErrChan <- err
				continue
			}
			s.LegStartChan <- *legStart
		case "round":
			round := new(Round)
			err := json.Unmarshal(msgData, round)
			if err != nil {
				s.ErrChan <- err
				continue
			}
			s.RoundChan <- *round
		case "leg_end":
			legEnd := new(LegEnd)
			err := json.Unmarshal(msgData, legEnd)
			if err != nil {
				s.ErrChan <- err
				continue
			}
			s.LegEndChan <- *legEnd
		case "game_over":
			s.GameOverChan <- struct{}{}
		default:
			s.ErrChan <- fmt.Errorf("Unknown msg name: %v with msg data:\n %v", msg.Name, msgData)
		}
	}
}

// Register is
func (s *CenterAgentStub) Register(registration *Registration) error {
	msg := new(APIMsg)
	msg.Name = "registration"
	msg.Data = &registration
	return sendMsg(s.writer, msg)
}

// Act is
func (s *CenterAgentStub) Act(action *Action) error {
	msg := new(APIMsg)
	msg.Name = "action"
	msg.Data = &action
	return sendMsg(s.writer, msg)
}
