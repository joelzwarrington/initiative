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
	if monster.HP() != 25 {
		t.Errorf("NewMonster().HP() = %v, want %v", monster.HP(), 25)
	}
	if monster.MaxHP() != 25 {
		t.Errorf("NewMonster().MaxHP() = %v, want %v", monster.MaxHP(), 25)
	}
}

func TestMonster_Health(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	monster := NewMonster("Orc", statBlock)

	if got := monster.HP(); got != 15 {
		t.Errorf("Monster.HP() = %v, want %v", got, 15)
	}
	if got := monster.MaxHP(); got != 15 {
		t.Errorf("Monster.MaxHP() = %v, want %v", got, 15)
	}
}

func TestMonster_SetHP(t *testing.T) {
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
			mon.SetHP(tt.newHP)

			if got := mon.HP(); got != tt.wantHP {
				t.Errorf("Monster.HP() after SetHP(%v) = %v, want %v", tt.newHP, got, tt.wantHP)
			}
			if got := mon.Name(); got != tt.wantName {
				t.Errorf("Monster.Name() after SetHP() = %v, want %v", got, tt.wantName)
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
	if got := result.HP(); got != 15 {
		t.Errorf("Monster.WithName().HP() = %v, want %v", got, 15)
	}
	if got := result.MaxHP(); got != 15 {
		t.Errorf("Monster.WithName().MaxHP() = %v, want %v", got, 15)
	}
}

func TestMonster_AdjustHP(t *testing.T) {
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
				m.SetHP(18) // Set to below max first
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
				m.SetHP(5) // Set to low health first
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
				m.SetHP(12) // Set to specific value first
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
			actualAdjustment := monster.AdjustHP(tt.adjustment)

			// Test that we get the actual adjustment amount
			_ = actualAdjustment // We could test this if desired

			if got := monster.HP(); got != tt.wantHP {
				t.Errorf("Monster.HP() after AdjustHP(%v) = %v, want %v", tt.adjustment, got, tt.wantHP)
			}
			if got := monster.Name(); got != tt.wantName {
				t.Errorf("Monster.Name() after AdjustHP() = %v, want %v", got, tt.wantName)
			}
		})
	}
}
