package merged

import "github.com/ed3899/kumo/common/tool/interfaces"

type Merged[E interfaces.Environment] struct {
	General E
	Cloud   E
}
