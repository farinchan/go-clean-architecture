package mail

import (
	"fmt"
	"time"
)

// SendLoginAlert sends a new login security alert email
func (m *Mailer) SendLoginAlert(to, name, ip, userAgent string, loginTime time.Time) error {
	subject := "Security Alert: New Login to Your Account"
	timeStr := loginTime.Format("2006-01-02 15:04:05 MST")

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><style>body{font-family: Arial, sans-serif; color: #333;} .box{border: 1px solid #e0e0e0; padding: 20px; border-radius: 8px; max-width: 500px; margin: auto;} .alert{color: #d9534f; font-weight: bold;}</style></head>
<body>
<div class="box">
	<h2>Security Alert: New Login</h2>
	<p>Hello <strong>%s</strong>,</p>
	<p>We detected a new login to your account with the following details:</p>
	<ul>
		<li><strong>IP Address:</strong> %s</li>
		<li><strong>Device / Browser:</strong> %s</li>
		<li><strong>Time:</strong> %s</li>
	</ul>
	<p>If this was you, you can safely ignore this email.</p>
	<p class="alert">If you did not initiate this login, please change your password immediately.</p>
</div>
</body>
</html>`, name, ip, userAgent, timeStr)

	return m.SendHTML(to, subject, body)
}

// Send2FACode sends a 2FA verification code email
func (m *Mailer) Send2FACode(to, name, code string) error {
	subject := "Your Two-Factor Authentication (2FA) Code"

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><style>body{font-family: Arial, sans-serif; color: #333;} .box{border: 1px solid #e0e0e0; padding: 20px; border-radius: 8px; max-width: 500px; margin: auto;} .code{font-size: 28px; font-weight: bold; letter-spacing: 5px; color: #0275d8; margin: 15px 0;}</style></head>
<body>
<div class="box">
	<h2>Two-Factor Authentication</h2>
	<p>Hello <strong>%s</strong>,</p>
	<p>Use the following code to complete your login. This code is valid for <strong>10 minutes</strong>:</p>
	<div class="code">%s</div>
	<p>Do not share this code with anyone.</p>
</div>
</body>
</html>`, name, code)

	return m.SendHTML(to, subject, body)
}

// SendPasswordResetCode sends a password reset OTP code email
func (m *Mailer) SendPasswordResetCode(to, name, code string) error {
	subject := "Password Reset Verification Code"

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><style>body{font-family: Arial, sans-serif; color: #333;} .box{border: 1px solid #e0e0e0; padding: 20px; border-radius: 8px; max-width: 500px; margin: auto;} .code{font-size: 28px; font-weight: bold; letter-spacing: 5px; color: #f0ad4e; margin: 15px 0;}</style></head>
<body>
<div class="box">
	<h2>Password Reset Request</h2>
	<p>Hello <strong>%s</strong>,</p>
	<p>We received a request to reset your password. Use the verification code below. It will expire in <strong>15 minutes</strong>:</p>
	<div class="code">%s</div>
	<p>If you did not request a password reset, please ignore this email.</p>
</div>
</body>
</html>`, name, code)

	return m.SendHTML(to, subject, body)
}

// SendPasswordChangedAlert sends a password changed confirmation email
func (m *Mailer) SendPasswordChangedAlert(to, name string, changeTime time.Time) error {
	subject := "Security Alert: Your Password Has Been Changed"
	timeStr := changeTime.Format("2006-01-02 15:04:05 MST")

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><style>body{font-family: Arial, sans-serif; color: #333;} .box{border: 1px solid #e0e0e0; padding: 20px; border-radius: 8px; max-width: 500px; margin: auto;} .alert{color: #d9534f; font-weight: bold;}</style></head>
<body>
<div class="box">
	<h2>Password Changed Successfully</h2>
	<p>Hello <strong>%s</strong>,</p>
	<p>Your password was changed at <strong>%s</strong>.</p>
	<p class="alert">If you did not make this change, please contact support immediately to secure your account.</p>
</div>
</body>
</html>`, name, timeStr)

	return m.SendHTML(to, subject, body)
}
