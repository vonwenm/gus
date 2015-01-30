package configure

// These routines allow interactive reading/writing of the configuration file
import (
	"fmt"
	"reflect"
)
type Cli struct {}

func (c * Cli )PromptForFields(config *Configure){
	fmt.Println("General service configuration")
	c.PromptForStructFields( config.Service)

	fmt.Println("Fields for the User store driver")
	c.PromptForStructFields(config.User)

	fmt.Println("Fields for the Client store driver")
	c.PromptForStructFields(config.Client)
}

func (c *Cli) PromptForStructFields( substruct  interface{}){
	val := reflect.Indirect( reflect.ValueOf(substruct))

	//val := reflect.ValueOf(substruct).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		help := tag.Get("help")
		name := tag.Get("name")
		// typeField.Name

		fmt.Printf("%2d: Field: %s (%s) Current value: %v\n", i+1, name , valueField.Type().Name(), valueField.Interface())
		fmt.Printf(".     (%s)\n", help )
		//fmt.Printf("%2d. Field %s (%s) Current value: \n", i, v.Field(i).Name(), v.Field(i).Type().Name())
	}
}
