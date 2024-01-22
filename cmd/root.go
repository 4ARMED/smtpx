/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"

	"github.com/spf13/cobra"
)

type SmtpOptions struct {
	Server     string
	Port       int
	Username   string
	Password   string
	Filename   string
	Sender     string
	Recipients []string
	Verbose    bool
}

func Execute() {

	o := SmtpOptions{}

	var rootCmd = &cobra.Command{
		Use:   "smtpx",
		Short: "SMTP utility to send Raw SMTP messages",
		Long:  `SMTP utility to send Raw SMTP messages`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Send(); err != nil {
				return fmt.Errorf("error sending SMTP message: %s", err)
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&o.Filename, "file", "f", "-", "Raw file containing SMTP message")
	rootCmd.Flags().StringVarP(&o.Server, "server", "s", "", "SMTP server hostname")
	rootCmd.Flags().IntVarP(&o.Port, "port", "p", 25, "SMTP server port")
	rootCmd.Flags().StringVar(&o.Username, "username", "", "SMTP username")
	rootCmd.Flags().StringVar(&o.Password, "password", "", "SMTP password")
	rootCmd.Flags().StringSliceVar(&o.Recipients, "to", []string{}, "Recipient email addresses (specify multiple times or use comma separated list)")
	rootCmd.Flags().StringVar(&o.Sender, "from", "", "Sender email address")
	rootCmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "Verbose output")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func (o *SmtpOptions) Send() error {
	if o.Server == "" {
		return fmt.Errorf("no SMTP server specified")
	}

	var auth smtp.Auth

	if o.Username != "" {
		auth = smtp.PlainAuth("", o.Username, o.Password, o.Server)
	}

	if o.Sender == "" {
		return fmt.Errorf("no sender specified")
	}

	if len(o.Recipients) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	var msg []byte
	var err error

	switch o.Filename {
	case "-":
		msg, err = io.ReadAll(os.Stdin)
	case "":
		return fmt.Errorf("no filename specified")
	default:
		msg, err = os.ReadFile(o.Filename)
	}

	if err != nil {
		return fmt.Errorf("error reading SMTP message: %s", err)
	}

	log.Printf("Sending SMTP message to %s", o.Recipients)
	if o.Verbose {
		log.Println(string(msg))
	}

	return smtp.SendMail(fmt.Sprintf("%s:%d", o.Server, o.Port), auth, o.Sender, o.Recipients, msg)
}
