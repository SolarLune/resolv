/*
Package resolv is a simple collision detection and resolution library mainly geared towards simpler 2D arcade-style games. Its
goal is to be lightweight, fast, simple, and easy-to-use for game development. Its goal is to also not become a physics engine
or physics library itself, but to always leave the actual physics implementation and "game feel" to the developer, while making
it very easy to do so.

Usage of resolv essentially centers around two main concepts: Spaces and Shapes.

A Shape can be used to test for collisions against another Shape. That's really all they have to do, but that capability is powerful
when paired with the resolv.Resolve() function. You can then check to see if a Shape would have a collision if it attempted to move
in a specified direction. If so, the Resolve() function would return a Collision object, which tells you some information about the
Collision, like how far the checking Shape would have to move to come into contact with the other, and which Shape it comes into
contact with.

A Space is just a slice that holds Shapes for detection. It doesn't represent any real physical space, and so there aren't any
units of measurement to remember when using Spaces. Similar to Shapes, Spaces are simple, but also very powerful. Spaces allow
you to easily check for collision with, and resolve collision against multiple Shapes within that Space. A Space being just a
collection of Shapes means that you can manipulate and filter them as necessary.
*/
package resolv
