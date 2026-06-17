package workers

import (
    "time"
    "github.com/google/uuid"
)

type Worker struct {
    ID                uuid.UUID  `json:"id"`
    TenantID          uuid.UUID  `json:"tenant_id"`
    EmployeeNumber    string     `json:"employee_number"`
    FullName          string     `json:"full_name"`
    Phone             string     `json:"phone"`
    PinHash           string     `json:"-"`
    SiteID            *uuid.UUID `json:"site_id,omitempty"`
    EmploymentType    string     `json:"employment_type"`
    PaymentSchedule   string     `json:"payment_schedule"`
    HourlyRate        *float64   `json:"hourly_rate,omitempty"`
    ContractEndDate   *time.Time `json:"contract_end_date,omitempty"`
    Status            string     `json:"status"`
    DateOfBirth       *time.Time `json:"date_of_birth,omitempty"`
    BiometricConsent  bool       `json:"biometric_consent"`
    ConsentTimestamp  *time.Time `json:"consent_timestamp,omitempty"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateWorkerRequest struct {
    EmployeeNumber    string     `json:"employee_number" validate:"required"`
    FullName          string     `json:"full_name" validate:"required"`
    Phone             string     `json:"phone" validate:"required"`
    Pin               string     `json:"pin" validate:"required,len=4,numeric"`
    SiteID            *uuid.UUID `json:"site_id,omitempty"`
    EmploymentType    string     `json:"employment_type" validate:"required,oneof=permanent temp contractor shift_worker"`
    PaymentSchedule   string     `json:"payment_schedule" validate:"required,oneof=monthly weekly adhoc hourly"`
    HourlyRate        *float64   `json:"hourly_rate,omitempty"`
    ContractEndDate   *time.Time `json:"contract_end_date,omitempty"`
    DateOfBirth       *time.Time `json:"date_of_birth,omitempty"`
    BiometricConsent  bool       `json:"biometric_consent"`
    ConsentTimestamp  *time.Time `json:"consent_timestamp,omitempty"`
}

type UpdateWorkerRequest struct {
    FullName          *string     `json:"full_name,omitempty"`
    Phone             *string     `json:"phone,omitempty"`
    SiteID            *uuid.UUID  `json:"site_id,omitempty"`
    EmploymentType    *string     `json:"employment_type,omitempty"`
    PaymentSchedule   *string     `json:"payment_schedule,omitempty"`
    HourlyRate        *float64    `json:"hourly_rate,omitempty"`
    ContractEndDate   *time.Time  `json:"contract_end_date,omitempty"`
    Status            *string     `json:"status,omitempty"`
    BiometricConsent  *bool       `json:"biometric_consent,omitempty"`
}

type WorkerResponse struct {
    ID                string     `json:"id"`
    TenantID          string     `json:"tenant_id"`
    EmployeeNumber    string     `json:"employee_number"`
    FullName          string     `json:"full_name"`
    Phone             string     `json:"phone"`
    SiteID            *string    `json:"site_id,omitempty"`
    EmploymentType    string     `json:"employment_type"`
    PaymentSchedule   string     `json:"payment_schedule"`
    HourlyRate        *float64   `json:"hourly_rate,omitempty"`
    ContractEndDate   *time.Time `json:"contract_end_date,omitempty"`
    Status            string     `json:"status"`
    DateOfBirth       *time.Time `json:"date_of_birth,omitempty"`
    BiometricConsent  bool       `json:"biometric_consent"`
    ConsentTimestamp  *time.Time `json:"consent_timestamp,omitempty"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}

type LoginRequest struct {
    Phone string `json:"phone"`
    Pin   string `json:"pin"`
}
