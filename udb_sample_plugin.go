package udb_plugin_samples

import (
	udb_plugin_library "github.com/benwiebe/udb-plugin-library"
	"github.com/benwiebe/udb-plugin-library/types"

	"github.com/benwiebe/udb-plugin-samples/boards"
	"github.com/benwiebe/udb-plugin-samples/datasources"
)

func init() {
	udb_plugin_library.Register(&UdbSamplePlugin{})
}

type UdbSamplePlugin struct{}

func (p *UdbSamplePlugin) GetId() string                             { return "udb-plugin-samples" }
func (p *UdbSamplePlugin) GetName() string                           { return "UDB Sample Plugin" }
func (p *UdbSamplePlugin) Configure(config types.PluginConfig) error { return nil }

func (p *UdbSamplePlugin) GetBoardMap() map[string]types.BoardFactory {
	return map[string]types.BoardFactory{
		"single-colour":  func() types.Board { return boards.NewSingleColourBoard() },
		"digital-clock":  func() types.Board { return boards.NewDigitalClockBoard() },
		"gradient":       func() types.Board { return boards.NewGradientBoard() },
		"rainbow":        func() types.Board { return boards.NewRainbowBoard() },
		"sprite":         func() types.Board { return boards.NewSpriteBoard() },
		"gif":            func() types.Board { return boards.NewGifBoard() },
		"scrolling-text": func() types.Board { return boards.NewScrollingTextBoard() },
	}
}

func (p *UdbSamplePlugin) GetDatasourceMap() map[string]types.DatasourceFactory {
	return map[string]types.DatasourceFactory{
		"current-time": func() types.Datasource { return &datasources.CurrentTimeDatasource{} },
	}
}
