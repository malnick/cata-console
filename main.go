package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

// Define flags
var verbose = flag.Bool("v", false, "Define verbose logging output.")

func main() {
	// Banner
	fmt.Printf(`                              .(%@@@@@@@@@@&%/.                                 	`)
	fmt.Printf(`                       *%@@@@@@@@@@@@@@@@@@@@@@@@@@@@%,                        		`)
	fmt.Printf(`                 .#@@@@@@@@@@@@@@&%(//*//(%&@@@@@@@@@@@@@@@&*                  		`)
	fmt.Printf(`         .(&@@@@@@@@@@@%,                         .*&@@@@@@@@@@@@#.           		`)
	fmt.Printf(`@@@@@@@@@@@@@@@@@@%,       ./&@@@@@@@@@@@@@@@@@@%*        ./%@@@@@@@@@@@@&&&&@&		`)
	fmt.Printf(`@@@@@@@&%/.         .(@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&*         .*%&@@@@@@@@		`)
	fmt.Printf(`               (&@@@@@@@@@&@@@@@@@@@@@@@@@@@@     *%@@@@@@@@@@&*               		`)
	fmt.Printf(`@#,//(%&@@@@@@@@@@@%*     @@@@@@@@@@@@@@@@@@@&           .(@@@@@@@@@@&*        		`)
	fmt.Printf(`@@@@@@@@@@@@@&*           %@@@@@@@@@@@@@@@@@@@                 ,#@@@@@@@@@@@@@@		`)
	fmt.Printf(`     ,(&@@@@@@@@#          @@@@@@@@@@@@@@@@@@#                     (@@@@@@@@@@@		`)
	fmt.Printf(`           .#@@@@@@@&.      &@@@@@@@@@@@@@@@(                 *&@@@@@@@&(,./&@@		`)
	fmt.Printf(`                &@@@@@@@@@#.  .&@@@@@@@@@&            ./%@@@@@@@@@(           .		`)
	fmt.Printf(`                 (@@@@@@@@@@@@@@@/*,.....,*(%&@@@@@@@@@@@@@#.                			`)
	fmt.Printf(`                 *@@@. &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%,                       		`)
	fmt.Printf(`                 @@@@    *@@@@&                                                		`)
	fmt.Printf(`      %*.      #@@@@@      .@@@@@                                              		`)
	fmt.Printf(`       .@@@@@@@@@@@@@        ,@@@@&                                            		`)
	fmt.Printf(`         ,@@@@@@@@@@@.         (@@@@@                             .,,,.        		`)
	fmt.Printf(`           @@@@@@@@@@.           &@@@@&                        &@@@@@@@@@@.    		`)
	fmt.Printf(`            .@@@@@@@@,             @@@@@@                    &@@@@@%#%@@@@@@.  		`)
	fmt.Printf(`              @@@@@@@,               @@@@@&                 (@@@@@@@@.  *@@@@  		`)
	fmt.Printf(`               %@@@@@*                 &@@@@@.              @@@@@@@@@@   &@@@* 		`)
	fmt.Printf(`                ,@@@@,                   %@@@@@*            (@@@@@@@@@   &@@@( 		`)
	fmt.Printf(`                 .@@@.                     ,@@@@@@(          *@@@@@@@.  .@@@@  		`)
	fmt.Printf(`                  .@@                         %@@@@@@&,                &@@@@*  		`)
	fmt.Printf(`                   .@                            /@@@@@@@@@%(,.  .,#&@@@@@&    		`)
	fmt.Printf(`                    *                                #@@@@@@@@@@@@@@@@@@#      		`)

	// Parse flags
	flag.Parse()
	// Parse the config here before doing anything else
	_ = ParseConfig()
	// Run the router
	router := NewRouter()
	// Handle a failure
	log.Fatal(http.ListenAndServe(":9000", router))
}
