package domain

type Tenant struct {
    BaseModelWithDeleted
    Code          string  `gorm:"size:50;not null;unique"`
    BusinessName  string  `gorm:"size:100;not null"`
    Address       string  `gorm:"not null"`
    Phone         string  `gorm:"size:25;not null"`
    Email         string  `gorm:"not null;unique"`
    InvoiceFooter *string
    IsActive      bool    `gorm:"default:true;not null"`
}