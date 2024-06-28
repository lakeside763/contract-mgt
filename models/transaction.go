package models

import (
	"time"
)

type ContractStatusType BaseType

const (
	CONTRACT_NEW 						ContractStatusType = "NEW"
	CONTRACT_IN_PROGRESS 		ContractStatusType = "IN_PROGRESS"
	CONTRACT_TERMINATED 		ContractStatusType = "TERMINATED"
)

type Contract struct {
	ID          	string              `json:"id" gorm:"type:uuid;default:uuid_generate_v4();column:id"`
	Terms       	string              `json:"terms" gorm:"type:varchar(255);column:terms"`
	Status      	ContractStatusType  `json:"status"`
	ClientID    	string              `json:"client_id" gorm:"type:uuid;column:client_id"`
	ContractorID 	string             	`json:"contractor_id" gorm:"type:uuid;column:contractor_id"`
	CreatedAt  time.Time  						`json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt  time.Time  						`json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

type Job struct {
	ID 						string 							`json:"id" gorm:"type:uuid;default:uuid_generate_v4();column:id"`
	Descripton		string							`json:"description"`
	Price					float64							`json:"price"`
	Paid					bool								`json:"paid"`
	PaymentDate		time.Time						`json:"payment_date" gorm:"type:timestamp;column:payment_date"`
	ContractID		string							`json:"contractId" gorm:"type:uuid;column:contractId"`
	Contract			Contract						`json:"contact"`
	CreatedAt  		time.Time  					`json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt  		time.Time  					`json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

