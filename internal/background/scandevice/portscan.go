package scandevice

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type portScanService struct {
	wg              sync.WaitGroup
	dialer          *net.Dialer
	mu              sync.Mutex
	portMin         int
	portMax         int
	caPath          string
	scannerKeyPath  string
	scannerCertPath string
	ipWhiteList     map[string]bool
	conns           map[string]*tls.Conn
	closed          bool
}

func NewPortScanService(
	portMin int,
	portMax int,
	caPath string,
	scannerKeyPath string,
	scannerCertPath string,
	ipWhiteList map[string]bool,
) DeviceScan {
	return &portScanService{
		dialer:          &net.Dialer{},
		mu:              sync.Mutex{},
		wg:              sync.WaitGroup{},
		conns:           make(map[string]*tls.Conn),
		portMin:         portMin,
		portMax:         portMax,
		caPath:          caPath,
		scannerKeyPath:  scannerKeyPath,
		scannerCertPath: scannerCertPath,
		ipWhiteList:     ipWhiteList,
		closed:          false,
	}
}

func (p *portScanService) ScanDevice(ctx context.Context) error {
	tlsConfig, err := p.loadTLSConfig()
	if err != nil {
		return fmt.Errorf("failed to load cert: %w", err)
	}

	for port := p.portMin; port <= p.portMax; port++ {
		log.Println("Scanning port", port)
		address := "127.0.0.1:" + strconv.Itoa(port)

		if p.checkWhitelist(address) {
			log.Printf("INFO: has been scanned and whitelisted: %s", address)
			continue
		}

		conn, err := p.connect(tlsConfig, address)
		if err != nil {
			log.Println("failed to connect: ", err)
			continue
		}

		p.register(ctx, address, conn)
		log.Println("port: ", port, " has whitelisted")
	}

	return nil
}

func (p *portScanService) register(ctx context.Context, addr string, conn *tls.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.closed {
		p.conns[addr] = conn
		p.ipWhiteList[addr] = true
		p.wg.Add(1)
		go p.ping(ctx, conn, addr)
	}
}

func (p *portScanService) release(addr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.conns, addr)
	delete(p.ipWhiteList, addr)
}

func (p *portScanService) checkWhitelist(addr string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, ok := p.ipWhiteList[addr]

	return ok
}

func (p *portScanService) loadTLSConfig() (*tls.Config, error) {
	return loadCert(p.scannerCertPath, p.scannerKeyPath, p.caPath, false)
}

func (p *portScanService) connect(tlsconfig *tls.Config, address string) (*tls.Conn, error) {
	conn, err := tls.DialWithDialer(p.dialer, "tcp", address, tlsconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return conn, err
}

func (p *portScanService) ping(ctx context.Context, conn *tls.Conn, addr string) error {
	defer p.wg.Done()
	errCh := make(chan error)

	defer conn.Close()
	defer p.release(addr)
	defer close(errCh)

	stop := context.AfterFunc(ctx, func() {
		conn.Close()
	})

	defer stop()

	err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)

	}

	_, err = conn.Write([]byte("localhost:8080\n"))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)

	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			return fmt.Errorf("failed to set write deadline: %w", err)
		}

		_, err = conn.Write([]byte("ping\n"))
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)

		}

		b := make([]byte, 5)
		err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		if err != nil {
			return fmt.Errorf("failed to set read deadline: %w", err)

		}

		_, err = io.ReadFull(conn, b)
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}

		if string(b) != "pong\n" {
			return fmt.Errorf("unexpected response: %s", string(b))
		}

		select {
		case <-time.After(5 * time.Second):
			continue
		case <-ctx.Done():
			log.Printf("shutdown ping [%s]....\n", addr)
			return ctx.Err()
		}
	}
}

func (p *portScanService) Run(ctx context.Context) error {
	for {
		p.ScanDevice(ctx)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
			continue
		}
	}
}

func (p *portScanService) Shutdown(ctx context.Context) error {
	p.mu.Lock()

	p.closed = true
	for _, conn := range p.conns {
		conn.Close()
	}
	p.mu.Unlock()

	// ensure ping goroutine returns
	done := make(chan struct{})
	go func() { p.wg.Wait(); close(done) }()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		log.Println("closing all worker connections...")
		return nil
	}
}
