![](https://i.imgur.com/BDQ2wWJ.gif)

# Resolv v0.6.0

[pkg.go.dev](https://pkg.go.dev/github.com/solarlune/resolv)

## What is Resolv?

Resolv is a 2D collision detection and resolution library, specifically created for simpler, arcade-y (non-realistic) video games. Resolv is written in Go, but the core concepts are fairly straightforward and could be easily adapted for use with other languages or game development frameworks.

Basically: It allows you to do simple physics easier, without actually doing the physics part - that's still on you.

## Why is it called that?

Because it's like... You know, collision resolution? To **resolve** a collision? So... That's the name. I juste took an e off because I seem to have misplaced it.

## Why did you create Resolv?

Because I was making games in Go and found that existing frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for simple, "arcade-y" games; if you need realistic physics, you have other options like Box2D or something.

____

As an aside, this actually used to be quite different; I decided to rework it when I was less than satisfied with my previous efforts, and after a few attempts over several months, I got it to this state (which, I think, is largely better). That said, there are breaking changes between the previous version, v0.4, and the current one (v0.5). These changes were necessary, in my opinion, to improve the library.

In comparison to the previous version of Resolv, v0.5 includes, among other things:

- A redesigned API from scratch. 
- The usage of floats instead of ints for position and movement, simplifying real usage of the library dramatically.
- Broadphase grid-based collision testing and querying for simple collisions, which means solid performance gains.
- ConvexPolygons for SAT intersection tests.

It's still a work-in-progress, but should be solid enough for usage in the field.

## How do I get it?

`go get github.com/solarlune/resolv`

## How do I use it?

There's a couple of ways to use Resolv. 

Firstly, you can create a Space, create Objects, add them to the Space, check the Space for collisions / intersections, and finally update the Objects as they move.

In Resolv v0.5, a Space represents a limited, bounded area which is separated into even Cells of predetermined size. Any Objects added to the Space fill at least one Cell (as long as the Object is within the Space). By checking a position in the Space, you can tell which Cells are occupied and so, where objects generally are. This is the broadphase, simpler portion of Resolv. 

Here's an example:

```go

var space *resolv.Space
var playerObj *resolv.Object

// As an example, in a game's initialization function that runs once
// when the game or level starts...

func Init() {

    // First, we want to create a Space. This represents the areas in our game world 
    // where objects can move about and check for collisions.

    // The first two arguments represent our Space's width and height, while the 
    // following two represent the individual Cells' sizes. The smaller the Cells' 
    // sizes, the finer the collision detection. Generally, each Cell should be 
    // reasonably be the size of a "unit", whatever that may be for the game.
    // For example, the player character, enemies, and collectibles could fit 
    // into one or more of these Cells.

    space = resolv.NewSpace(640, 480, 16, 16)

    // Next, we can start creating things and adding it to the Space.

    // Here's some level geometry. resolv.NewObject() takes the X and Y
    // position, and width and height to create a new *resolv.Object.
    // You can also specify tags when creating an Object. Tags can be used 
    // to filter down objects when checking the Space for a collision.

    space.Add(
        resolv.NewObject(0, 0, 640, 16),
        resolv.NewObject(0, 480-16, 640, 16),
        resolv.NewObject(0, 16, 16, 480-32),
        resolv.NewObject(640-16, 16, 16, 480-32),
    )

    // We'll keep a reference to the player's Object to move it later.
    playerObj = resolv.NewObject(32, 32, 16, 16)

    // Finally, we add the Object to the Space, and we're good to go!
    space.Add(playerObj)

}

// Later on, in the game's update loop, which runs once per game frame...

func Update() {

    // Let's say we are attempting to move the player to the right by 2 
    // pixels. Here's how we could do it.
    dx := 2.0

    // To start, we check to see if there would be a collision if the 
    // playerObj were to move to the right by 2 pixels. The Check function 
    // returns a Collision object if so.

    if collision := playerObj.Check(dx, 0); collision != nil {
        
        // If there was a collision, the "playerObj" Object can't move fully 
        // to the right by 2, and Object.Check() would return a *Collision object.
        // A *Collision object contains the Objects and Cells that the calling 
        // *resolv.Object ran into when it called Check().

        // To resolve (haha) this collision, we probably want to move the player into
        // contact with that Object. So, we call Collision.ContactWithObject() on the 
        // first Object that we came into contact with (which is stored in the Collision).

        // Collision.ContactWithObject() will return a Vector, indicating how much
        // distance to move to come into contact with the specified Object.

        // We could also come into contact with the cell to the right using 
        // Collision.ContactWithCell(collision.Cells[0]).
        dx = collision.ContactWithObject(collision.Objects[0]).X()

    }

    // If there wasn't a collision, then dx will just be 2, as set above, and the 
    // movement will go through unimpeded.
    playerObj.X += dx

    // Lastly, when we move an Object, we need to call Object.Update() so it can be 
    // updated within the Space as well. For static / unmoving Objects, this is
    // unnecessary, as Object.Update() is called once when an Object is first added to a Space.
    playerObj.Update()

    // If we were making a platformer, you could then check for the Y-axis as well. 
    // Conceptually, this is decomposing movement into two separate axes, and is a familiar 
    // and well-used approach for handling movement in a standard tile-based platformer. 
    // See this fantastic post on the subject:
    // http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

    // If you want to filter out types of Objects to check for, add tags on the objects 
    // you want to filter using Object.AddTags(), or when the Object is created 
    // with resolv.NewObject(), and specify them in Object.Check().

    onlySolidOrHazardous := playerObj.Check(dx, 0, "hazard", "solid")

}

// That's it!

```

The second way to use Resolv is to check for a more accurate shape intersection test by assigning two Objects Shapes, and then checking for an intersection between them. Checking for an intersection between Shapes internally performs separating axis theorum (SAT) collision testing (when checking against ConvexPolygons), and represents the more inefficient narrow-phase portion of Resolv. If you can get by without doing Shape-based intersection testing, it would be most performant to do so.

```go

playerObj *resolv.Object
stairs *resolv.Object
space *resolv.Space

func Init() {
    
    space = resolv.NewSpace(640, 480, 16, 16)

    // Create the Object as usual, but then...
    playerObj = resolv.NewObject(32, 128, 16, 16)
    // Assign the Object a Shape. A Rectangle is, for now, a ConvexPolygon that's simply 
    // rectangular, rather than a specific, separate Shape.
    playerObj.SetShape(resolv.NewRectangle(0, 0, 16, 16))
    // Then we add the Object to the Space.
    space.Add(playerObj)

    // Note that we can just use the shapes directly as well.

    stairs = resolv.NewObject(96, 128, 16, 16)

    // Here, we use resolv.NewConvexPolygon() to create a new ConvexPolygon Shape. It takes 
    // a series of float64 values indicating the X and Y positions of each vertex; the call 
    // below, for example, creates a triangle.

    stairs.SetShape(resolv.NewConvexPolygon(
        0, 0, // Position of the polygon

        16, 0, // (x, y) pair for the first vertex
        16, 16, // (x, y) pair for the second vertex
        0, 16, // (x, y) pair for the third and last vertex
    ))

    //     0
    //    /|
    //   / |
    //  /  |
    // 2---1

    // Note that the vertices are in clockwise order. They can be in either clockwise or 
    // counter-clockwise order as long as it's consistent throughout your application. 
    // As an aside, resolv.NewRectangle() defines the vertices in clockwise order.
    space.Add(stairs)

}

func Update() {

    dx := 1.0

    // Shape.Intersection() returns a ContactSet, representing information 
    // regarding the intersection between two Shapes (i.e. the point(s) of
    // collision, the distance to move to get out, etc).
    if intersection := playerObj.Shape.Intersection(dx, 0, stairs.Shape); intersection != nil {
        
        // We are colliding with the stairs shape, so we can move according
        // to the delta (MTV) to get out of it.
        dx = intersection.MTV.X()

        // You might want to move a bit less (say, 0.1) than the delta to
        // avoid "bouncing", depending on your application.

    }

    playerObj.X += dx

    // When Object.Update() is called, the Object's Shape is also moved
    // accordingly.
    playerObj.Update()

}

```

___

Welp, that's about it. If you want to see more info, feel free to examine the examples in the `examples` folder; the `worldPlatformer.go` example is particularly in-depth when it comes to movement and collision response. You can run them by just calling `go run .` from the examples folder.

[You can check out the documentation here, as well.](https://pkg.go.dev/github.com/solarlune/resolv)

## Dependencies?

Resolv requires just quartercastle's nice and clean [vector](https://github.com/quartercastle/vector) library, and the built-in `math` package. For the examples, [ebiten](https://github.com/hajimehoshi/ebiten), as well as tanema's [gween](https://github.com/tanema/gween) library are also required.

## Shout-out Time!

Thanks to the people who stopped by on my [YouTube Golang gamedev streams](https://www.youtube.com/c/SolarLune) - they helped out a lot with a couple of the technical aspects of getting Go to do what I needed to, haha.

If you want to support development, feel free to throw me a couple of bones over on my [Patreon](https://www.patreon.com/SolarLune) or [itch.io](https://solarlune.itch.io/) / [Steam](https://store.steampowered.com/app/1269310/MasterPlan/) pages. I really appreciate it - thanks!

## To-do List

- [ ] Allow for cells that are less than 1 unit large (and so Spaces can have cell sizes of, say, 0.1 units)
- [ ] Custom Vector struct for speed, consistency, and to reduce third-party imports
- [ ] Intersection MTV works properly for external normals, but not internal normals of a polygon