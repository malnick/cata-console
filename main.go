package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

// Define flags
var verbose = flag.Bool("v", false, "Define verbose logging output.")
var port = flag.String("p", "9000", "Port to run.")

func main() {
	// Banner
	fmt.Println(`                              .(%@@@@@@@@@@&%/.                                 	`)
	fmt.Println(`                       *%@@@@@@@@@@@@@@@@@@@@@@@@@@@@%,                        		`)
	fmt.Println(`                 .#@@@@@@@@@@@@@@&%(//*//(%&@@@@@@@@@@@@@@@&*                  		`)
	fmt.Println(`         .(&@@@@@@@@@@@%,                         .*&@@@@@@@@@@@@#.           		`)
	fmt.Println(`@@@@@@@@@@@@@@@@@@%,       ./&@@@@@@@@@@@@@@@@@@%*        ./%@@@@@@@@@@@@&&&&@&		`)
	fmt.Println(`@@@@@@@&%/.         .(@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&*         .*%&@@@@@@@@		`)
	fmt.Println(`               (&@@@@@@@@@&@@@@@@@@@@@@@@@@@@     *%@@@@@@@@@@&*               		`)
	fmt.Println(`@#,//(%&@@@@@@@@@@@%*     @@@@@@@@@@@@@@@@@@@&           .(@@@@@@@@@@&*        		`)
	fmt.Println(`@@@@@@@@@@@@@&*           %@@@@@@@@@@@@@@@@@@@                 ,#@@@@@@@@@@@@@@		`)
	fmt.Println(`     ,(&@@@@@@@@#          @@@@@@@@@@@@@@@@@@#                     (@@@@@@@@@@@		`)
	fmt.Println(`           .#@@@@@@@&.      &@@@@@@@@@@@@@@@(                 *&@@@@@@@&(,./&@@		`)
	fmt.Println(`                &@@@@@@@@@#.  .&@@@@@@@@@&            ./%@@@@@@@@@(           .		`)
	fmt.Println(`                 (@@@@@@@@@@@@@@@/*,.....,*(%&@@@@@@@@@@@@@#.                			`)
	fmt.Println(`                 *@@@. &@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%,                       		`)
	fmt.Println(`                 @@@@    *@@@@&                                                		`)
	fmt.Println(`      %*.      #@@@@@      .@@@@@                                              		`)
	fmt.Println(`       .@@@@@@@@@@@@@        ,@@@@&                                            		`)
	fmt.Println(`         ,@@@@@@@@@@@.         (@@@@@                             .,,,.        		`)
	fmt.Println(`           @@@@@@@@@@.           &@@@@&                        &@@@@@@@@@@.    		`)
	fmt.Println(`            .@@@@@@@@,             @@@@@@                    &@@@@@%#%@@@@@@.  		`)
	fmt.Println(`              @@@@@@@,               @@@@@&                 (@@@@@@@@.  *@@@@  		`)
	fmt.Println(`               %@@@@@*                 &@@@@@.              @@@@@@@@@@   &@@@* 		`)
	fmt.Println(`                ,@@@@,                   %@@@@@*            (@@@@@@@@@   &@@@( 		`)
	fmt.Println(`                 .@@@.                     ,@@@@@@(          *@@@@@@@.  .@@@@  		`)
	fmt.Println(`                  .@@                         %@@@@@@&,                &@@@@*  		`)
	fmt.Println(`                   .@                            /@@@@@@@@@%(,.  .,#&@@@@@&    		`)
	fmt.Println(`                    *                                #@@@@@@@@@@@@@@@@@@#      		`)
	fmt.Println(``)
	fmt.Println(`|/  _  _ _|_ __ _  |    /   _ __  _  _  |  _ `)
	fmt.Println(`|\ (/__>  |_ | (/_ |    \__(_)| |_> (_) | (/_`)
	fmt.Println(``)
	// Parse flags
	flag.Parse()
	// Parse the config here before doing anything else
	_ = ParseConfig()
	// Run the router
	router := NewRouter()
	// Handle a failure
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), router))
}
