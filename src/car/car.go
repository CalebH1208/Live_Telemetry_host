package car

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Telem_value struct {
	Name        string  `json:"N"`
	Unit        string  `json:"U"`
	Value       float64 `json:"V"`
	Filer_value float64 `json:"F"`
}

func (tv *Telem_value) update_value(newval float64) {
	tv.Filer_value = tv.Value
	tv.Value = newval
}

type Car struct {
	Car_num      int `json:"CN"`
	last_active  time.Time
	Active       bool          `json:"AF"`
	Telem_values []Telem_value `json:"TV"`
}

// func (c *Car) update_Car_Num(newNum int) {
// 	c.Car_num = newNum
// }

func (c *Car) update_last_active() {
	c.last_active = time.Now()
}

func (c *Car) Update_active_flag(timeout time.Duration) {
	current_time := time.Now()
	elapsed_time := current_time.Sub(c.last_active)
	if elapsed_time > timeout {
		c.Active = false
	} else {
		c.Active = true
	}
}

// func (c *Car) add_Telem_Value(tv Telem_value) {
// 	c.Telem_values = append(c.Telem_values, tv)
// }

func (c *Car) Serialize() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Car) Update_Car(data string) error {
	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return errors.New("empty input")
	}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return fmt.Errorf("invalid telemetry format in line: %s", line)
		}
		name := strings.TrimSpace(parts[0])
		unit := strings.TrimSpace(parts[1])
		valStr := strings.TrimSpace(parts[2])
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return fmt.Errorf("invalid value for telemetry '%s': %v", name, err)
		}

		found := false
		for idx := range c.Telem_values {
			if c.Telem_values[idx].Name == name {
				// Update existing telemetry value.
				c.Telem_values[idx].update_value(val)
				// If unit differs, update it.
				if c.Telem_values[idx].Unit != unit {
					c.Telem_values[idx].Unit = unit
				}
				found = true
				break
			}
		}

		if !found {
			newTelem := Telem_value{
				Name:        name,
				Unit:        unit,
				Value:       val,
				Filer_value: 0, // Initialize as needed.
			}
			c.Telem_values = append(c.Telem_values, newTelem)
		}
	}

	c.update_last_active()

	return nil

}
