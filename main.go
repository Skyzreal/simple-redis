package main

import (
	"log"
	"log/slog"
	"net"
)

// config for the server
type Config struct {
	ListenAddr string
}

// TCP server
type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
}

const defaultListenAddr = ":5001"

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	return s.acceptLoop()
}
func (s *Server) loop() {
	for {
		select {
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}

	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Accept error", "error", err)
			continue
		}
		go s.handleConn(conn)
	}
}
func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer

	peer.readLoop()
}
func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}
