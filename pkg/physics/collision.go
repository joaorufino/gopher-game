package physics

// CheckCollisionOnX checks if two rigid bodies are colliding on the X axis.
func CheckCollisionOnX(a, b *RigidBody) bool {
	return a.Position.X < b.Position.X+b.Size.X &&
		a.Position.X+a.Size.X > b.Position.X
}

// CheckCollisionOnY checks if two rigid bodies are colliding on the Y axis.
func CheckCollisionOnY(a, b *RigidBody) bool {
	return a.Position.Y < b.Position.Y+b.Size.Y &&
		a.Position.Y+a.Size.Y > b.Position.Y
}

// CheckIfOnTop checks if one rigid body is on top of another.
func CheckIfOnTop(a, b *RigidBody) bool {
	// Check if 'a' is directly above 'b'
	return a.Position.Y+a.Size.Y <= b.Position.Y && // a's bottom is above b's top
		a.Position.Y+a.Size.Y >= b.Position.Y-1 && // a's bottom is not too far above b's top (tolerance of 1 unit)
		a.Position.X < b.Position.X+b.Size.X && // a's right edge is to the left of b's right edge
		a.Position.X+a.Size.X > b.Position.X // a's left edge is to the right of b's left edge
}

// ResolveCollision resolves collision between two rigid bodies.
func ResolveCollision(a, b *RigidBody) {
	if a.IsStatic && b.IsStatic {
		return
	}

	// Determine which body is movable and which is static/pushable
	var movable, static *RigidBody
	if a.IsStatic {
		movable, static = b, a
	} else if b.IsStatic {
		movable, static = a, b
	} else if a.IsPushable && b.IsPushable {
		if a.Mass > b.Mass {
			movable, static = a, b
		} else {
			movable, static = b, a
		}
	} else if a.IsPushable {
		movable, static = a, b
	} else {
		movable, static = b, a
	}

	// Check if the movable body can push the static body
	if static.IsPushable {
		visited := make(map[*RigidBody]bool)
		if !canPush(movable, static, visited) {
			return
		}
	}

	// Calculate the overlap on each axis
	overlapX := (movable.Size.X/2 + static.Size.X/2) - abs((movable.Position.X+movable.Size.X/2)-(static.Position.X+static.Size.X/2))
	overlapY := (movable.Size.Y/2 + static.Size.Y/2) - abs((movable.Position.Y+movable.Size.Y/2)-(static.Position.Y+static.Size.Y/2))

	// Resolve collision based on the smaller overlap
	if overlapX < overlapY {
		if CheckCollisionOnX(movable, static) {
			if movable.Position.X < static.Position.X {
				movable.Position.X = static.Position.X - movable.Size.X
			} else {
				movable.Position.X = static.Position.X + static.Size.X
			}
			movable.Velocity.X = 0
		}
	} else {
		if CheckCollisionOnY(movable, static) {
			if movable.Position.Y < static.Position.Y {
				movable.Position.Y = static.Position.Y - movable.Size.Y
				if CheckIfOnTop(movable, static) {
					movable.OnGround = true
				}
			} else {
				movable.Position.Y = static.Position.Y + static.Size.Y
			}
			movable.Velocity.Y = 0
		}
	}
}

// canPush checks if a can push b (and recursively checks if a can push b's pushable chain).
func canPush(a *RigidBody, b *RigidBody, visited map[*RigidBody]bool) bool {
	if a.Mass <= b.Mass || visited[b] {
		return false
	}
	visited[b] = true

	for _, rb := range b.CollidingBodies {
		if rb.IsPushable && !canPush(a, rb, visited) {
			return false
		}
	}

	return true
}

// abs is a helper function to calculate the absolute value of a float64.
func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
