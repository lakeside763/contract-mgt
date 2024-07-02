package handlers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakeside763/contract-mgt/models"
	"gorm.io/gorm"
)

func GetContracts(c *gin.Context, db *gorm.DB) {
	var contracts []models.Contract

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user ID in context"})
		return
	}

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user role in context"})
		return
	}

	switch userRole {
	case "CONTRACTOR":
		query := `SELECT * FROM contracts WHERE "contractorId" = ?`
		if err := db.Where(query, userID).Find(&contracts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contractor ID"})
			return
		}
	case "CLIENT":
		query := `SELECT * FROM contracts WHERE "clientId" = ?`
		if err := db.Raw(query, userID).Find(&contracts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return
		}
	default:
		if err := db.Find(&contracts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, contracts)
}


func GetContract(c *gin.Context, db *gorm.DB) {
	contractID := c.Param("id")
	contract := models.Contract{ID: contractID}

	if err := db.First(&contract).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contract)
}

func CreateContract(c *gin.Context, db *gorm.DB) {
	var contract models.Contract

	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&contract).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func CreateJob(c *gin.Context, db *gorm.DB) {
	var job models.Job
	
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var contract models.Contract
	if err := db.First(&contract, "id = ?", job.ContractID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	if err := db.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	job.Contract = contract

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context, db *gorm.DB) {
	var jobs []models.Job

	if err := db.Preload("Contract").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func JobPayment(c *gin.Context, db *gorm.DB) {
	var jobPayment models.JobPayment

	if err := c.ShouldBindJSON(&jobPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var job = models.Job{ID: jobPayment.JobID}
	if err := db.Preload("Contract").First(&job, "id = ?", jobPayment.JobID).Error;  err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid job ID was provided"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(job.Contract)

	if job.Paid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The payment for the job has already been processed"})
		return
	}

	var client models.Profile
	if err := db.First(&client, "id = ?", jobPayment.ClientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID was provided"})
		return
	}

	if client.Balance < job.Price {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient client balance"})
		return
	}

	var contractor models.Profile
	if err := db.First(&contractor, "id = ?", job.Contract.ContractorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contractor profile not found"})
		return
	}

	clientUpdatedBalance := math.Round((client.Balance - job.Price) * 100) / 100
	contractorUpdatedBalance := math.Round((contractor.Balance + job.Price) * 100) / 100

	err := db.Transaction(func(tx *gorm.DB) error {
		// update job status
		if err := tx.Model(&job).Update("paid", true).Error; err != nil {
			return err
		}

		// update client balance
		if err := tx.Model(&client).Update("balance", clientUpdatedBalance).Error; err != nil {
			return err
		}

		// update contractor balance
		if err := tx.Model(&contractor).Update("balance", contractorUpdatedBalance).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Job payment processed successfully"})
}

func PaymentDeposit(c *gin.Context, db *gorm.DB) {
	var deposit models.PaymentDeposit

	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var client models.Profile
	if err := db.First(&client, "id = ?", deposit.ClientID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Client ID was provided"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Calculate the total sum of unpaid jobs
	var totalSum float64
	err := db.Model(&models.Job{}).
		Joins("JOIN contracts ON contracts.id = jobs.contract_id").
		Where("contracts.client_id = ? AND contracts.status = ?", deposit.ClientID, "IN_PROGRESS").
		Select("SUM(jobs.price)").
		Scan(&totalSum).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_25Percent := totalSum * 0.25
	if client.Balance + deposit.Amount > _25Percent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deposit cannot be more than 25% of total unpaid job prices"})
		return
	}

	updatedBalance := math.Round((client.Balance + deposit.Amount) * 100) / 100
	err = db.Model(&client).Update("balance", updatedBalance).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deposited successfully"})
}