/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mock

import (
	"context"
	"fmt"
	pluginDriver "github.com/brokercap/Bifrost/plugin/driver"
	"github.com/brokercap/Bifrost/sdk/pluginTestData"
	"time"
)

type NormalTable struct {
	SchemaName    string
	TableName     string
	LongStringLen int
	NoUint64      bool
	NoMapping     bool
	NoPks         bool
	TwoPks        bool
	ch            chan *pluginDriver.PluginDataType
}

func (t *NormalTable) GetSchemaName() string {
	return t.SchemaName
}

func (t *NormalTable) GetTableName() string {
	return t.TableName
}

func (t *NormalTable) Start(ctx context.Context, ch chan *pluginDriver.PluginDataType) {
	t.ch = ch
	event := pluginTestData.NewEvent()
	event.SetSchema(t.SchemaName)
	event.SetTable(t.TableName)
	event.SetLongStringLen(t.LongStringLen)
	event.SetNoUint64(t.NoUint64)
	t.Callback(event.GetTestInsertData())
	t.Callback(event.GetTestUpdateData(true))
	t.Callback(event.GetTestDeleteData())
	t.Callback(event.GetTestInsertData())
	t.Callback(event.GetTestUpdateData(true))
	t.Callback(event.GetTestInsertData())
	t.Callback(event.GetTestInsertData())

	timeDuration := 8 * time.Second
	timer := time.NewTimer(timeDuration)
	<-timer.C
	t.Callback(event.GetTestUpdateData(true))

	timer.Reset(timeDuration)
	<-timer.C

	t.Callback(event.GetTestDeleteData())
	t.Callback(event.GetTestInsertData())

}

func (t *NormalTable) Callback(data *pluginDriver.PluginDataType) {
	if t.ch == nil {
		return
	}
	if t.NoMapping {
		data.ColumnMapping = nil
	}
	if t.NoPks {
		data.Pri = make([]string, 0)
	}
	if t.TwoPks {
		secondPk := fmt.Sprintf("%s_mock_pk_2", data.Pri[0])
		data.Pri = append(data.Pri, secondPk)
		if !t.NoMapping {
			data.ColumnMapping[secondPk] = data.ColumnMapping[data.Pri[0]]
		}
		for k, _ := range data.Rows {
			i := k
			data.Rows[i][secondPk] = data.Rows[i][data.Pri[0]]
		}
	}
	t.ch <- data
}
