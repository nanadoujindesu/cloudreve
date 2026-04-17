package migrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudreve/Cloudreve/v4/application/migrator/model"
	"github.com/cloudreve/Cloudreve/v4/ent/node"
	"github.com/cloudreve/Cloudreve/v4/inventory/types"
	"github.com/cloudreve/Cloudreve/v4/pkg/boolset"
)

func (m *Migrator) migrateNode() error {
	m.l.Info("Migrating nodes...")

	var nodes []model.Node
	if err := model.DB.Find(&nodes).Error; err != nil {
		return fmt.Errorf("failed to list v3 nodes: %w", err)
	}

	for _, n := range nodes {
		nodeType := node.TypeSlave
		nodeStatus := node.StatusSuspended
		if n.Type == model.MasterNodeType {
			nodeType = node.TypeMaster
		}
		if n.Status == model.NodeActive {
			nodeStatus = node.StatusActive
		}

		cap := &boolset.BooleanSet{}
		settings := &types.NodeSetting{
			Provider: types.DownloaderProviderA2Rpc,
		}

		if n.A2RpcEnabled {
			boolset.Sets(map[types.NodeCapability]bool{
				types.NodeCapabilityRemoteDownload: true,
			}, cap)

			a2rpcOptions := &model.A2RpcOption{}
			if err := json.Unmarshal([]byte(n.A2RpcOptions), a2rpcOptions); err != nil {
				return fmt.Errorf("failed to unmarshal a2 options: %w", err)
			}

			downloaderOptions := map[string]any{}
			if a2rpcOptions.Options != "" {
				if err := json.Unmarshal([]byte(a2rpcOptions.Options), &downloaderOptions); err != nil {
					return fmt.Errorf("failed to unmarshal a2 options: %w", err)
				}
			}

			settings.A2RpcSetting = &types.A2RpcSetting{
				Server:   a2rpcOptions.Server,
				Token:    a2rpcOptions.Token,
				Options:  downloaderOptions,
				TempPath: a2rpcOptions.TempPath,
			}
		}

		if n.Type == model.MasterNodeType {
			boolset.Sets(map[types.NodeCapability]bool{
				types.NodeCapabilityExtractArchive: true,
				types.NodeCapabilityCreateArchive:  true,
			}, cap)
		}

		stm := m.v4client.Node.Create().
			SetRawID(int(n.ID)).
			SetCreatedAt(formatTime(n.CreatedAt)).
			SetUpdatedAt(formatTime(n.UpdatedAt)).
			SetName(n.Name).
			SetType(nodeType).
			SetStatus(nodeStatus).
			SetServer(n.Server).
			SetSlaveKey(n.SlaveKey).
			SetCapabilities(cap).
			SetSettings(settings).
			SetWeight(n.Rank)

		if err := stm.Exec(context.Background()); err != nil {
			return fmt.Errorf("failed to create node %q: %w", n.Name, err)
		}

	}

	return nil
}
