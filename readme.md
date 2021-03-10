# Resolv

![Bouncing](https://thumbs.gfycat.com/UnhappySadEquine-size_restricted.gif)

[Other gifs](https://gfycat.com/@solarlune/collections/LxBK2KGU/resolv-gfys)

[Docs](https://pkg.go.dev/github.com/solarlune/resolv)

## What is Resolv?

Resolv is a library specifically created for simple arcade (non-realistic) collision detection and resolution for video games. resolv is created in the Go language, but the core concepts are very straightforward and could be easily adapted for use with other languages or game engines and frameworks.

Basically: It allows you to do simple physics easier, without it doing the physics part - that's still on you.

## Why is it called that?

Because it's like... You know, collision resolution? To resolve a collision? So... That's the name. I juste took an e off because I misplaced it somewhere.

## Why did you create Resolv?

Because I was making games and frequently found that most frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for simple, "arcade-y" games. Note that Resolv is generally recommended for use for simpler games with non-grid-based objects. If your game has objects already contained in or otherwise aligned to a 2D array, then it would most likely be more efficient to use that array for collision detection instead of using Resolv.

## How do I install it?

It should be as simple as just go getting it and importing it in your game application, or importing it and ensuring it's within your `go.mod` file.

`go get github.com/solarlune/resolv`

## How do I use it?

There's two ways to use Resolv. One way is to simply create two Shapes, and then check for a collision between them, or attempt to resolve a movement of one into the other, like below:

```go

// Here, we keep a pointer to the shapes we want to check, since we want to create them just once
// in the Init() function, and then check them every frame in the Update() function.

var shape1 *resolv.Rectangle
var shape2 *resolv.Rectangle

func Init() {

    // Create one rectangle, with an X and Y of 10 each, and a Width and Height of 16 each.
    shape1 = resolv.NewRectangle(10, 10, 16, 16)

    // Create another rectangle, as well.
    shape2 = resolv.NewRectangle(11, 100, 16, 16)

}

func Update() {

    // Later on, in the game's update loop...

    // Let's say we were trying to move the shape to the right by 2 pixels. We'll create a delta X movement variable
    // that stores a value of 2.
    dx := 2

    // Here, we check to see if there's a collision should shape1 try to move to the right by the delta movement. The Resolve()
    // functions return a Collision object that has information about whether the attempted movement resulted in a collision 
    // or not.
    resolution := resolv.Resolve(shape1, shape2, dx, 0)

    if resolution.Colliding() {
        
        // If there was a collision, then shape1 couldn't move fully to the right. It came into contact with shape2,
        // and the variable "resolution" now holds a Collision struct with helpful information, like how far to move to be touching.
        
        // Here we just move the shape over to the right by the distance reported by the Collision struct so it'll come into contact 
        // with shape2.
        shape1.X += resolution.ResolveX

    } else {

        // If there wasn't a collision, shape1 should be able to move fully to the right, so we move it.
        shape1.X += dx

    }

    // We can also do collision testing only pretty simply:

    colliding := shape1.IsColliding(shape2)

    if colliding {
        fmt.Println("WHOA! shape1 and shape2 are colliding.")
    }

}

// That's it!

```

This is fine for simple testing, but if you have even a slightly more complex game with a lot more Shapes, then you would have to check each Shape against each other Shape. This is a bit awkward for the developer to code, so I also added Spaces. 

A Space represents a container for Shapes to exist in and test against. The fundamentals are the same, but it scales up more easily, since you don't have to do manual checking everywhere you want to test a Shape against any others. 

A Space is just a pointer to a slice of Shapes, so feel free to use as many as you need to. You could split up a level into multiple Spaces, or have everything in one Space if it works for your game. Spaces also contain functions to filter them out as necessary to easily test a smaller selection of Shapes when desired.

Here's an example using a Space to check one Shape against others:

```go

var space *resolv.Space
var playerRect *resolv.Rectangle

// Here, in the game's init loop...

func Init() {

    // Create a space for Shapes to occupy.
    space = resolv.NewSpace()

    // Create one rectangle - we'll say this one represents our player.
    playerRect = resolv.NewRectangle(40, 40, 16, 16)

    /* Note that we don't HAVE to add the Player Rectangle to the Space; this is only if we
    want it to also be checked for collision testing and resolution within the Space by 
    other Shapes.*/
    space.Add(playerRect)

    // Now we're going to create some more Rectangles to represent level bounds.
    
    // We can also add multiple shapes in a single Add() call...

    space.Add(resolv.NewRectangle(16, 0, 320, 16),
        resolv.NewRectangle(16, 16, 320, 16),
        resolv.NewRectangle(16, 240-16, 320, 16),
        resolv.NewRectangle(320-16, 16, 16, 240-16)
    )

    // Note that this is a bit verbose - in reality, you'd probably be loading the necessary data 
    // to construct the Shapes by looping through a for-loop when reading data in from a 
    // level format, like Tiled's TMX format. Anyway...

    // A Space also has the ability to easily add tags to its Shapes.
    space.AddTags("solid")
    
}

func Update() {

    // This time, we want to see if we're going to collide with something solid when 
    // moving down-right by 4 pixels on each axis.

    dx := 4
    dy := 4

    // To check for Shapes with a specific tag, we can filter out the Space they exist 
    // in with either the Space.FilterByTags() or Space.Filter() functions. Space.Filter() 
    // allows us to provide a function to filter out the Shapes; Space.FilterByTags() 
    // takes tags themselves to filter out the Shapes by.

    // This gives us just the Shapes with the "solid" tag.
    solids := space.FilterByTags("solid")

    // You can provide multiple tags in the same function to filter by all of them at the same
    // time, as well. ( i.e. deathZones := space.FilterByTags("danger", "zone") )

    // Now we check each axis individually against the Space (or, in other words, against
    // all Shapes conatined within the Space). This is done to allow a collision on one 
    // axis to not stop movement on the other as necessary. Note that Space.Resolve() 
    // takes the checking Shape as the first argument, and returns the first collision 
    // that it comes into contact with.

    collision := solids.Resolve(playerRect, dx, 0)

    if collision.Colliding() {
        playerRect.X += collision.ResolveX
    } else {
        playerRect.X += dx
    }

    collision = solids.Resolve(playerRect, 0, dy)

    if collision.Colliding() {
        playerRect.Y += collision.ResolveY
    } else {
        playerRect.Y += dy
    }

}

// Done-zo!

```

Shapes also have functions to get and set a pointer to a data object (of type `interface{}` ), which can be used to point to a user-created object, like a GameObject or Player struct.

```go
type Player struct {
    Rect *resolv.Rectangle
}

func Init() {
    player := &Player{ Rect: resolv.newRectangle(0, 0, 16, 16) }
    player.Rect.SetData(player) 
    // Now if we ever need to get a reference to the Player from its Rect, we 
    // can do so by using the Rectangle's GetData() function. 
}
```

Also of interest, a Space itself satisfies the requirements for a Shape, so they can be checked against like any other Shape. This works like a complex Shape composed of smaller Shapes, where doing collision testing and resolution simply does the equivalent functions for each Shape contained within the Space. This means that you can make complex Shapes out of simple Shapes easily by adding them to a Space, and then using that Space wherever you would use a normal Shape.

```go

var ship *resolv.Space
var world *resolv.Space

func Init() {

    world = resolv.NewSpace()

    ship = resolv.NewSpace()
    
    // Construct the ship!
    ship.Add(
        resolv.NewRectangle(0, 0, 16, 16), 
        resolv.NewLine(16, 0, 32, 16), 
        resolv.NewLine(32, 16, 16, 16))

    // Add the Ship to the game world!
    world.Add(ship)

    // Make something to dodge!
    bullet := resolv.NewRectangle(64, 8, 2, 2)
    bullet.SetTags("bullet")
    world.Add(bullet)

}

func Update() {

    // To make using Spaces as compound Shapes easier, you can use the Space's Move() 
    // function to move all Shapes contained within the Space by the specified delta
    // X and Y values.

    ship.Move(2, 0)

    bullets := world.FilterByTags("bullet")

    // Now this line will run if any bullet touches any part of our ship Space.
    if bullets.IsColliding(ship) {
        fmt.Println("OW! I ran into a bullet!")
    }

    // There are additional functions present in the Space struct to assist with finding
    // colliding objects, as well:
    for _, bullet := range bullets.getCollidingShapes(ship) {
        fmt.Println(bullet)
    }

}

```

Welp, that's about it. If you want to see more info, feel free to examine the examples folder.

[You can check out the documentation here, as well.](https://pkg.go.dev/github.com/solarlune/resolv)

## Dependencies?

For using resolv with your projects, there are no external dependencies. resolv just uses the built-in "fmt" and "math" packages.

For the resolv tests, resolv requires raylib-go to handle input and draw the screen.

## Known Issues

Line-line collision works, but doesn't resolve movement properly if done on edges with non-perpendicular lines.

Circle-line collisions aren't implemented yet.

## Shout-out Time!

Massive thanks to the developer of raylib; it's an awesome framework, and simplifies things CONSIDERABLY. 

Thanks to gen2brain for maintaining raylib-go, as well!

Props to whoever made arcadepi.ttf! It's a nice font.

Thanks a lot to the SDL2 team for development, and thanks to veandco for maintaining the Golang SDL2 port, as well!

Thanks to the people who stopped by on my stream - they helped out a lot with a couple of the technical aspects of getting Go to do what I needed to, haha.
