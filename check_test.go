package resourcemultiply_test

import (
	"testing"

	resourcemultiply "github.com/gmidorii/resource-multiply"
)

func TestCheck(t *testing.T) {

	// スキーマをコピーして指定した数作成
	err, cleanup := resourcemultiply.MultiplySchema("hoge", 3)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		err := cleanup()
		if err != nil {
			t.Error(err)
		}
	})

	// スキーマごとのコネクションを作成
}
