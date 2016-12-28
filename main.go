package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tkanos/serverlogsreplay/replayer"
)

func main() {
	if len(os.Args) < 2 || (len(os.Args)-1)%2 != 0 {
		showHelp()
		return
	}

	args, err := getArgs(os.Args[1:])
	if err {
		showHelp()
		return
	}

	var arguments replayer.Arguments
	if arg, ok := checkArgs(args); !ok {
		showHelp()
		return
	} else {
		arguments = arg
	}

	replayer.Replay(arguments)
}

func getArgs(args []string) (map[string]string, bool) {
	ret := make(map[string]string)
	i := 0
	for i < len(args) {
		if strings.HasPrefix(args[i], "-") {
			ret[args[i]] = args[i+1]
		} else {
			return nil, true
		}
		i = i + 2
	}
	return ret, false
}

func checkArgs(args map[string]string) (replayer.Arguments, bool) {
	var arguments replayer.Arguments

	if len(args) < 8 {
		return arguments, false
	}

	if args["-p"] == "" || args["-d"] == "" || args["-pn"] == "" || args["-qsn"] == "" || args["-bn"] == "" || args["-vn"] == "" || args["-s"] == "" || args["-bt"] == "" {
		return arguments, false
	}

	arguments.File = replayer.FileArgument{FilePath: args["-p"], FileType: args["-ft"]}
	arguments.Http = replayer.HttpArgument{BaseUri: args["-s"], Headers: args["-H"], Cookies: args["-C"]}
	arguments.Regexp = replayer.RegexpArgument{Regexp: args["-mr"], Pattern: args["-mp"], Replace: args["-r"]}
	arguments.Parse = replayer.ParseArgument{Delimiter: args["-d"]}
	// check if it is an interger by converting from string into an int.
	if val, err := strconv.Atoi(args["-pn"]); err != nil {
		return arguments, false
	} else {
		arguments.Parse.UriStemColumn = val
	}

	if val, err := strconv.Atoi(args["-qsn"]); err != nil {
		return arguments, false
	} else {
		arguments.Parse.UriQueryColumn = val
	}

	if val, err := strconv.Atoi(args["-bn"]); err != nil {
		return arguments, false
	} else {
		arguments.Parse.BodyColumn = val
	}

	if val, err := strconv.Atoi(args["-vn"]); err != nil {
		return arguments, false
	} else {
		arguments.Parse.VerbColumn = val
	}

	//check if optional parameters that must be int, are effectively int
	if args["-bl"] != "" {
		if val, err := strconv.Atoi(args["-bl"]); err != nil {
			return arguments, false
		} else {
			arguments.Parse.BeginLine = val
		}
	}

	if args["-uan"] != "" {
		if val, err := strconv.Atoi(args["-uan"]); err != nil {
			return arguments, false
		} else {
			arguments.Parse.UserAgentColumn = val
		}
	}

	if args["-tm"] != "" {
		if val, err := strconv.Atoi(args["-tm"]); err != nil {
			return arguments, false
		} else {
			arguments.ThreadMax = val
		}
	}

	return arguments, true

}

func showHelp() {

	fmt.Println("serverlogsreplay -p path -d delimiter -pn pathNb -qsn queryStringNb -vn verbNr -s server [-ft fileType] [-bl beginLine] [-uan userAgentNb] [-H headers] [-C cookies] [-mr matchRequest] [-mp modifyPattern] [-r replacement] ")
	fmt.Println("\n\n")
	fmt.Println("Mandatory Parameters")
	fmt.Println("\n")
	fmt.Println("-p  \t : path where is located the file/directory")
	fmt.Println("-d  \t : delimiter")
	fmt.Println("-pn \t : int that locate the uri-stem")
	fmt.Println("-s  \t : base uri address")
	fmt.Println("-qsn\t : int that locate the uri-query")
	fmt.Println("-bn\t : int that locate the body for POST")
	fmt.Println("-bt\t : int that locate the body type for POST")
	fmt.Println("-vn \t : int that locate the method (verb)")

	fmt.Println("\n\n")
	fmt.Println("Optional Parameters")
	fmt.Println("\n")
	fmt.Println("-ft \t : if in -p you inform the directory path, ft is needed to inform the filetype (example csv)")
	fmt.Println("-bl \t : int that tell to the program in which line we begin")
	fmt.Println("-uan\t : int that locate the user-agent")
	fmt.Println("-H  \t : to specify others (http)Headers")
	fmt.Println("-C  \t : to specify (http)Cookies")
	fmt.Println("-mr \t : to specify a regexp to execute only the request that match -mr")
	fmt.Println("-mp \t : to specify a regexp pattern for Replace somthing on path or queryString")
	fmt.Println("-r  \t : if you have specified -mp, -r is to specify by what you want to replace your -mp")
	fmt.Println("-tm \t : int representing the Thread max (in parallelization) you want to use, by default it's sequencial (1)")

	fmt.Println("\n\n")
	fmt.Println("Example")
	fmt.Println("\n")
	fmt.Println("serverlogsreplay -p \"D:\\Logs\\mylogs.log\" -d ' ' -pn 3 -qsn 4 -vn 9 -s http://mybetaserver.com")
	fmt.Println("\n")
	fmt.Println("serverlogsreplay -p \"D:\\Logs\" -d ' ' -pn 3 -qsn 4 -vn 9 -s http://mybetaserver.com -ft .log -bl 5")
	fmt.Println("\n")
	fmt.Println("serverlogsreplay -p \"D:\\Logs\" -d ' ' -pn 3 -qsn 4 -vn 9 -s http://mybetaserver.com -ft .log -bl 5 -uan 2 -H \"HeaderName: HeaderValue\" -C \"cookieName1: CookieValue1 CookieName2:CookieValue2\" -mr \"v1\" -mp \"v1\" -r \"v2\" ")

}
