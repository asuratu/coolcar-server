package mongotesting

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"testing"
)

const (
	Image         = "mongo"
	ContainerPort = "27017/tcp"
)

// RunWithMongoInDocker runs the given test function with a mongo container
func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
	// 初始化数据库
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// 开启docker service端口，端口始终是27017，27017端口映射到Host端口
	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: Image,
		ExposedPorts: nat.PortSet{
			ContainerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			ContainerPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0", //0，系统分配端口
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	containerID := resp.ID
	defer func() {
		//time.Sleep(10 * time.Minute)

		// kill docker service
		err := c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()

	// 开启docker services
	err = c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	iPAndPort := inspRes.NetworkSettings.Ports["27017/tcp"][0]
	*mongoURI = fmt.Sprintf("mongodb://%s:%s/coolcar?readPreference=primary&ssl=false", iPAndPort.HostIP, iPAndPort.HostPort)

	return m.Run()
}
