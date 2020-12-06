
# Resolv

GIF GOES HERE

[Documentation](https://pkg.go.dev/github.com/SolarLune/resolv/resolv?tab=doc) (It's also, of course, documented directly in the code.)

## What is Resolv?

Resolv is a collision detection library, specifically created for simple, arcade (non-realistic) collision detection and resolution for video games. Resolv is created in the Go language, but the core concepts are very straightforward and could be easily adapted for use with other languages or game engines and frameworks.

Basically: It allows you to do simple physics easier, without it doing the physics part - that's still on you.

## Why is it called that?

Because it's like... You know, collision resolution? To resolve a collision? So... That's the name. I juste took an e off because I misplaced it somewhere.

## Why did you create resolv?

Because I was making games and frequently found that most frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for simple, "arcade-y" games. Note that Resolv is generally recommended for use for simpler games with non-grid-based objects. If your game has objects already contained in or otherwise aligned to a 2D array, then it would most likely be more efficient to use that array for collision detection instead of using Resolv.

## How do I install it?

It should be as simple as just go getting it and using it, or simply importing it in your game application if you're using Go modules.

`go get github.com/SolarLune/resolv`

## How do I use it?

There's one relatively simple way to use resolv - create a Space, create Objects, update the Objects as they move, and check for collisions.

```go

var space *resolv.Space
var playerBody *resolv.Object

func Init() {

    // First, we want to create a Space. This represents the areas in our game world 
    // where objects can move about and check for collisions.

    // The first two arguments represent our grid width and height, while the following two
    // represent our cells' sizes. This grid represents how fine the collision detection is.
    space = resolv.NewSpace(40, 20, 16, 16)

    // Next, we can start creating things.

    // Here's some level geometry; we don't need to actually store it anywhere unless we plan on 
    // moving it around or reinstantiating it at any point AFTER removing it from the space.
    
    // NewObject takes the X, Y, width, and height of the object, and the space to put it in.
    resolv.NewObject(0, 0, 1280, 16, space)
    resolv.NewObject(0, 720-16, 1280, 16, space)
    resolv.NewObject(0, 16, 16, 720-32, space)
    resolv.NewObject(1280-16, 16, 16, 720-32, space)

    // We'll keep a reference to the player's body to move it later.
    playerBody = resolv.NewObject(32, 32, 16, 16, space)

}

func Update() {

    // Later on, in the game's update loop, which runs once per game frame...

    // Let's say we are attempting to move the player to the right. Here's how we could do it.

    dx := 2.0

    // Here, we check to see if there's a collision if playerBody were to move to the right by 2 pixels. The function returns
    // a Collision object.

    collision := playerBody.Check(dx, 0)

    if collision.Valid() {
        
        // If there was a collision, then it was a valid collision, and the playerBody can't move fully to the right by 2. It came into contact 
        // with an occupied Cell to the right, and the 'collision' variable now holds a Collision struct with helpful information, like how 
        // far to move to be touching the cell to the right.
        
        // Here we just move the shape over to the right by the distance reported by the Collision struct so it'll come into contact 
        // with the colliding Cell.
        playerBody.X += collision.ContactX

    } else {
        // If there wasn't a collision, playerBody can move the full distance, so we simply apply the motion to the X coordinate.
        playerBody.X += dx
    }

    // If this was for a platformer, you could then check this for the Y-axis as well.

    // If you want to filter out types of Objects to check for, add tags on the objects you want to filter using Object.AddTags() 
    // and specify them in the last arguments of Object.Check.

    onlySolidHazardous := playerBody.Check(dx, 0, "hazard", "solid")

}

// That's it!

```

You can do much more than this with the provided API, but this is the basic concept.

## Wasn't this... different before?

That is correct; Resolv did have more to it, initially. It had multiple object types and geometry tests, like line-line intersection. However, I was generally unhappy with it for a couple of reasons:

- Object positions were represented using integer numbers. This was pretty bad, as this made it difficult to make objects that moved gradually, using speeds stored in floats.
- Relatively poor performance. It was fine for small games, but collision testing should be _extremely_ lightweight, especially considering how frequently objects in games need to collide. The previous iteration basically boiled down to _checking every shape that moves against every other moving shape_, which is, well, bad.
- No concept of "physical space". Originally, each object existed independently of each other, and was not tied down to any spatial grid or partitioning system. This meant that the only real way to check a space for an object was to loop through every shape, again, to see if it collides with that point or area.
- Relatively difficult to change and improve it, as my knowledge of high school-level geometry is lacking at best. (COUGH COUGH HACK SNEEZE)

So, I decided to overhaul it; after several attempts and three separate attempts at revamps, I came up with this as a simple, but effective, improvement. It still does what the original version of Resolv did well; namely, being very simple and easy to use, and being made primarily for arcade game physics. It still stays out of your way, but now it's more performant and easier to use for actual use-cases. The codebase is also much smaller and simpler, which is nice.

----

Welp, that's about it. If you want to see more info, feel free to examine the main.go and world#.go tests to see how a couple of quick example tests are set up.

[You can check out the documentation here, as well.](https://pkg.go.dev/github.com/SolarLune/resolv/resolv?tab=doc)

## Dependencies?

For using resolv with your projects, there are no external dependencies. resolv just uses the built-in "fmt" and "math" packages.

For the resolv tests, resolv requires [ebiten](github.com/hajimehoshi/ebiten) and [vector](github.com/kvartborg/vector). Both of these are modules, so you should be able to simply run `go run ./examples` from the base directory for Go to download them (and Resolv) to run tests successfully.

## Shout-out Time!

Thanks to the people who stopped by on my stream - they helped out a lot with a couple of the technical aspects of getting Go to do what I needed to, haha.
