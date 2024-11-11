
<p align="center">
<img src="logo.png">
</p>

# Resolv v0.8.0

[pkg.go.dev](https://pkg.go.dev/github.com/solarlune/resolv)

## What is Resolv?

Resolv is a 2D collision detection and resolution library, specifically created for simpler, arcade-y (non-realistic) video games. Resolv is written in pure Go, but the core concepts are fairly straightforward and could be easily adapted for use with other languages or game development frameworks.

Basically: It allows you to do simple physics easier, without actually _doing_ the physics part - that's still on you and your game's use-case.

## Why is it called that?

Because it's like... You know, collision resolution? To **resolve** a collision? So... That's the name. I juste seem to have misplaced the "e", so I couldn't include it in the name - how odd.

## Why did you create Resolv?

Because I was making games in Go and found that existing frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for simple, "arcade-y" games; if you need realistic physics, you have other options like [cp](https://github.com/jakecoffman/cp) or [Box2D](https://github.com/ByteArena/box2d).

____

As an aside, this actually used to be quite different; I decided to rework it a couple of times. This is now the second rework, and should be _significantly_ easier to use and more accurate. (Thanks a lot to everyone who contributed their time to submit PRs and issues!)

It's still not _totally_ complete, but it should be solid enough for usage in the field.

## Dependencies?

Resolv has no external dependencies. It requires Go 1.20 or above.

## How do I get it?

`go get github.com/solarlune/resolv`

## How do I use it?

There's a couple of ways to use Resolv. One way is to just create Shapes then use functions to check for intersections.

```go
func main() {

    // Create a rectangle at 200, 100 with a width and height of 32x32
    rect := resolv.NewRectangle(200, 100, 32, 32)

    // Create a circle at 200, 120 with a radius of 8
    circle := resolv.NewCircle(200, 120, 8)

    // Check for intersection
    if intersection, ok := rect.Intersection(circle); ok {
        fmt.Println("They're touching! Here's the data:", intersection)
    }

}
```

You can also get the intersection with `Shape.Intersection(other)`.

However, you'll probably want to check intersection with a larger group of objects, which you can do with `Spaces` and `ShapeFilters`. You create a Space, add Shapes to the space, and then call `Shape.IntersectionTest()` with more advanced settings:

```go

type Game struct {
    Rect *resolv.ConvexPolygon
    Space *resolv.Space
}

func (g *Game) Init() {

    // Create a space that is 640x480 large and that has a cellular size of 16x16. The cell size is mainly used to
    // determine internally how close objects are together to qualify for intersection testing. Generally, this should
    // be the size of the maximum speed of your objects (i.e. objects shouldn't move faster than 1 cell in size each
    // frame).
    g.Space = resolv.NewSpace(640, 480, 16, 16)

    // Create a rectangle at 200, 100 with a width and height of 32x32
    g.Rect = resolv.NewRectangle(200, 100, 32, 32)

    // Create a circle at 200, 120 with a radius of 8
    circle := resolv.NewCircle(200, 120, 8)

    // Add the shapes to allow them to be detected by other Shapes.
    g.Space.Add(rect)
    g.Space.Add(circle)
}

func (g *Game) Update() {

    // Check for intersection and do something for each intersection
    g.Space.Rect.IntersectionTest(resolv.IntersectionTestSettings{
        TestAgainst: rect.SelectTouchingCells(1).FilterShapes(), // Check only shapes that are near the rectangle (within 1 cell's margin)
        OnIntersect: func(set resolv.IntersectionSet, index, max int) bool {
            fmt.Println("There was an intersection with some other object! Here's the data:", set)
            return true
        }
    })

}

```

You can also do line tests and shape-based line tests, to see if there would be a collision in a given direction - this is more-so useful for movement and space checking.
___

If you want to see more info, feel free to examine the examples in the `examples` folder; the platformer example is particularly in-depth when it comes to movement and collision response. You can run them and switch between them just by calling `go run .` from the examples folder and pressing Q or E to switch between the example worlds.

[You can check out the documentation here, as well.](https://pkg.go.dev/github.com/solarlune/resolv)

## To-do List

- [x] Rewrite to be significantly easier and simpler
- [ ] Allow for cells that are less than 1 unit large (and so Spaces can have cell sizes of, say, 0.1 units)
- [x] Custom Vector struct for speed, consistency, and to reduce third-party imports
    - [ ] Implement Matrices as well for parenting?
- [ ] Intersection MTV works properly for external normals, but not internal normals of a polygon
- [ ] Properly implement moving around inside a circle (?)