package main

import "fmt"

type SafetyPlacer interface {
	placeSafeties()
}

type RockClimber struct {
	rocksClimed int
	sp          SafetyPlacer
}

type IceSafetyPlacer struct{}

func (rc *RockClimber) climbRock() {
	rc.rocksClimed++
	if rc.rocksClimed == 10 {
		rc.sp.placeSafeties()
	}
}

func newRockClimber(sp SafetyPlacer) *RockClimber {
	return &RockClimber{
		sp: sp,
	}
}

func (sp *IceSafetyPlacer) placeSafeties() {
	fmt.Println("PLacing my ICE safeties...")
}

type NOPSafetyPlacer struct{}

func (sp *NOPSafetyPlacer) placeSafeties() {
	fmt.Println("PLacing my NOOP safeties...")
}

func main() {
	rc := newRockClimber(&IceSafetyPlacer{})

	for i := 0; i < 11; i++ {
		rc.climbRock()
	}
}
