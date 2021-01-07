package amqp

import (
  "github.com/streadway/amqp"
)

// Session  一个amqp会话
type Session struct {
  DSN     string
  Conn    *amqp.Connection
  ErrChan chan *amqp.Error
}

// NewSession 得到amqp会话对象
func NewSession(dsn string) (*Session, error) {
  session := &Session{DSN: dsn}
  err := session.connect()
  if err != nil {
    return nil, err
  }
  return session, nil
}

// Connect 建立amqp会话
func (s *Session) connect() error {
  errChan := make(chan *amqp.Error)
  s.ErrChan = errChan
  conn, err := amqp.Dial(s.DSN)
  if err != nil {
    return err
  }
  go func() {
    errs := conn.NotifyClose(make(chan *amqp.Error))
    for err := range errs {
      s.ErrChan <- err
    }
  }()

  s.Conn = conn
  return nil
}

// Close 关闭会话
func (s *Session) Close() error {
  return s.Conn.Close()
}
