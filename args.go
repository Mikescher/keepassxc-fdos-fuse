package kxcfuse

import (
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"strings"
)

type SSSpec struct {
	Filename string
	IdentKey string
	IdentVal string
	Attr     string
}

type ProgArgs struct {
	Spec     []SSSpec
	MountDir string
}

func ParseArgs(args []string) ProgArgs {
	var dir *string = nil
	var specs = make([]SSSpec, 0)

	for i := 0; i < len(args); {

		arg := args[i]

		if strings.ToLower(arg) == "--help" {

			PrintHelp("")
			os.Exit(0)

		} else if strings.ToLower(arg) == "--version" {

			fmt.Println(Version)
			os.Exit(0)

		} else if strings.ToLower(arg) == "--info" {

			fmt.Printf("Version        := %s\n", Version)
			fmt.Printf("VCSType        := %s\n", langext.Coalesce(VCSType, "?"))
			fmt.Printf("CommitHash     := %s\n", langext.Coalesce(CommitHash, "?"))
			fmt.Printf("CommitTime     := %s\n", langext.Coalesce(CommitTime, "?"))
			fmt.Printf("CommitModified := %s\n", langext.Coalesce(CommitModified, "?"))
			os.Exit(0)

		} else if strings.ToLower(arg) == "--mount" {

			if i+1 >= len(args) {
				PrintHelp("Missing values for --mount")
				os.Exit(1)
			}

			dir = langext.Ptr(args[i+1])

			i += 1

		} else if strings.ToLower(arg) == "--secret" {

			if i+4 >= len(args) {
				PrintHelp("Missing values for --secret")
				os.Exit(1)
			}

			specs = append(specs, SSSpec{
				Filename: args[i+1],
				IdentKey: args[i+2],
				IdentVal: args[i+3],
				Attr:     args[i+4],
			})

			i += 4

		} else {

			PrintHelp(fmt.Sprintf("Unknown parameter: '%s'", arg))
			os.Exit(1)

		}

		i++
	}

	if dir == nil {
		PrintHelp("Missing parameter `--mount [dir]`")
		os.Exit(1)
	}

	return ProgArgs{
		Spec:     specs,
		MountDir: *dir,
	}
}

func PrintHelp(err string) {
	if err != "" {
		fmt.Println(err)
		fmt.Println("")
	}
	fmt.Println("Usage: ./kxc-fdos-fuse --mount [dir] --secret [filename] [search-attribute-key] [search-attribute-value] [data-attribute-key]")
	fmt.Println(" --version")
	fmt.Println(" --info")
	fmt.Println(" --help")
	fmt.Println("")
}
