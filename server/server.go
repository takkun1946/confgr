// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/elliottpolk/confgr/datastore"
)

const (
	DefaultStdPort string = "8080"
	DefaultTlsPort string = "8443"

	DefaultCertFile string = ".ssl/cert.pem"
	DefaultKeyFile  string = ".ssl/key.pem"
)

var stdPort, tlsPort, certFile, keyFile string

func Start() {
	if err := datastore.Start(); err != nil {
		panic(err)
	}
	fmt.Println("confgr datastore started")

	//  configure listener ports
	stdPort = DefaultStdPort
	if os.Getenv("HTTP_PORT") != "" {
		stdPort = os.Getenv("HTTP_PORT")
	}

	tlsPort = DefaultTlsPort
	if os.Getenv("HTTPS_PORT") != "" {
		tlsPort = os.Getenv("HTTPS_PORT")
	}

	//  configure ssl
	certFile = DefaultCertFile
	keyFile = DefaultKeyFile
	if os.Getenv("TLS_CERT") != "" {
		certFile = os.Getenv("TLS_CERT")
	}

	if os.Getenv("TLS_KEY") != "" {
		keyFile = os.Getenv("TLS_KEY")
	}

	fmt.Println("confgr server starting")

	if startHttps(certFile, keyFile) {
		fmt.Println("HTTPS started")
	}

	if err := http.ListenAndServe(":"+stdPort, nil); err != nil {
		fmt.Printf("unable to serve http: %v\n", err)
	}
}

func startHttps(certFile, keyFile string) bool {
	certInfo, certErr := os.Stat(certFile)
	if certErr != nil && !os.IsNotExist(certErr) {
		fmt.Printf("unable to access cert file %s: %v\n", certFile, certErr)
	}

	keyInfo, keyErr := os.Stat(keyFile)
	if keyErr != nil && !os.IsNotExist(keyErr) {
		fmt.Printf("unable to access key file %s: %v\n", keyFile, keyErr)
	}

	if certInfo != nil && keyInfo != nil {
		//  run HTTPS listener in goroutine to allow HTTP server
		go func(port, cert, key string) {
			if err := http.ListenAndServeTLS(":"+port, cert, key, nil); err != nil {
				fmt.Printf("unable to serve https: %v\n", err)
			}
		}(tlsPort, certFile, keyFile)

		return true
	}

	return false
}
