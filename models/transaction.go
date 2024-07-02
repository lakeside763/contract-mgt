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
	ID          	string              `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Terms       	string              `json:"terms" gorm:"type:varchar(255)"`
	Status      	ContractStatusType  `json:"status"`
	ClientID    	string              `json:"client_id" gorm:"type:uuid;column:client_id"`
	ContractorID 	string              `json:"contractor_id" gorm:"type:uuid;column:contractor_id"`
	CreatedAt   	time.Time           `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   	time.Time           `json:"updated_at" gorm:"type:timestamp"`
}

type Job struct {
	ID 						string 							`json:"id" gorm:"type:uuid;default:uuid_generate_v4();column:id"`
	Description		string							`json:"description"`
	Price					float64							`json:"price"`
	Paid					bool								`json:"paid"`
	PaymentDate		time.Time						`json:"payment_date" gorm:"type:timestamp;column:payment_date"`
	ContractID		string							`json:"contract_id" gorm:"type:uuid;column:contract_id"`
	Contract			Contract						`json:"contract"`
	CreatedAt  		time.Time  					`json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt  		time.Time  					`json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

type JobPayment struct {
	JobID			string			`json:"job_id"`
	ClientID	string			`json:"client_id"`
}

type PaymentDeposit struct {
	ClientID 	string		`json:"client_id"`
	Amount		float64	 	`json:"amount"`
}
