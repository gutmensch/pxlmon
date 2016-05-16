/*
  pxlmon - terminal monitoring and oncall tool

  (c) 2016 Robert Schumann <rs@n-os.org>

  License: BSD
*/
package main


import (
    "fmt"
    "log"

    "github.com/gutmensch/pxlmon/oncall"

    "github.com/kelseyhightower/envconfig"
)


type Configuration struct {
    Debug      bool
    Port       int
    Oncall     string
    Monitoring string
}


func main() {

    var c Configuration

    err := envconfig.Process("pixelmon", &c)
    if err != nil {
        log.Fatal(err.Error())
    }

    fmt.Println(oncall.GetDutyOfficer())
}

