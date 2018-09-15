
# resolv

![peek 2018-09-15 02-26](https://user-images.githubusercontent.com/4733521/45585063-0c3fa100-b893-11e8-93df-7a14be9992ae.gif)

## What is resolv?

resolv is a library specifically created for simple arcade (non-realistic) collision detection and resolution for video games. resolv is created in the Go language, but the core concepts are very straightforward and could be easily adapted for use with other languages or game engines and frameworks.

Basically: It allows you to do simple physics easier, without it doing the physics part - that's still on you.

## Why is it called that?

Because it's like... You know, collision resolution? To resolve a collision? So... That's the name. I juste took an e off because I misplaced it somewhere.

## Why did you create resolv?

Because I was making games and frequently found that most frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most 2D games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for arcade applications.

## How do you install it?

It should be as simple as just go getting it and importing it in your game application.

`go get github.com/SolarLune/resolv`

## How do you use it?

There's two ways to use resolv. One way is to simply create two Shapes, and then attempt to resolve a movement of one into the other, like below. (Note that this is untested pseudo-code, but the idea is fine.)

```go

// Here, we keep a pointer to the shapes we want to check, since we want to create them just once
// in the Init() function, and then check them every frame in the Update() function.

var shape1 *resolv.Rectangle
var shape2 *resolv.Rectangle

func Init() {

    // Create one rectangle, with an X and Y of 10 each, and a Width and Height of 16 each.
    shape1 = resolv.NewRectangle(10, 10, 16, 16)

    // Create another rectangle, as well.
    shape2 = resolv.NewRectangle(11, 100, 0, 0)

}

func Update() {

    // Later on, in the game's update loop...

    // Let's say we were trying to move the shape to the right by 2 pixels. We'll create a delta X movement variable
    // that stores a value of 2.
    dx := 2

    // Here, we check to see if there's a collision should shape1 try to move to the right by 10 pixels. The Resolve()
    // functions return a Collision object that has information about whether the attempted movement would work,
    // and whether it resulted in a collision or not.
    collision := shape1.Resolve(shape2, dx, true)

    if collision.Colliding() {
        
        // If there was a collision, then shape1 couldn't move fully to the right. It came into contact with shape2,
        // and the variable collision now holds a Collision object with helpful information, like how far it was able to move.
        // Move the shape over to the right by the distance that it can to come into full contact with shape2.
        shape1.X += collision.ResolveDistance

    } else {

        // If there wasn't a collision, shape1 should be able to move fully to the right, so we move it.
        shape1.X += dx

    }

}

// That's it!

```

This is fine for simple testing, but if you have a more complex game with a lot more Shapes, then you would have to check each Shape against each other Shape. This is awkward for the developer to code, so I also added Spaces. 

A Space represents a container for Shapes to exist in and test against. This way, the fundamentals are the same, but it should scale up more easily, since you don't have to do manual for checking everywhere you want to test a Shape against others. A Space is just a collection of Shapes, so feel free to use as many as you need to (i.e. you could split up a level into multiple Spaces, or have everything in one Space if it works for your game). Here's an example using a Space to check one shape against others more easily:

```go

// Here, in the game's init loop...

var space resolv.Space
var playerRect *resolv.Rectangle

func Init() {

    // Create a space.
    space = resolv.NewSpace()

    // Create one rectangle - we'll say this one represents our player.
    playerRect = resolv.NewRectangle(40, 40, 16, 16)

    // Add the shape to the space
    space.AddShape(playerRect)

    // Create some more Rectangles to represent level bounds.
    s := resolv.NewRectangle(0, 0, 16, 240)

    // We set tags on the Rectangles to allow us to more easily check for collisions by specific "type".
    s.SetTag("solid")
    
    // Then we add the shape to the space.
    space.AddShape(s)

    // Note that this is a bit verbose - in reality, you'd probably be loading the necessary data to construct the Shapes
    // by looping through a for loop, reading data in from a level format, like Tiled's TMX format. Then you'd just do it once in
    // the for loop to have it be done for each Shape you need to represent your level geometry, rather than hand-coding the shapes
    // like this. Anyway...

    s = resolv.NewRectangle(16, 0, 320, 16)
    s.SetTag("solid")
    space.AddShape(s)

    s = resolv.NewRectangle(16, 240-16, 320, 16)
    s.SetTag("solid")
    space.AddShape(s)
    
    s = resolv.NewRectangle(320-16, 16, 16, 240-16)
    s.SetTag("solid")
    space.AddShape(s)

}

func Update() {

    // This time, we want to see if we're going to collide with something moving down-right by 2 pixels, each axis.
    dx := 2
    dy := 2

    // Now we check each axis individually. This is done to allow a collision on one axis to not stop movement on the other
    // as necessary. The "solid" tag goes here, so we only resolve a collision against Shapes that have that tag.
    collision := space.Resolve(playerRect, dx, true, "solid")

    if collision.Colliding() {
        playerRect.X += collision.ResolveDistance
    } else {
        playerRect.X += dx
    }

    collision = space.Resolve(playerRect, dy, false, "solid")

    if collision.Colliding() {
        playerRect.Y += collision.ResolveDistance
    } else {
        playerRect.Y += dy
    }

}

// Done-zo!

```

Welp, that's about it. If you want to see more info, feel free to examine the main.go and world.go tests to see how a couple of quick example tests are set up.

## Dependencies?

For the actual package, there are no external dependencies. resolv just uses the built-in "fmt" and "math" packages. It also exists in just one Go file (currently, anyway), which is also a good thing, I suppose, haha.

For the tests, resolv requires veandco's sdl2 port to create the window, handle input, and draw the shapes.

## Shout-out Time!

Props to whoever made arcadepi.ttf!

Thanks a lot to the SDL2 team for development.

Thanks to veandco for maintaining the Golang SDL2 port, as well!
