/*
 *Copyright ClusterHQ Inc.  See LICENSE file for details.
 *
 */

package logger

import (
	"log"
	"io"
)

var (
	Message *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

func Init( 
	messageHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Message = log.New(messageHandle,
		" ", 0)

	Info = log.New(infoHandle,
		"INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"[error] ", 0)
}

