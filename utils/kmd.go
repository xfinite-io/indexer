package utils

import (
	"fmt"
	"os"
	"crypto/sha256"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/algorand/go-algorand-sdk/crypto"
    	"github.com/algorand/go-algorand-sdk/mnemonic"
)

// These constants represent the kmdd REST endpoint and the corresponding API
// token. You can retrieve these from the `kmd.net` and `kmd.token` files in
// the kmd data directory.
var kmdAddress = os.Getenv("KMD_ADDRESS")
var kmdToken = os.Getenv("KMD_TOKEN")

func CreateUserAlgoAddress() (string, error) {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return "", err
	}
	fmt.Println("Made a kmd client")

	// Create the example wallet, if it doesn't already exist
	cwResponse, err := kmdClient.CreateWallet("testwallet", "testpassword", kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		return "", err
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	exampleWalletID := cwResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, exampleWalletID)

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err := kmdClient.InitWalletHandle(exampleWalletID, "testpassword")
	if err != nil {
		return "", err
	}

	// Extract the wallet handle
	exampleWalletHandleToken := initResponse.WalletHandleToken

	// Generate a new address from the wallet handle
	genResponse, err := kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		return "", err
	}
	return genResponse.Address, nil
}

func CreateUserStandaloneAccount(user_id string) (types.Address, error) {
	account := crypto.GenerateAccount()
	passphrase, err := mnemonic.FromPrivateKey(account.PrivateKey)

	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err) 
	} else {
		fmt.Printf("My address: %s\n", account.Address)
		fmt.Printf("My passphrase: %s\n", passphrase)
	}

	hash := sha256.Sum256(append([]byte(user_id), []byte("mzaalo")...))
	_, err = SetSecret("algo", fmt.Sprintf("%s_publickey", hash), account.Address.String())
	if err != nil {
                fmt.Printf("Error adding public key to kms: %s\n", err)
	}
	_, err = SetSecret("algo", fmt.Sprintf("%s_privatekey", hash), passphrase)
	if err != nil {
                fmt.Printf("Error adding private key to kms: %s\n", err)
	}
	return account.Address, nil
}
