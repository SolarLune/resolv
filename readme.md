
# Resolv v0.5

[pkg.go.dev](https://pkg.go.dev/github.com/solarlune/resolv/resolv?tab=doc)

## What is Resolv?

Resolv is a 2D collision detection and resolution library, specifically created for simple, arcade (non-realistic) collision detection and resolution for video games. Resolv is written in Go, but the core concepts are fairly straightforward and could be easily adapted for use with other languages or game development frameworks.

Basically: It allows you to do simple physics easier, without actually doing the physics part - that's still on you.

## Why is it called that?

Because it's like... You know, collision resolution? To **resolve** a collision? So... That's the name. I juste took an e off because I seem to have misplaced it.

## Why did you create resolv?

Because I was making games in Go and found that existing frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for simple, "arcade-y" games; if you need realistic physics, you have other options like Box2D or something.

As an aside, this actually used to be quite different; I decided to rework it to its current incarnation, and after a few attempts over several months, I got it to this state (which, I think, is largely better). That said, there are breaking changes between v0.4.0 and v0.5.0; this was necessary, in my opinion, to update and improve the library.

In comparison to the previous version of resolv, v0.5.0 includes:

- An API that has been completely reworked from scratch. It now uses floats instead of ints for position and movement, simplifying usage of the library dramatically.
- Addition of broadphase grid-based collision testing and querying for simple collisions, which means performance gains.
- A reimplementation of SAT intersection-based intersection tests.
- Addition of ConvexPolygons for intersection tests.

That said, it's still a work-in-progress, but should be solid enough for usage in the field.

## How do I get it?

`go get github.com/solarlune/resolv`

## How do I use it?

There's a couple of ways to use Resolv. 

Firstly, you can create a Space, create Objects, update the Objects as they move, and check for collisions / intersections. 

In Resolv, a Space represents a limited, bounded area in which Objects are placed, and which is separated into even Cells of predetermined size. Each Object fills at least one Cell as long as it exists within the Space, and by checking its position against the Space, it can tell which Cells are occupied, and therefore, where other objects are. This is the broadphase portion of Resolv.

```go

var space *resolv.Space
var playerObj *resolv.Object

// As an example, in a game's initialization function that runs once when the game or level starts...

func Init() {

    // First, we want to create a Space. This represents the areas in our game world 
    // where objects can move about and check for collisions.

    // The first two arguments represent our Space's width and height, while the following two
    // represent the cells' sizes. The smaller the cells' sizes, the finer (and less efficient) the collision detection.
    space = resolv.NewSpace(640, 480, 16, 16)

    // Next, we can start creating things and adding it to the Space.

    // Here's some level geometry.
    // NewObject takes the X and Y position, and width and height of the object.
    space.Add(
        resolv.NewObject(0, 0, 640, 16),
        resolv.NewObject(0, 480-16, 640, 16),
        resolv.NewObject(0, 16, 16, 480-16),
        resolv.NewObject(640-16, 16, 16, 480-16),
    )

    // We'll keep a reference to the player's Object to move it later.
    playerObj = resolv.NewObject(32, 32, 16, 16, space)

    // Finally, we add the Object to the Space, and we're good to go!
    space.Add(playerObj)

}

// Later on, in the game's update loop, which runs once per game frame...

func Update() {

    // Let's say we are attempting to move the player to the right by 2 pixels. Here's how we could do it.
    dx := 2.0

    // To start, we check to see if there would be a collision if playerObj were to move to the right by 2 pixels. The Check function returns
    // a Collision object if so.

    if collision := playerObj.Check(dx, 0); collision != nil {
        
        // If there was a collision, the player's Object can't move fully to the right by 2.

        // To resolve (haha) this collision, we probably want to move the player into contact with that Object. So, we call Collision.ContactToObject() on the first
        // Object that we came into contact with (which is stored in the Collision). It will return a vector.Vector, which indicates how much distance to move in 
        // to come into contact with the specified Object.

        // We could also come into contact with the cell to the right using Collision.ContactWithCell(collision.Cells[0]).
        dx = collision.ContactWithObject(collision.Objects[0]).X()

    }

    // If there wasn't a collision, then dx will just be 2, as set above, and the movement will go through unimpeded.
    playerObj.X += dx

    // Lastly, when we move an Object, we need to call Object.Update() so it can be updated within the Space as well. For static / unmoving Objects, this is
    // unnecessary, as Object.Update() is called once when an Object is first added to a Space.
    playerObj.Update()

    // If we were making a platformer, you could then check for the Y-axis as well - conceptually, this is decomposing movement into two separate axes,
    // and is a familiar and well-used approach for handling movement in a standard tile-based platformer. See this fantastic post on the subject:
    // http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

    // If you want to filter out types of Objects to check for, add tags on the objects you want to filter using Object.AddTags(), or when the Object is created 
    // with resolv.NewObject(), and specify them in Object.Check.

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
    // Assign the Object a Shape. A Rectangle is, for now, a ConvexPolygon that's simply rectangular.
    playerObj.Shape = resolv.NewRectangle(0, 0, 16, 16)
    // Then we add the Object to the Space. Note that it's important that you do this last so that the Shape is properly updated. (You can also simply call Object.Update() later.)
    space.Add(playerObj)

    stairs = resolv.NewObject(96, 128, 16, 16)
    // NewConvexPolygon() takes a series of float64 values indicating the X and Y positions of each vertex; the call below, for example, creates a triangle.
    stairs.Shape = resolv.NewConvexPolygon(
        16, 0, // (x, y) for the first vertex
        16, 16, // (x, y) for the second vertex
        0, 16, // (x, y) for the last vertex
    )
    // Note that the vertices are in clockwise order. They can be in whatever order as long as it's consistent throughout your application. NewRectangle specifies
    // the vertices in clockwise order.
    space.Add(stairs)

}

func Update() {

    dx := 1.0

    // Shape.Intersection() returns a ContactSet, representing the intersection between two Shapes (i.e. the point(s) of collision, the normal, the distance to move to get out, etc).
    if intersection := playerObj.Shape.Intersection(dx, 0, stairs.Shape); intersection != nil {
        
        // We are colliding with the stairs shape, so we can move according to the delta to get out of it.
        dx = intersection.MTV.X()

        // You might want to move a bit less (say, 0.1) than the delta to avoid "bouncing", depending on your application.

    }

    playerObj.X += dx

    // When Object.Update() is called, the Object's Shape is also moved accordingly.
    playerObj.Update()

}

```
### Quick question, is it inefficient to have a ton of Objects in a Space?

Checking an object against a Space just operates using cellular positions, which are either filled, or not, with Objects. Having a lot of Objects doesn't really influence this process, as having more Objects doesn't complicate the collision checking any more than having fewer Objects, assuming they're sharing Cells.

----

Welp, that's about it. If you want to see more info, feel free to examine the main.go and world#.go tests to see how a couple of quick example tests are set up.

[You can check out the documentation here, as well.](https://pkg.go.dev/github.com/solarlune/resolv/resolv?tab=doc)

## Dependencies?

Resolv requires just kvartborg's nice and clean [vector](github.com/kvartborg/vector) library, and the built-in `fmt` and `math` packages. For the tests, [ebiten](github.com/hajimehoshi/ebiten) is also required.

## Shout-out Time!

Thanks to the people who stopped by on my [YouTube Golang gamedev streams](https://www.youtube.com/c/SolarLune) - they helped out a lot with a couple of the technical aspects of getting Go to do what I needed to, haha.
