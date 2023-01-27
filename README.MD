# Azure Subscription Switcher
This is a command-line tool for switching between Azure subscriptions and logging in to them.

### `Prerequisites`
In order to run this script, you will need to have Azure CLI - This script utilizes the Azure CLI to switch subscriptions. You can download and install the Azure CLI by following the instructions [here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) and have it added to your PATH.
You need to have the subscription names in your Azure CLI config file. You can find more information about this [here](https://docs.microsoft.com/en-us/cli/azure/manage-azure-subscriptions-azure-cli?view=azure-cli-latest#sign-in).

---

### `Usage`
You can use the following command to see the available options:

```
azss --help
```
---
### `List all subscriptions`
To list all the subscriptions, you can use the --list or -l flag. This will print a table of all the subscriptions and their details.


```
azss --list
```
---
### `Switch to a subscription`
To switch to a specific subscription, you can use the --subscription or -s flag followed by the name of the subscription.

```
azss --subscription "My Subscription"
```
#### `(!) Non-interactive shells`
When running this command in a non-interactive shell, you will need to provide your Azure login credentials. You can do this by using the --username or -u flag for your username and the --password or -p flag for your password.

```
azss --subscription "My Subscription" --username "myusername" --password "mypassword"
```
Alternatively, you can use the --device-code or -d flag to use device code authentication.

```
azss --subscription "My Subscription" --device-code
```
---
<br>

#### `Output`

If the provided subscription name is valid, the command will set the subscription and log in to it. If the provided subscription name is invalid, the command will print an error message and list all the available subscriptions.

---
#### `Error Handling`

If the command encounters an error, it will print the error message and exit with a non-zero status code.