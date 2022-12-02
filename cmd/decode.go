/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	vin    string
	fields string
)

type Response struct {
	Count          int
	Message        string
	SearchCriteria string
	Results        []Result
}

type Result struct {
	Value      string
	ValueID    string
	Variable   string
	VariableId int
}

func call_api(vin string, outputType string) []byte {
	url := "https://vpic.nhtsa.dot.gov/api/vehicles/DecodeVin/" + vin
	response, err := http.Get(url + "/?format=" + outputType)

	if err != nil {
		fmt.Println(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

func format_response_from_data(data []byte) Response {
	var responseObject Response
	json.Unmarshal(data, &responseObject)

	return responseObject
}

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a vehicle's VIN",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("decode called")

		rawFlag, _ := cmd.Flags().GetBool("raw")
		yamlFlag, _ := cmd.Flags().GetBool("yaml")
		metaFlag, _ := cmd.Flags().GetBool("meta")
		sparseFlag, _ := cmd.Flags().GetBool("sparse")
		fieldsFlag, _ := cmd.Flags().GetString("fields")

		var response Response
		if yamlFlag {
			fmt.Println("Calling api with Yaml")
			data := call_api(vin, "yaml")
			fmt.Println(string(data))
		} else if rawFlag {
			fmt.Println("Calling api with Json")
			data := call_api(vin, "json")
			fmt.Println(string(data))
		}

		data := call_api(vin, "json")
		response = format_response_from_data(data)

		if metaFlag {
			count := fmt.Sprint("There were ", response.Count, " items returned")
			fmt.Println(count)
		}
		if sparseFlag {
			for _, field := range response.Results {
				if field.Value != "" && field.Value != "Not Applicable" {
					fmt.Println(field)
				}
			}
		}
		if fieldsFlag != "" {
			fields := strings.Split(fieldsFlag, ",")

			for _, field := range response.Results {
				if field.Value == "" {
					continue
				}
				for _, query := range fields {
					cleaned_query := strings.ReplaceAll(query, " ", "")
					fully_cleaned := strings.TrimSuffix(cleaned_query, ".*")
					if strings.Contains(field.Variable, fully_cleaned) {
						fmt.Println(field)
					}
				}
			}
		}
	},
}

func init() {

	decodeCmd.Flags().StringVarP(&vin, "vin", "v", "", "The VIN to decode (required)")
	decodeCmd.Flags().BoolP("raw", "r", false, "Output the raw JSON response of the API")
	decodeCmd.Flags().BoolP("yaml", "y", false, "Output the results in YAML format")
	decodeCmd.Flags().BoolP("meta", "m", false, "Output how many fields were returned in the response")
	decodeCmd.Flags().BoolP("sparse", "s", false, "Only output fields that have data in them and aren't empty or null")
	// fields Example: --fields="Manufacturer Name, Make, Model, Plant.*"
	decodeCmd.Flags().StringVarP(&fields, "fields", "f", "", "Return fields that match the expression provied to the argument")

	if err := decodeCmd.MarkFlagRequired("vin"); err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(decodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
