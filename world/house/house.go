package house

import "github.com/mikelangelon/town-sweet-town/common"

type House struct {
	ID       string
	Position common.Position
	Owner    *string
}
