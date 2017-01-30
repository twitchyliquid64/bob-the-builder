package config



type Config struct {
	Name string							//canonical name to help identify this server
	TLS struct {						//Relative file addresses of the .pem files needed for TLS.
		PrivateKey string
		Cert string
	}

	Web struct{							//Details needed to get the website part working.
		Domain string					//Domain should be in the form example.com
		Listener string				//Address:port (address can be omitted) where the HTTPS listener will bind.
		RequireBasicAuth bool //If set, will require HTTP Basic authentication with one of the user:pass pairs in AuthPairs
		AuthPairs map[string]string
	}

	RaspberryPi struct {
		Enable bool
		BuildLedPin int
		DataLedPin int
		CycleFlashers []int
	}

	AWS struct {
		Enable bool
		AccessKey string
		SecretKey string
		Token string
	}
}
