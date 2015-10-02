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
	fmt.Println(`  ____  __.         __            _________                                 .__           `)
	fmt.Println(` |    |/ _|_____  _/  |_ _____    \_   ___ \   ____    ____    ______ ____  |  |    ____  `)
	fmt.Println(` |      <  \__  \ \   __\\__  \   /    \  \/  /  _ \  /    \  /  ___//  _ \ |  |  _/ __ \ `)
	fmt.Println(` |    |  \  / __ \_|  |   / __ \_ \     \____(  <_> )|   |  \ \___ \(  <_> )|  |__\  ___/ `)
	fmt.Println(` |____|__ \(____  /|__|  (____  /  \______  / \____/ |___|  //____  >\____/ |____/ \___  >`)
	fmt.Println(`         \/     \/            \/          \/              \/      \/                   \/ `)
	fmt.Println(``)
	// Parse flags
	flag.Parse()
	// Parse the config here before doing anything else
	_ = ParseConfig()
	//log.Debug("Grafana URL: ", c.GrafanaUrl)
	//log.Debug("Grafana Auth: ", c.GrafanaAuth)
	//log.Debug("Kata Home: ", c.KataHome)

	// Run the router
	router := NewRouter()
	// Handle a failure
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), router))
}
