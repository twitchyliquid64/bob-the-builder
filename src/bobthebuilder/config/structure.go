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

	RaspberryPi struct { //Intended to flash LEDs using GPIOs on a raspberry pi. For most uses, set Enable = false. Pin is the pin number as written on the BCM2835 pinout.
		Enable bool
		BuildLedPin int	//LED to flash while a build is running unless Disable Physical Indicators is set on the build.
		DataLedPin int //LED to flash as a phase writes data.
		CycleFlashers []int	//LED pins to cycle through while a build is in progress unless Disable Physical Indicators is set on the build.
	}

	AWS struct {	//Optional settings for S3 integration. Necessary for phases which talk to S3.
		Enable bool
		AccessKey string
		SecretKey string
		Token string
	}
}
