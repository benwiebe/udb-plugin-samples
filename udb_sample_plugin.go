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
func (p *UdbSamplePlugin) GetPluginType() types.PluginType           { return types.PluginTypeCombined }
func (p *UdbSamplePlugin) Configure(config types.PluginConfig) error { return nil }

func (p *UdbSamplePlugin) GetBoardMap() map[string]types.Board[any] {
	return map[string]types.Board[any]{
		"single-colour":  boards.NewSingleColourBoard("single-colour"),
		"digital-clock":  boards.NewDigitalClockBoard("digital-clock"),
		"gradient":       boards.NewGradientBoard("gradient"),
		"rainbow":        boards.NewRainbowBoard("rainbow"),
		"sprite":         boards.NewSpriteBoard("sprite"),
		"gif":            boards.NewGifBoard("gif"),
		"scrolling-text": boards.NewScrollingTextBoard("scrolling-text"),
	}
}

func (p *UdbSamplePlugin) GetAllBoards() []types.Board[any] {
	m := p.GetBoardMap()
	result := make([]types.Board[any], 0, len(m))
	for _, b := range m {
		result = append(result, b)
	}
	return result
}

func (p *UdbSamplePlugin) GetDatasourceMap() map[string]types.Datasource[any] {
	return map[string]types.Datasource[any]{
		"current-time": &datasources.CurrentTimeDatasource{},
	}
}

func (p *UdbSamplePlugin) GetAllDatasources() []types.Datasource[any] {
	m := p.GetDatasourceMap()
	result := make([]types.Datasource[any], 0, len(m))
	for _, ds := range m {
		result = append(result, ds)
	}
	return result
}
