package core

import (
	"log"
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
	CenterAgent CenterAgent
	MsgCh       chan Message
	ErrCh       chan error
	StopCh      chan struct{}
}

// NewClientService creates a ClientService instance
func NewClientService(client Client, centerAgent CenterAgent) *ClientService {
	return &ClientService{
		Client:      client,
		CenterAgent: centerAgent,
		MsgCh:       make(chan Message, 10),
		ErrCh:       make(chan error, 10),
		StopCh:      make(chan struct{}, 1),
	}
}

// Start is
func (s *ClientService) Start() error {
	// s.CenterAgent = NewCenterAgentImpl(s.Client.ID, s.Client.Name, s.ServerIP, s.ServerPort)
	// err := s.CenterAgent.Connect()
	// if err != nil {
	// 	log.Println(err) // todo
	// 	return err
	// }

	/*
	 * bind center-message handler and error handler
	 */
	s.CenterAgent.OnMessage(s.MsgCh)
	s.CenterAgent.OnError(s.ErrCh)

	/*
	 * register team
	 */
	reg := s.Client.GetRegistration()
	err := s.CenterAgent.Register(reg)
	if err != nil {
		log.Println(err) // todo
		return err
	}

	/*
	 * handle the coming messages
	 */
	go s.Handle()

	return nil
}

// Stop is
func (s *ClientService) Stop() {
	s.StopCh <- struct{}{}
	s.CenterAgent.OffMessage(s.MsgCh)
	s.CenterAgent.OffError(s.ErrCh)
}

// Handle is
func (s *ClientService) Handle() {
loop:
	for {
		select {
		case msg := <-s.MsgCh:
			switch msg.Name {
			case LegStartName:
				legStart, err := msg.LegStart()
				if err != nil {
					s.ErrCh <- err
					continue loop
				}
				if err := s.Client.LegStart(legStart); err != nil {
					s.ErrCh <- err
				}
			case LegEndName:
				legEnd, err := msg.LegEnd()
				if err != nil {
					s.ErrCh <- err
					continue loop
				}
				if err := s.Client.LegEnd(legEnd); err != nil {
					s.ErrCh <- err
				}
			case RoundName:
				round, err := msg.Round()
				if err != nil {
					s.ErrCh <- err
					continue loop
				}
				action, err := s.Client.Round(round)
				if err != nil {
					s.ErrCh <- err
					continue loop
				}
				err = s.CenterAgent.Act(action)
				if err != nil {
					s.ErrCh <- err
				}
			case GameOverName:
				_, err := msg.GameOver()
				if err != nil {
					s.ErrCh <- err
				}
				err = s.Client.GameOver()
				if err != nil {
					s.ErrCh <- err
				}
				go s.Stop()
			default:
				log.Printf("message error - unknown message %q with msg data:\n%v", msg.Name, msg.String())
			}
		case err := <-s.ErrCh:
			log.Println(err) // todo
		case <-s.StopCh:
			break loop
		}
	}
	// select {
	// case legStart := <-s.CenterAgent.LegStartCh:
	// 	if err := s.Client.LegStart(&legStart); err != nil {
	// 		s.CenterAgent.ErrCh <- err
	// 	}
	// case legEnd := <-s.CenterAgent.LegEndCh:
	// 	if err := s.Client.LegEnd(&legEnd); err != nil {
	// 		s.CenterAgent.ErrCh <- err
	// 	}
	// case round := <-s.CenterAgent.RoundCh:
	// 	action, err := s.Client.Round(&round)
	// 	if err != nil {
	// 		s.CenterAgent.ErrCh <- err
	// 		break
	// 	}
	// 	if err := s.CenterAgent.Act(action); err != nil {
	// 		s.CenterAgent.ErrCh <- err
	// 	}
	// case <-s.CenterAgent.GameOverCh:
	// 	if err := s.Client.GameOver(); err != nil {
	// 		s.CenterAgent.ErrCh <- err
	// 	}
	// 	break
	// case err := <-s.CenterAgent.ErrCh:
	// 	fmt.Printf("ERROR: %v", err) // todo
	// }
	// }
}
