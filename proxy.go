package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type UDPProxy struct {
	localAddr    string
	remoteAddr   string
	listener     *net.UDPConn
	sessions     map[string]*session
	sessionsLock sync.RWMutex
	stopChan     chan struct{}
	running      bool
	runningLock  sync.RWMutex
}

type session struct {
	clientAddr *net.UDPAddr
	remoteConn *net.UDPConn
	lastActive time.Time
}

func NewUDPProxy(localAddr, remoteAddr string) *UDPProxy {
	return &UDPProxy{
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		sessions:   make(map[string]*session),
	}
}

func (p *UDPProxy) Start() error {
	p.runningLock.Lock()
	if p.running {
		p.runningLock.Unlock()
		return fmt.Errorf("proxy is already running")
	}
	p.running = true
	p.runningLock.Unlock()

	// resolve remote
	remoteUDPAddr, err := net.ResolveUDPAddr("udp4", p.remoteAddr)
	if err != nil {
		p.runningLock.Lock()
		p.running = false
		p.runningLock.Unlock()
		return fmt.Errorf("failed to resolve remote address: %v", err)
	}
	log(fmt.Sprintf("Resolved remote address: %s", remoteUDPAddr.String()))

	// local listen
	localUDPAddr, err := net.ResolveUDPAddr("udp4", p.localAddr)
	if err != nil {
		p.runningLock.Lock()
		p.running = false
		p.runningLock.Unlock()
		return fmt.Errorf("failed to resolve local address: %v", err)
	}

	listener, err := net.ListenUDP("udp4", localUDPAddr)
	if err != nil {
		p.runningLock.Lock()
		p.running = false
		p.runningLock.Unlock()
		return fmt.Errorf("failed to bind to port %s: %v", p.localAddr, err)
	}
	p.listener = listener
	p.stopChan = make(chan struct{})

	log(fmt.Sprintf("Started forwarding %s to %s (UDP)", p.localAddr, p.remoteAddr))

	go p.handleIncoming(remoteUDPAddr)

	return nil
}

func (p *UDPProxy) Stop() {
	p.runningLock.Lock()
	if !p.running {
		p.runningLock.Unlock()
		return
	}
	p.running = false
	p.runningLock.Unlock()

	close(p.stopChan)

	if p.listener != nil {
		p.listener.Close()
	}

	p.sessionsLock.Lock()
	for _, sess := range p.sessions {
		if sess.remoteConn != nil {
			sess.remoteConn.Close()
		}
	}
	p.sessions = make(map[string]*session)
	p.sessionsLock.Unlock()

	log("Stopped forwarding")
}

func (p *UDPProxy) IsRunning() bool {
	p.runningLock.RLock()
	defer p.runningLock.RUnlock()
	return p.running
}

func (p *UDPProxy) handleIncoming(remoteAddr *net.UDPAddr) {
	buffer := make([]byte, 65535)

	for {
		select {
		case <-p.stopChan:
			return
		default:
			p.listener.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, clientAddr, err := p.listener.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				select {
				case <-p.stopChan:
					return
				default:
					log(fmt.Sprintf("Error reading from client: %v", err))
					continue
				}
			}

			sess := p.getOrCreateSession(clientAddr, remoteAddr)
			if sess == nil {
				continue
			}

			// remote forwarding
			go func(data []byte, size int, s *session) {
				_, err := s.remoteConn.Write(data[:size])
				if err != nil {
					log(fmt.Sprintf("Error forwarding to remote: %v", err))
				}
			}(append([]byte{}, buffer[:n]...), n, sess)
		}
	}
}

func (p *UDPProxy) getOrCreateSession(clientAddr *net.UDPAddr, remoteAddr *net.UDPAddr) *session {
	key := clientAddr.String()

	p.sessionsLock.RLock()
	sess, exists := p.sessions[key]
	p.sessionsLock.RUnlock()

	if exists {
		sess.lastActive = time.Now()
		return sess
	}

	// Create new session
	remoteConn, err := net.DialUDP("udp4", nil, remoteAddr)
	if err != nil {
		log(fmt.Sprintf("Error creating connection to remote: %v", err))
		return nil
	}

	sess = &session{
		clientAddr: clientAddr,
		remoteConn: remoteConn,
		lastActive: time.Now(),
	}

	p.sessionsLock.Lock()
	p.sessions[key] = sess
	p.sessionsLock.Unlock()

	go p.handleRemoteResponse(sess)

	return sess
}

func (p *UDPProxy) handleRemoteResponse(sess *session) {
	buffer := make([]byte, 65535)

	for {
		sess.remoteConn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := sess.remoteConn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// exit if we don't exist anymore
				p.sessionsLock.RLock()
				_, exists := p.sessions[sess.clientAddr.String()]
				p.sessionsLock.RUnlock()
				if !exists {
					return
				}
				continue
			}
			return
		}

		// back to client
		_, err = p.listener.WriteToUDP(buffer[:n], sess.clientAddr)
		if err != nil {
			log(fmt.Sprintf("Error forwarding to client: %v", err))
		}

		sess.lastActive = time.Now()
	}
}
