// Package anon allows users to send anonymous messages to rooms via bot.
package anon

import (
	"github.com/matrix-org/go-neb/types"
	"github.com/matrix-org/gomatrix"
)

const ServiceType = "anon"

type Service struct {
	types.DefaultService
}

func (s *Service) Commands(client *gomatrix.Client) []types.Command {
	return []types.Command{
		types.Command{
			Path: []string{"anon"},
			Command: func(roomID, userID string, args []string) (interface{}, error) {
				return s.send(client, roomID, userID, args)
			},
		},
	}
}

func (s *Service) send(client *gomatrix.Client, roomID, userID string, args []string) (interface{}, error) {
	return nil, nil
}

func init() {
	types.RegisterService(func(serviceID, serviceUserID, webhookEndpointURL string) types.Service {
		return &Service{
			DefaultService: types.NewDefaultService(serviceID, serviceUserID, ServiceType),
		}
	})
}
