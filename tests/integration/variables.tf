variable "login_password" {
  description = "Dummy password for a test login that doesn't exist."
  sensitive   = true
  type        = string
  default     = "KeepMeSecr3t"
}