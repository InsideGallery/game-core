package relations

import (
	"testing"

	"github.com/InsideGallery/core/ecs"
	"github.com/InsideGallery/core/testutils"
)

var (
	exampleParentKey = "0"
	exampleChildKey  = "1"
)

type ExampleParent struct {
	*ecs.BaseEntity
	*RelationComponent
}

func NewExampleParent() *ExampleParent {
	p := &ExampleParent{
		BaseEntity:        ecs.NewBaseEntityWithID(1),
		RelationComponent: NewRelationComponent(map[string]uint64{}, false),
	}
	err := store.Add(exampleParentKey, p.GetID(), p)
	if err != nil {
		panic(err)
	}
	return p
}

type ExampleChild struct {
	*ecs.BaseEntity
	*RelationComponent
}

func NewExampleChild() *ExampleChild {
	c := &ExampleChild{
		BaseEntity: ecs.NewBaseEntityWithID(2),
		RelationComponent: NewRelationComponent(map[string]uint64{
			exampleParentKey: 1,
		}, false),
	}
	err := store.Add(exampleChildKey, c.GetID(), c)
	if err != nil {
		panic(err)
	}
	return c
}

func (e *ExampleChild) Construct() error {
	return e.ConstructChild(exampleChildKey, e)
}

func (e *ExampleChild) Destroy() error {
	return e.DestroyChild(exampleChildKey, e)
}

func TestRelations(t *testing.T) {
	p := NewExampleParent()
	c := NewExampleChild()
	testutils.Equal(t, c.GetParentID(exampleParentKey), uint64(1))
	ch := p.GetChildren(exampleChildKey)
	testutils.Equal(t, len(ch), 1)
	testutils.Equal(t, ch[0], c)

	err := store.Remove(exampleChildKey, uint64(2))
	if err != nil {
		t.Fatal(err)
	}

	ch = p.GetChildren(exampleChildKey)
	testutils.Equal(t, len(ch), 0)
}
