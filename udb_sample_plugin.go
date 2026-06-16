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

func (p *UdbSamplePlugin) GetBoardMap() map[string]types.Board {
	return map[string]types.Board{
		"single-colour":  boards.NewSingleColourBoard("single-colour"),
		"digital-clock":  boards.NewDigitalClockBoard("digital-clock"),
		"gradient":       boards.NewGradientBoard("gradient"),
		"rainbow":        boards.NewRainbowBoard("rainbow"),
		"sprite":         boards.NewSpriteBoard("sprite"),
		"gif":            boards.NewGifBoard("gif"),
		"scrolling-text": boards.NewScrollingTextBoard("scrolling-text"),
	}
}

func (p *UdbSamplePlugin) GetDatasourceMap() map[string]types.Datasource {
	return map[string]types.Datasource{
		"current-time": &datasources.CurrentTimeDatasource{},
	}
}
