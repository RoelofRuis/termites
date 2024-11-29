package termites

import (
	"fmt"
)

type connectionConfig struct {
	from    *OutPort
	to      *InPort
	adapter *adapter
	mailbox MailboxConfig
}

type ConnectionOption func(conn *connectionConfig)

func Via[A any, B any](transform func(A) (B, error)) ConnectionOption {
	untypedTransform, inDataType, outDataType := extractFunc(transform)

	return func(conn *connectionConfig) {
		conn.adapter = &adapter{
			name:        fmt.Sprintf("As %s", outDataType.Name()),
			inDataType:  inDataType,
			outDataType: outDataType,
			transform:   untypedTransform,
		}
	}
}

func ViaNamed[A any, B any](transform func(A) (B, error), name string) ConnectionOption {
	untypedTransform, inDataType, outDataType := extractFunc(transform)

	return func(conn *connectionConfig) {
		conn.adapter = &adapter{
			name:        name,
			inDataType:  inDataType,
			outDataType: outDataType,
			transform:   untypedTransform,
		}
	}
}

func To(in *InPort) ConnectionOption {
	if in == nil {
		panic(fmt.Errorf("invalid connection option: in port cannot be nil"))
	}
	return func(conn *connectionConfig) {
		conn.to = in
	}
}

func WithMailbox(conf MailboxConfig) ConnectionOption {
	return func(conn *connectionConfig) {
		conn.mailbox = conf
	}
}
