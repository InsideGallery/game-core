package scene

import (
	"context"

	"github.com/FrogoAI/memory/registry"
	"github.com/InsideGallery/game-core/client/camera"
	"github.com/InsideGallery/game-core/client/core"
	"github.com/InsideGallery/game-core/client/event"
	"github.com/InsideGallery/game-core/geometry/shapes"
	"github.com/InsideGallery/game-core/rtree"
	"github.com/hajimehoshi/ebiten/v2"
)

// BaseScene provides shared ECS infrastructure for all scenes.
// Embed this in concrete scenes to get Systems, Registry, RTree, Camera, Bus, and a World buffer.
type BaseScene struct {
	Ctx      context.Context
	Registry *registry.Registry[string, uint64, any]
	Systems  *core.Systems
	RTree    *rtree.RTree
	Camera   *camera.Camera
	Bus      *event.Bus
	World    *ebiten.Image // offscreen World buffer; systems draw here in world coords
	vpW, vpH int           // current viewport dimensions; guards World reallocation
}

// NewBaseScene creates a BaseScene with all ECS infrastructure initialized.
// vpW and vpH set the camera viewport and the World buffer size; pass 0 to skip World creation.
func NewBaseScene(ctx context.Context, bus *event.Bus, vpW, vpH int) *BaseScene {
	cam := camera.NewCamera(shapes.NewPoint())
	cam.SetViewPort(vpW, vpH)

	var world *ebiten.Image
	if vpW > 0 && vpH > 0 {
		world = ebiten.NewImage(vpW, vpH)
	}

	return &BaseScene{
		Ctx:      ctx,
		Systems:  core.NewSystems(),
		Registry: registry.NewRegistry[string, uint64, any](),
		RTree:    rtree.NewRTree(rtree.DefaultMinRTreeOption, rtree.DefaultMaxRTreeOption),
		Bus:      bus,
		Camera:   cam,
		World:    world,
		vpW:      vpW,
		vpH:      vpH,
	}
}

// viewportNeedsResize reports whether the World buffer should be recreated.
// Extracted as a pure helper so it can be tested without an *ebiten.Image.
func viewportNeedsResize(curW, curH, newW, newH int, hasWorld bool) bool {
	return !hasWorld || newW != curW || newH != curH
}

// SetViewPort resizes the camera viewport and recreates the World buffer only when dimensions change.
func (b *BaseScene) SetViewPort(w, h int) {
	b.Camera.SetViewPort(w, h)
	if w > 0 && h > 0 && viewportNeedsResize(b.vpW, b.vpH, w, h, b.World != nil) {
		b.World = ebiten.NewImage(w, h)
		b.vpW = w
		b.vpH = h
	}
}

// Update iterates all systems in order.
func (b *BaseScene) Update() error {
	for _, sys := range b.Systems.Get() {
		if err := sys.Update(b.Ctx); err != nil {
			return err
		}
	}

	return nil
}

// Draw composites world-space systems into the World buffer (via camera matrix),
// then calls ScreenDraw on SystemWindow systems for screen-space UI overlays.
// If World is nil, systems draw directly to screen without the camera transform.
func (b *BaseScene) Draw(screen *ebiten.Image) {
	systems := b.Systems.Get()

	if b.World != nil {
		b.World.Clear()

		for _, sys := range systems {
			sys.Draw(b.Ctx, b.World)
		}

		screen.DrawImage(b.World, &ebiten.DrawImageOptions{
			GeoM: b.Camera.WorldMatrix(),
		})
	} else {
		for _, sys := range systems {
			sys.Draw(b.Ctx, screen)
		}
	}

	for _, sys := range systems {
		if w, ok := sys.(core.SystemWindow); ok {
			w.ScreenDraw(b.Ctx, screen)
		}
	}
}
