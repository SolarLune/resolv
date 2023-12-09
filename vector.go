package resolv

import (
	"fmt"
	"math"
)

// This is essentially a 2D version of my 3D Vectors used in Tetra3D.

// WorldRight represents a unit vector in the global direction of WorldRight on the right-handed OpenGL / Tetra3D's coordinate system (+X).
var WorldRight = NewVector(1, 0)

// WorldLeft represents a unit vector in the global direction of WorldLeft on the right-handed OpenGL / Tetra3D's coordinate system (-X).
var WorldLeft = WorldRight.Invert()

// WorldUp represents a unit vector in the global direction of WorldUp on the right-handed OpenGL / Tetra3D's coordinate system (+Y).
var WorldUp = NewVector(0, 1)

// WorldDown represents a unit vector in the global direction of WorldDown on the right-handed OpenGL / Tetra3D's coordinate system (+Y).
var WorldDown = WorldUp.Invert()

// Vector represents a 2D Vector, which can be used for usual 2D applications (position, direction, velocity, etc).
// Any Vector functions that modify the calling Vector return copies of the modified Vector, meaning you can do method-chaining easily.
// Vectors seem to be most efficient when copied (so try not to store pointers to them if possible, as dereferencing pointers
// can be more inefficient than directly acting on data, and storing pointers moves variables to heap).
type Vector struct {
	X float64 // The X (1st) component of the Vector
	Y float64 // The Y (2nd) component of the Vector
}

// NewVector creates a new Vector with the specified x, y, and z components. The W component is generally ignored for most purposes.
func NewVector(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

// NewVectorZero creates a new "zero-ed out" Vector, with the values of 0, 0, 0, and 0 (for W).
func NewVectorZero() Vector {
	return Vector{}
}

// Modify returns a ModVector object (a pointer to the original vector).
func (vec *Vector) Modify() ModVector {
	ip := ModVector{Vector: vec}
	return ip
}

// String returns a string representation of the Vector, excluding its W component (which is primarily used for internal purposes).
func (vec Vector) String() string {
	return fmt.Sprintf("{%.2f, %.2f}", vec.X, vec.Y)
}

// Plus returns a copy of the calling vector, added together with the other Vector provided (ignoring the W component).
func (vec Vector) Add(other Vector) Vector {
	vec.X += other.X
	vec.Y += other.Y
	return vec
}

// Sub returns a copy of the calling Vector, with the other Vector subtracted from it (ignoring the W component).
func (vec Vector) Sub(other Vector) Vector {
	vec.X -= other.X
	vec.Y -= other.Y
	return vec
}

// Expand expands the Vector by the margin specified, in absolute units, if each component is over the minimum argument.
// To illustrate: Given a Vector of {1, 0.1, -0.3}, Vector.Expand(0.5, 0.2) would give you a Vector of {1.5, 0.1, -0.8}.
// This function returns a copy of the Vector with the result.
func (vec Vector) Expand(margin, min float64) Vector {
	if vec.X > min || vec.X < -min {
		vec.X += math.Copysign(margin, vec.X)
	}
	if vec.Y > min || vec.Y < -min {
		vec.Y += math.Copysign(margin, vec.Y)
	}
	return vec
}

// Invert returns a copy of the Vector with all components inverted.
func (vec Vector) Invert() Vector {
	vec.X = -vec.X
	vec.Y = -vec.Y
	return vec
}

// Magnitude returns the length of the Vector.
func (vec Vector) Magnitude() float64 {
	return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)
}

// MagnitudeSquared returns the squared length of the Vector; this is faster than Length() as it avoids using math.Sqrt().
func (vec Vector) MagnitudeSquared() float64 {
	return vec.X*vec.X + vec.Y*vec.Y
}

// ClampMagnitude clamps the overall magnitude of the Vector to the maximum magnitude specified, returning a copy with the result.
func (vec Vector) ClampMagnitude(maxMag float64) Vector {
	if vec.Magnitude() > maxMag {
		vec = vec.Unit().Scale(maxMag)
	}
	return vec
}

// SubMagnitude subtracts the given magnitude from the Vector's existing magnitude.
// If the vector's magnitude is less than the given magnitude to subtract, a zero-length Vector will be returned.
func (vec Vector) SubMagnitude(mag float64) Vector {
	if vec.Magnitude() > mag {
		return vec.Sub(vec.Unit().Scale(mag))
	}
	return Vector{0, 0}

}

// Distance returns the distance from the calling Vector to the other Vector provided.
func (vec Vector) Distance(other Vector) float64 {
	return vec.Sub(other).Magnitude()
}

// Distance returns the squared distance from the calling Vector to the other Vector provided. This is faster than Distance(), as it avoids using math.Sqrt().
func (vec Vector) DistanceSquared(other Vector) float64 {
	return vec.Sub(other).MagnitudeSquared()
}

// Mult performs Hadamard (component-wise) multiplication on the calling Vector with the other Vector provided, returning a copy with the result (and ignoring the Vector's W component).
func (vec Vector) Mult(other Vector) Vector {
	vec.X *= other.X
	vec.Y *= other.Y
	return vec
}

// Unit returns a copy of the Vector, normalized (set to be of unit length).
// It does not alter the W component of the Vector.
func (vec Vector) Unit() Vector {
	l := vec.Magnitude()
	if l < 1e-8 || l == 1 {
		// If it's 0, then don't modify the vector
		return vec
	}
	vec.X, vec.Y = vec.X/l, vec.Y/l
	return vec
}

// SetX sets the X component in the vector to the value provided.
func (vec Vector) SetX(x float64) Vector {
	vec.X = x
	return vec
}

// SetY sets the Y component in the vector to the value provided.
func (vec Vector) SetY(y float64) Vector {
	vec.Y = y
	return vec
}

// Set sets the values in the Vector to the x, y, and z values provided.
func (vec Vector) Set(x, y float64) Vector {
	vec.X = x
	vec.Y = y
	return vec
}

// Floats returns a [2]float64 array consisting of the Vector's contents.
func (vec Vector) Floats() [2]float64 {
	return [2]float64{vec.X, vec.Y}
}

// Equals returns true if the two Vectors are close enough in all values (excluding W).
func (vec Vector) Equals(other Vector) bool {

	eps := 1e-4

	if math.Abs(float64(vec.X-other.X)) > eps || math.Abs(float64(vec.Y-other.Y)) > eps {
		return false
	}

	return true

}

// IsZero returns true if the values in the Vector are extremely close to 0 (excluding W).
func (vec Vector) IsZero() bool {

	eps := 1e-4

	if math.Abs(float64(vec.X)) > eps || math.Abs(float64(vec.Y)) > eps {
		return false
	}

	// if !onlyXYZ && math.Abs(vec.W-other.W) > eps {
	// 	return false
	// }

	return true

}

// Rotate returns a copy of the Vector, rotated around an axis Vector with the x, y, and z components provided, by the angle
// provided (in radians), counter-clockwise.
// The function is most efficient if passed an orthogonal, normalized axis (i.e. the X, Y, or Z constants).
// Note that this function ignores the W component of both Vectors.
func (vec Vector) Rotate(angle float64) Vector {
	x := vec.X
	y := vec.Y
	vec.X = x*math.Cos(angle) - y*math.Sin(angle)
	vec.Y = x*math.Sin(angle) + y*math.Cos(angle)
	return vec
}

// Angle returns the angle between the calling Vector and the provided other Vector (ignoring the W component).
func (vec Vector) Angle(other Vector) float64 {
	d := vec.Unit().Dot(other.Unit())
	d = clamp(d, -1, 1) // Acos returns NaN if value < -1 or > 1
	return math.Acos(float64(d))
}

func (vec Vector) AngleRotation() float64 {
	return vec.Angle(WorldRight)
}

// Scale scales a Vector by the given scalar (ignoring the W component), returning a copy with the result.
func (vec Vector) Scale(scalar float64) Vector {
	vec.X *= scalar
	vec.Y *= scalar
	return vec
}

// Divide divides a Vector by the given scalar (ignoring the W component), returning a copy with the result.
func (vec Vector) Divide(scalar float64) Vector {
	vec.X /= scalar
	vec.Y /= scalar
	return vec
}

// Dot returns the dot product of a Vector and another Vector (ignoring the W component).
func (vec Vector) Dot(other Vector) float64 {
	return vec.X*other.X + vec.Y*other.Y
}

// Round rounds off the Vector's components to the given space in world unit increments, returning a clone
// (e.g. Vector{0.1, 1.27, 3.33}.Snap(0.25) will return Vector{0, 1.25, 3.25}).
func (vec Vector) Round(snapToUnits float64) Vector {
	vec.X = round(vec.X/snapToUnits) * snapToUnits
	vec.Y = round(vec.Y/snapToUnits) * snapToUnits
	return vec
}

// ClampAngle clamps the Vector such that it doesn't exceed the angle specified (in radians).
// This function returns a normalized (unit) Vector.
func (vec Vector) ClampAngle(baselineVec Vector, maxAngle float64) Vector {

	mag := vec.Magnitude()

	angle := vec.Angle(baselineVec)

	if angle > maxAngle {
		vec = baselineVec.Slerp(vec, maxAngle/angle).Unit()
	}

	return vec.Scale(mag)

}

// Lerp performs a linear interpolation between the starting Vector and the provided
// other Vector, to the given percentage (ranging from 0 to 1).
func (vec Vector) Lerp(other Vector, percentage float64) Vector {
	percentage = clamp(percentage, 0, 1)
	vec.X = vec.X + ((other.X - vec.X) * percentage)
	vec.Y = vec.Y + ((other.Y - vec.Y) * percentage)
	return vec
}

// Slerp performs a spherical linear interpolation between the starting Vector and the provided
// ending Vector, to the given percentage (ranging from 0 to 1).
// This should be done with directions, usually, rather than positions.
// This being the case, this normalizes both Vectors.
func (vec Vector) Slerp(targetDirection Vector, percentage float64) Vector {

	vec = vec.Unit()
	targetDirection = targetDirection.Unit()

	// Thank you StackOverflow, once again! : https://stackoverflow.com/questions/67919193/how-does-unity-implements-vector3-slerp-exactly
	percentage = clamp(percentage, 0, 1)

	dot := vec.Dot(targetDirection)

	dot = clamp(dot, -1, 1)

	theta := math.Acos(dot) * percentage
	relative := targetDirection.Sub(vec.Scale(dot)).Unit()

	return (vec.Scale(math.Cos(theta)).Add(relative.Scale(math.Sin(theta)))).Unit()

}

// ModVector represents a reference to a Vector, made to facilitate easy method-chaining and modifications on that Vector (as you
// don't need to re-assign the results of a chain of operations to the original variable to "save" the results).
// Note that a ModVector is not meant to be used to chain methods on a vector to pass directly into a function; you can just
// use the normal vector functions for that purpose. ModVectors are pointers, which are allocated to the heap. This being the case,
// they should be slower relative to normal Vectors, so use them only in non-performance-critical parts of your application.
type ModVector struct {
	*Vector
}

// Add adds the other Vector provided to the ModVector.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Add(other Vector) ModVector {
	ip.X += other.X
	ip.Y += other.Y
	return ip
}

// Sub subtracts the other Vector from the calling ModVector.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Sub(other Vector) ModVector {
	ip.X -= other.X
	ip.Y -= other.Y
	return ip
}

func (ip ModVector) SetZero() ModVector {
	ip.X = 0
	ip.Y = 0
	return ip
}

// Scale scales the Vector by the scalar provided.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Scale(scalar float64) ModVector {
	ip.X *= scalar
	ip.Y *= scalar
	return ip
}

// Divide divides a Vector by the given scalar (ignoring the W component).
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Divide(scalar float64) ModVector {
	ip.X /= scalar
	ip.Y /= scalar
	return ip
}

// Expand expands the ModVector by the margin specified, in absolute units, if each component is over the minimum argument.
// To illustrate: Given a ModVector of {1, 0.1, -0.3}, ModVector.Expand(0.5, 0.2) would give you a ModVector of {1.5, 0.1, -0.8}.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Expand(margin, min float64) ModVector {
	exp := ip.Vector.Expand(margin, min)
	ip.X = exp.X
	ip.Y = exp.Y
	return ip
}

// Mult performs Hadamard (component-wise) multiplication with the Vector on the other Vector provided.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Mult(other Vector) ModVector {
	ip.X *= other.X
	ip.Y *= other.Y
	return ip
}

// Unit normalizes the ModVector (sets it to be of unit length).
// It does not alter the W component of the Vector.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Unit() ModVector {
	l := ip.Magnitude()
	if l < 1e-8 || l == 1 {
		// If it's 0, then don't modify the vector
		return ip
	}
	ip.X, ip.Y = ip.X/l, ip.Y/l
	return ip
}

// Rotate rotates the calling Vector by the angle provided (in radians).
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Rotate(angle float64) ModVector {
	rot := (*ip.Vector).Rotate(angle)
	ip.X = rot.X
	ip.Y = rot.Y
	return ip
}

// Invert inverts all components of the calling Vector.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Invert() ModVector {
	ip.X = -ip.X
	ip.Y = -ip.Y
	return ip
}

// Round snaps the Vector's components to the given space in world units, returning a clone (e.g. Vector{0.1, 1.27, 3.33}.Snap(0.25) will return Vector{0, 1.25, 3.25}).
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Round(snapToUnits float64) ModVector {
	snapped := (*ip.Vector).Round(snapToUnits)
	ip.X = snapped.X
	ip.Y = snapped.Y
	return ip
}

// ClampMagnitude clamps the overall magnitude of the Vector to the maximum magnitude specified.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) ClampMagnitude(maxMag float64) ModVector {
	clamped := (*ip.Vector).ClampMagnitude(maxMag)
	ip.X = clamped.X
	ip.Y = clamped.Y
	return ip
}

// SubMagnitude subtacts the given magnitude from the Vector's. If the vector's magnitude is less than the given magnitude to subtract,
// a zero-length Vector will be returned.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) SubMagnitude(mag float64) ModVector {
	if ip.Magnitude() > mag {
		ip.Sub(ip.Vector.Unit().Scale(mag))
	} else {
		ip.X = 0
		ip.Y = 0
	}
	return ip

}

// Lerp performs a linear interpolation between the starting Vector and the provided
// other Vector, to the given percentage (ranging from 0 to 1).
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Lerp(other Vector, percentage float64) ModVector {
	lerped := (*ip.Vector).Lerp(other, percentage)
	ip.X = lerped.X
	ip.Y = lerped.Y
	return ip
}

// Slerp performs a linear interpolation between the starting Vector and the provided
// other Vector, to the given percentage (ranging from 0 to 1).
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Slerp(targetDirection Vector, percentage float64) ModVector {
	slerped := (*ip.Vector).Slerp(targetDirection, percentage)
	ip.X = slerped.X
	ip.Y = slerped.Y
	return ip
}

// ClampAngle clamps the Vector such that it doesn't exceed the angle specified (in radians).
// This function returns the calling ModVector for method chaining.
func (ip ModVector) ClampAngle(baselineVector Vector, maxAngle float64) ModVector {
	clamped := (*ip.Vector).ClampAngle(baselineVector, maxAngle)
	ip.X = clamped.X
	ip.Y = clamped.Y
	return ip
}

// String converts the ModVector to a string. Because it's a ModVector, it's represented with a *.
func (ip ModVector) String() string {
	return fmt.Sprintf("*{%.2f, %.2f}", ip.X, ip.Y)
}

// Clone returns a ModVector of a clone of its backing Vector.
// This function returns the calling ModVector for method chaining.
func (ip ModVector) Clone() ModVector {
	v := *ip.Vector
	return v.Modify()
}

func (ip ModVector) ToVector() Vector {
	return *ip.Vector
}
