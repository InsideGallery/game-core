package relations

import (
	"testing"

	"github.com/FrogoAI/testutils"
	"github.com/InsideGallery/core/ecs"
)

// ---------- SetParent / RemParent / GetParentID tests ----------

func TestSetParent(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	rc.SetParent("typeA", 42)
	testutils.Equal(t, rc.GetParentID("typeA"), uint64(42))
}

func TestRemParent(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{"typeA": 42}, false)
	rc.RemParent("typeA")
	testutils.Equal(t, rc.GetParentID("typeA"), uint64(0))
}

func TestGetParentIDNonExistent(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	testutils.Equal(t, rc.GetParentID("noSuchType"), uint64(0))
}

// ---------- GetParents returns copy ----------

func TestGetParentsReturnsCopy(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{"typeA": 1, "typeB": 2}, false)
	parents := rc.GetParents()
	testutils.Equal(t, len(parents), 2)
	testutils.Equal(t, parents["typeA"], uint64(1))
	testutils.Equal(t, parents["typeB"], uint64(2))

	// Mutating the returned map should not affect the component
	parents["typeC"] = 3
	testutils.Equal(t, rc.GetParentID("typeC"), uint64(0))
}

// ---------- Attach / Detach / GetChildren / GetAllChildren tests ----------

type mockChild struct {
	*ecs.BaseEntity
	id uint64
}

func (m *mockChild) Parent(_ string) (Parent, error) { return nil, nil }
func (m *mockChild) SetParent(_ string, _ uint64)    {}
func (m *mockChild) RemParent(_ string)              {}
func (m *mockChild) Construct() error                { return nil }
func (m *mockChild) Destroy() error                  { return nil }

func newTestBaseEntityWithID(id uint64) *ecs.BaseEntity {
	return ecs.NewRegistry().NewBaseEntityWithID(id)
}

func newMockChild(id uint64) *mockChild {
	return &mockChild{
		BaseEntity: newTestBaseEntityWithID(id),
		id:         id,
	}
}

func TestAttachAndGetChildren(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	c1 := newMockChild(10)
	c2 := newMockChild(11)

	rc.Attach("child", c1)
	rc.Attach("child", c2)

	children := rc.GetChildren("child")
	testutils.Equal(t, len(children), 2)
}

func TestDetach(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	c1 := newMockChild(10)
	c2 := newMockChild(11)

	rc.Attach("child", c1)
	rc.Attach("child", c2)
	rc.Detach("child", c1)

	children := rc.GetChildren("child")
	testutils.Equal(t, len(children), 1)
	testutils.Equal(t, children[0], c2)
}

func TestDetachNonExistentChild(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	c1 := newMockChild(10)
	c2 := newMockChild(11)

	rc.Attach("child", c1)
	// Detaching c2 which was never attached should be a no-op
	rc.Detach("child", c2)

	children := rc.GetChildren("child")
	testutils.Equal(t, len(children), 1)
}

func TestGetChildrenNonExistentType(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	children := rc.GetChildren("nonexistent")
	testutils.Equal(t, len(children), 0)
}

func TestGetAllChildren(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	c1 := newMockChild(10)
	c2 := newMockChild(11)

	rc.Attach("typeA", c1)
	rc.Attach("typeB", c2)

	all := rc.GetAllChildren()
	testutils.Equal(t, len(all), 2)
	testutils.Equal(t, len(all["typeA"]), 1)
	testutils.Equal(t, len(all["typeB"]), 1)
}

func TestGetAllChildrenReturnsCopy(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	c1 := newMockChild(10)
	rc.Attach("typeA", c1)

	all := rc.GetAllChildren()
	// Mutate returned map
	all["typeA"] = append(all["typeA"], newMockChild(99))

	// Original should not change
	testutils.Equal(t, len(rc.GetChildren("typeA")), 1)
}

func TestGetChildrenReturnsCopy(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	c1 := newMockChild(10)
	rc.Attach("typeA", c1)

	children := rc.GetChildren("typeA")
	children = append(children, newMockChild(99))
	_ = children
	testutils.Equal(t, len(rc.GetChildren("typeA")), 1)
}

// ---------- Parent not found error ----------

func TestParentNotFound(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{"typeA": 0}, false)
	_, err := rc.Parent("typeA")
	testutils.Equal(t, err, ErrNotFoundParent)
}

func TestParentNotFoundMissingKey(t *testing.T) {
	rc := NewRelationComponent(map[string]uint64{}, false)
	_, err := rc.Parent("typeA")
	testutils.Equal(t, err, ErrNotFoundParent)
}
