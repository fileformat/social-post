package cmd

import (
	"fmt"
	"net/smtp"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	from          string
	smtp_host     string
	smtp_port     int
	smtp_username string
	subject       string
	toAddresses   []string
)

var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Send an email",
	Long:  `Send an email via SMTP\n\nNote:\n * the SMTP password is read from the SMTP_PASSWORD environment variable\n * the to addresses are only read from the arguments.`,
	Args:  cobra.ExactArgs(1),
	RunE:  emailSend,
}

func InitEmail() {
	slog.Debug("email init")
	rootCmd.AddCommand(emailCmd)

	emailCmd.Flags().StringVar(&from, "from", "", "From email address")
	viper.BindPFlag("from", emailCmd.PersistentFlags().Lookup("from"))

	emailCmd.Flags().StringVar(&smtp_host, "smtp-host", "", "SMTP host")
	viper.BindPFlag("smtp-host", emailCmd.PersistentFlags().Lookup("smtp-host"))

	emailCmd.Flags().IntVar(&smtp_port, "smtp-port", 25, "SMTP port")
	viper.BindPFlag("smtp-port", emailCmd.PersistentFlags().Lookup("smtp-port"))

	emailCmd.Flags().StringVar(&smtp_username, "smtp-username", "", "SMTP username")
	viper.BindPFlag("smtp-username", emailCmd.PersistentFlags().Lookup("smtp-username"))

	emailCmd.Flags().StringVar(&subject, "subject", "", "Subject")
	viper.BindPFlag("subject", emailCmd.PersistentFlags().Lookup("subject"))

	viper.BindEnv("SMTP_PASSWORD")

	emailCmd.Flags().StringSliceVar(&toAddresses, "to", []string{}, "Destination email addresses")
	viper.BindPFlag("to", emailCmd.PersistentFlags().Lookup("to"))
}

func emailSend(cmd *cobra.Command, args []string) error {

	from = viper.GetString("from")
	subject = viper.GetString("subject")
	//LATER: why doesn't this work? toAddresses = viper.GetStringSlice("to")
	smtp_host = viper.GetString("smtp_host")
	smtp_port = viper.GetInt("smtp_port")
	smtp_username = viper.GetString("smtp_username")

	smtp_password := viper.GetString("SMTP_PASSWORD")
	var auth smtp.Auth = nil
	if smtp_username != "" && smtp_password != "" {
		auth = smtp.PlainAuth("", smtp_username, smtp_password, smtp_host)
	} else {
		slog.Debug("No SMTP username or password provided, not using authentication")
	}

	body, bodyErr := getInput(args[0])
	if bodyErr != nil {
		return bodyErr
	}

	msg := []byte("From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	slog.Info("Sending email", "to", toAddresses, "from", from, "subject", subject,
		"smtp_host", smtp_host, "smtp_port", smtp_port, "smtp_username", smtp_username, "smtp_password", Mask(smtp_password))

	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtp_host, smtp_port), auth, from, toAddresses, msg)

	return err
}
