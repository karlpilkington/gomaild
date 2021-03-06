//Package cipher provides cryptographic functionality, such as TLS.
package cipher

import (
	"crypto/rand"
	"crypto/tls"
	"github.com/trapped/gomaild/config"
	"log"
	"net"
)

var TLSAvailable bool     //Whether TLS certificates have been loaded successfully
var TLSConfig *tls.Config //Static configuration for TLS

//Loads the TLS certificate provided in the configuration file.
func TLSLoadCertificate() {
	cert, err := tls.LoadX509KeyPair(config.Configuration.TLS.CertificateFile,
		config.Configuration.TLS.CertificateKeyFile)
	if err != nil {
		log.Println("TLS:", "Failed loading SSL certificate:", err)
		return
	}
	TLSAvailable = true
	TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		Rand:         rand.Reader,
	}
}

//Converts a normal connection to a TLS-protected one, keeping the object type (net.Conn).
func TLSTransmuteConn(c net.Conn) net.Conn {
	tc := tls.Server(c, TLSConfig)
	err := tc.Handshake()
	if err != nil {
		log.Println("TLS:", "Error handshaking for", c.RemoteAddr().String()+":", err)
	}
	log.Println("TLS:", "Handshake completed for", c.RemoteAddr().String()+"; success:", tc.ConnectionState().HandshakeComplete)
	return tc
}
