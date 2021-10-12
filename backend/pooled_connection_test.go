package backend

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseMaster(t *testing.T) {
	s := new(Slice)
	s.Cfg = models.Slice{
		Name:     "slice-0",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3309",
		Slaves: []string{
			"192.168.122.2:3310",
			"192.168.122.2:3311",
		},
		Capacity:    12,
		MaxCapacity: 24,
		IdleTimeout: 60,
	}
	s.charset = "utf8mb4"
	s.collationID = 46
	masterStr := "192.168.122.2:3309"
	err := s.ParseMaster(masterStr)
	require.Equal(t, err, nil)
}
