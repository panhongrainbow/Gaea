package containerdTest

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainerdPath(t *testing.T) {
	r := Load{prefix: "./test/", client: nil}
	path := r.ContainerdPath()
	require.Equal(t, path, "test/containerd")
}

func TestListContainerD(t *testing.T) {
	r := Load{prefix: "./example/", client: nil}
	files, err := r.listContainerD()
	require.Nil(t, err)
	require.Equal(t, len(files), 2)
	require.Contains(t, files, "mariadb-sakila.json")
	require.Contains(t, files, "mysql.json")
}

func TestLoadContainerD(t *testing.T) {
	r := Load{prefix: "./example/", client: nil}
	sakila, err := r.loadContainerD("mariadb-sakila.json")
	require.Nil(t, err)
	require.Equal(t, sakila.Name, "mariadb-sakila")
}

func TestLoadAllContainerD(t *testing.T) {
	r := Load{prefix: "./example/", client: nil}
	configs, err := r.loadAllContainerD()
	require.Nil(t, err)
	for key, value := range configs {
		require.Equal(t, key, value.Name)
	}
}
