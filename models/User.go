package models

type User struct {
	Name     string `json:"name" bind:"required"`
	PassWord string `json:"password,omitempty" bind:"Password"`
	Describe string `json:"Describe,omitempty"`
}
