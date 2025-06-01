package voronoi

import "errors"

// All kind of errors
var (
	ErrCircleSitesOrder              = errors.New("circle sites wrong order")
	ErrNoCircleFoundConnectionPoints = errors.New("no circle found connecting points")
)
