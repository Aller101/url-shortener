package random

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}

	for _, res := range tests {
		t.Run(res.name, func(t *testing.T) {
			nRS1 := NewRandomString(res.size)

			//sleep - был добавлен, потому что функции генерят
			//одну и ту же строку, т к проходит слижком мало времени,
			//между вызовами этих функция

			//а без нее результат nRS1 и nRS2 были бы равны(не всегда)
			//и тест провалился бы
			time.Sleep(1 * time.Microsecond)
			nRS2 := NewRandomString(res.size)

			assert.Len(t, nRS1, res.size)
			assert.Len(t, nRS2, res.size)

			assert.NotEqual(t, nRS1, nRS2)
		})
	}

}
