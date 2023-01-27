package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

type AzureOutput struct {
	CloudName        string        `json:"cloudName"`
	HomeTenantID     string        `json:"homeTenantId"`
	ID               string        `json:"id"`
	IsDefault        bool          `json:"isDefault"`
	ManagedByTenants []interface{} `json:"managedByTenants"`
	Name             string        `json:"name"`
	State            string        `json:"state"`
	TenantID         string        `json:"tenantId"`
	User             User          `json:"user"`
}
type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {
	var azureSubscriptions []AzureOutput
	var subscriptionName string
	var username string
	var password string
	var list bool
	var deviceCode bool

	app := cli.NewApp()
	app.Name = "Azure Subscription Switcher"
	app.Usage = "Switch Azure subscriptions and login"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "subscription, s",
			Usage:       "Name of the subscription to switch to",
			Destination: &subscriptionName,
		},
		cli.StringFlag{
			Name:        "username, u",
			Usage:       "Username for Azure login | Required for non-interactive shells",
			Destination: &username,
		},
		cli.StringFlag{
			Name:        "password, p",
			Usage:       "Password for Azure login | Required for non-interactive shells",
			Destination: &password,
		},
		cli.BoolFlag{
			Name:        "device-code, d",
			Usage:       "Use device code authentication | Required for non-interactive shells",
			Destination: &deviceCode,
		},
		cli.BoolFlag{
			Name:        "list, l",
			Usage:       "List all subscriptions",
			Destination: &list,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// This code is used to set the default Azure subscription.
	// It uses the Azure CLI to find the list of subscriptions, and then uses the subscription name to find the ID of the desired subscription.
	// It then sets the subscription using the Azure CLI.
	subscriptions, _ := exec.Command("az", "account", "list").Output()
	json.Unmarshal(subscriptions, &azureSubscriptions)
	var subscriptionMatch AzureOutput
	if subscriptionName != "" {
		for _, subscription := range azureSubscriptions {
			if subscription.Name == subscriptionName {
				subscriptionMatch = subscription
				break
			}
		}
	}
	switch {
	case list:
		ListSubscriptions(azureSubscriptions)
		return
	case subscriptionMatch.Name == "" && subscriptionName != "":
		os.Stdout.WriteString("\033[1m")
		fmt.Println("\nPlease provide a valid subscription name.")
		os.Stdout.WriteString("\033[0m")
		ListSubscriptions(azureSubscriptions)
		cli.ShowAppHelpAndExit(cli.NewContext(app, nil, nil), 0)
		return
	case subscriptionMatch.Name == "" && subscriptionName == "":
		os.Stdout.WriteString("\033[1m")
		fmt.Printf("\n\n[!] No subscription provided.\n\n")
		os.Stdout.WriteString("\033[0m")
		return
	default:
		AzureSetSubscription(subscriptionMatch, deviceCode, username, password)
	}
}

// AzureSetSubscription sets the Azure subscription and logs in
func AzureSetSubscription(subscriptionMatch AzureOutput, deviceCode bool, username string, password string) {
	exec.Command("az", "account", "set", "--subscription", subscriptionMatch.ID, "--tenant", subscriptionMatch.TenantID).Run()

	switch {
	case deviceCode:
		exec.Command("az", "login", "--use-device-code", "--tenant", subscriptionMatch.TenantID).Run()
	case username != "" && password != "":
		exec.Command("az", "login", "-u", username, "-p", password, "--tenant", subscriptionMatch.TenantID).Run()
	default:
		cmd := exec.Command("az", "login")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

// ListSubscriptions prints a list of Azure subscriptions
// The function takes a list of AzureOutput objects, and prints a table with the name, user, and state of each subscription
func ListSubscriptions(azureSubscriptions []AzureOutput) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "User", "State"})
	for _, subscription := range azureSubscriptions {
		table.Append([]string{subscription.Name, subscription.User.Name, subscription.State})
	}
	table.Render()
}
