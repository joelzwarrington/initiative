package dnd

import "testing"

func TestNewMonster(t *testing.T) {
	statBlock := StatBlock{
		HitPoints: HitPoints{Fixed: 25},
		Challenge: Challenge{Rating: 1, ExperiencePoints: 200},
	}

	monster := NewMonster("Goblin", statBlock)

	if monster.Name() != "Goblin" {
		t.Errorf("NewMonster().Name() = %v, want %v", monster.Name(), "Goblin")
	}
	if monster.HitPoints() != 25 {
		t.Errorf("NewMonster().HitPoints() = %v, want %v", monster.HitPoints(), 25)
	}
	if monster.MaximumHitPoints() != 25 {
		t.Errorf("NewMonster().MaximumHitPoints() = %v, want %v", monster.MaximumHitPoints(), 25)
	}
}

func TestMonster_Health(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	monster := NewMonster("Orc", statBlock)

	if got := monster.HitPoints(); got != 15 {
		t.Errorf("Monster.HitPoints() = %v, want %v", got, 15)
	}
	if got := monster.MaximumHitPoints(); got != 15 {
		t.Errorf("Monster.MaximumHitPoints() = %v, want %v", got, 15)
	}
}

func TestMonster_SetHitPoints(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	baseMonster := NewMonster("Orc", statBlock)

	tests := []struct {
		name     string
		monster  *Monster
		newHP    int
		wantHP   int
		wantName string
	}{
		{
			name:     "set normal hit points",
			monster:  baseMonster,
			newHP:    10,
			wantHP:   10,
			wantName: "Orc",
		},
		{
			name:     "set hit points to zero",
			monster:  baseMonster,
			newHP:    0,
			wantHP:   0,
			wantName: "Orc",
		},
		{
			name:     "set negative hit points (should clamp to 0)",
			monster:  baseMonster,
			newHP:    -3,
			wantHP:   0,
			wantName: "Orc",
		},
		{
			name:     "set hit points above maximum (should clamp)",
			monster:  baseMonster,
			newHP:    20,
			wantHP:   15,
			wantName: "Orc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mon := tt.monster
			mon.SetHitPoints(tt.newHP)

			if got := mon.HitPoints(); got != tt.wantHP {
				t.Errorf("Monster.HitPoints() after SetHitPoints(%v) = %v, want %v", tt.newHP, got, tt.wantHP)
			}
			if got := mon.Name(); got != tt.wantName {
				t.Errorf("Monster.Name() after SetHitPoints() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestMonster_WithName(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	original := NewMonster("Orc", statBlock)

	result := original.WithName("Elite Orc")

	if got := result.Name(); got != "Elite Orc" {
		t.Errorf("Monster.WithName().Name() = %v, want %v", got, "Elite Orc")
	}
	if got := result.HitPoints(); got != 15 {
		t.Errorf("Monster.WithName().HitPoints() = %v, want %v", got, 15)
	}
	if got := result.MaximumHitPoints(); got != 15 {
		t.Errorf("Monster.WithName().MaximumHitPoints() = %v, want %v", got, 15)
	}
}

func TestMonster_AdjustHitPoints(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 20}}
	baseMonster := NewMonster("Orc", statBlock)

	tests := []struct {
		name       string
		monster    *Monster
		adjustment int
		wantHP     int
		wantName   string
	}{
		{
			name:       "healing (positive adjustment)",
			monster:    NewMonster("Orc", statBlock).WithName("Orc"),
			adjustment: 5,
			wantHP:     20, // Clamped to maximum
			wantName:   "Orc",
		},
		{
			name:       "damage (negative adjustment)",
			monster:    baseMonster,
			adjustment: -8,
			wantHP:     12,
			wantName:   "Orc",
		},
		{
			name: "healing beyond maximum (should clamp)",
			monster: func() *Monster {
				m := NewMonster("Orc", statBlock)
				m.SetHitPoints(18) // Set to below max first
				return m
			}(),
			adjustment: 5,
			wantHP:     20,
			wantName:   "Orc",
		},
		{
			name: "damage below zero (should clamp to 0)",
			monster: func() *Monster {
				m := NewMonster("Orc", statBlock)
				m.SetHitPoints(5) // Set to low health first
				return m
			}(),
			adjustment: -10,
			wantHP:     0,
			wantName:   "Orc",
		},
		{
			name: "zero adjustment (no change)",
			monster: func() *Monster {
				m := NewMonster("Orc", statBlock)
				m.SetHitPoints(12) // Set to specific value first
				return m
			}(),
			adjustment: 0,
			wantHP:     12,
			wantName:   "Orc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monster := tt.monster
			actualAdjustment := monster.AdjustHitPoints(tt.adjustment)

			// Test that we get the actual adjustment amount
			_ = actualAdjustment // We could test this if desired

			if got := monster.HitPoints(); got != tt.wantHP {
				t.Errorf("Monster.HitPoints() after AdjustHitPoints(%v) = %v, want %v", tt.adjustment, got, tt.wantHP)
			}
			if got := monster.Name(); got != tt.wantName {
				t.Errorf("Monster.Name() after AdjustHitPoints() = %v, want %v", got, tt.wantName)
			}
		})
	}
}
