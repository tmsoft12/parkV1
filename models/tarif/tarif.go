package tarif

import "time"

type Tarif struct {
	Id         int       `json:"id"`
	Plate      string    `json:"plate"`
	Name       string    `json:"name"`
	Start_time time.Time `json:"start_time"`
	End_time   time.Time `json:"end_time"`
	Price      int       `json:"price"`
}
