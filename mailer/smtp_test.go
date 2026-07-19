package mailer

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net"
// 	"net/http"
// 	"net/mail"
// 	"os"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/dankobg/fluffly-api/config"
// 	"github.com/ory/dockertest/v3"
// 	"github.com/ory/dockertest/v3/docker"
// 	"github.com/rs/zerolog"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// type IntegrationMailerSuite struct {
// 	suite.Suite

// 	mailer   *smtpClient
// 	resource *dockertest.Resource
// }

// func TestIntegrationMailerSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationMailerSuite))
// }

// func (s *IntegrationMailerSuite) SetupSuite() {
// 	resource, err := setupMailhogResource(s.T())
// 	s.NoError(err)

// 	mailer, err := setupMailer(resource)
// 	s.NoError(err)

// 	s.mailer = mailer
// 	s.resource = resource
// }

// func (s *IntegrationMailerSuite) TestIntegrationMailer() {

// 	if testing.Short() {
// 		s.T().Skip("skipping integration tests for mailer")
// 	}

// 	s.Run("send email", func() {
// 		//

// 		md := MailData{
// 			plainTextContent: "test plain",
// 			htmlContent:      "<h1>test html</h1>",
// 			subject:          "test",
// 			replyTo:          "test@test.com",
// 			fromAddress:      mail.Address{Name: "fluffly", Address: "office@fluffly.xyz"},
// 			toAddress:        mail.Address{Name: "fluffly", Address: "office@fluffly.xyz"},
// 		}

// 		err := s.mailer.Send(context.Background(), &md)
// 		s.NoError(err)
// 	})

// 	s.Run("confirm email delivered", func() {
// 		//

// 		url := fmt.Sprintf("http://localhost:%s/api/v2/messages", s.resource.GetPort("8025/tcp"))

// 		resp, err := http.Get(url)
// 		s.NoError(err)

// 		s.T().Cleanup(func() {
// 			resp.Body.Close()
// 		})

// 		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
// 			s.FailNow("bad status code")
// 		}

// 		bb, err := io.ReadAll(resp.Body)
// 		s.NoError(err)

// 		type mailhogMessagesResult struct {
// 			Total int           `json:"total"`
// 			Count int           `json:"count"`
// 			Start int           `json:"start"`
// 			Items []interface{} `json:"items"`
// 		}

// 		var result mailhogMessagesResult
// 		uerr := json.Unmarshal(bb, &result)
// 		s.NoError(uerr)

// 		if result.Total == 0 || len(result.Items) == 0 {
// 			s.FailNow("0 total emails received")
// 		}
// 	})
// }

// func setupMailhogResource(tb testing.TB) (*dockertest.Resource, error) {
// 	asserts := assert.New(tb)

// 	pool, err := dockertest.NewPool("")
// 	asserts.NoError(err)

// 	pool.MaxWait = 10 * time.Second

// 	options := &dockertest.RunOptions{
// 		Repository: "mailhog/mailhog",
// 		Tag:        "latest",
// 	}

// 	hostConfigOptions := []func(hc *docker.HostConfig){
// 		func(hc *docker.HostConfig) {
// 			hc.AutoRemove = true
// 			hc.RestartPolicy = docker.RestartPolicy{
// 				Name: "no",
// 			}
// 		},
// 	}

// 	resource, err := pool.RunWithOptions(options, hostConfigOptions...)
// 	asserts.NoError(err)

// 	rerr := pool.Retry(func() error {
// 		smtpPort := resource.GetPort("1025/tcp")
// 		addr := net.JoinHostPort("localhost", smtpPort)
// 		_, err := net.Dial("tcp", addr)
// 		return err
// 	})
// 	asserts.NoError(rerr)

// 	asserts.NoError(resource.Expire(30))

// 	tb.Cleanup(func() {
// 		asserts.NoError(pool.Purge(resource))
// 	})

// 	return resource, nil
// }

// func setupMailer(resource *dockertest.Resource) (*smtpClient, error) {
// 	devSMTPHost := resource.GetBoundIP("1025/tcp")
// 	devSMTPPort, err := strconv.Atoi(resource.GetPort("1025/tcp"))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse resource mailhog port: %w", err)
// 	}

// 	cfg := &config.EmailSettings{
// 		Enabled:         false,
// 		TLS:             false,
// 		FromName:        "fluffly",
// 		FromAddress:     "office@fluffly.xyz",
// 		DevSMTPHost:     devSMTPHost,
// 		DevSMTPPort:     devSMTPPort,
// 		DevSMTPUsername: "test",
// 		DevSMTPPassword: "test",
// 	}
// 	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

// 	client, err := NewSmtpClient(cfg, &logger)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client, nil
// }
