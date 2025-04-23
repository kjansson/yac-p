package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/kjansson/yac-p/pkg/types"
)

func TestMetricsProcessing(t *testing.T) {

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l := &SlogLogger{}

	l.Init(types.Config{
		Debug: true,
	})

	// c := controller.Controller{
	// 	Config: types.Config{
	// 		//Debug:            true,
	// 	},
	// }

	// c.Init()
	// c := controller.NewController(types.Config{
	// 	Debug: true,
	// })

	// c := &controller.Controller{
	// 	Logger:     &SlogLogger{},
	// 	Collector:  &test_utils.YaceMockClient{},
	// 	YaceConfig: &yace.YaceOptions{},
	// 	Persister:  &prom.PromClient{},
	// 	Config: types.Config{
	// 		// RemoteWriteURL:   "http://localhost:1234",
	// 		// ConfigFileLoader: test_utils.GetTestConfigLoader(),

	// 		Debug: true,
	// 	},
	// }

	l.Log("info", "test message", "key1", "value1")

	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// // back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	// reading our temp stdout
	fmt.Println("previous output:")
	fmt.Print(out)

}
