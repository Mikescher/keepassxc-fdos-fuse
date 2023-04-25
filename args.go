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
	Extra    []string
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

			if i+5 >= len(args) {
				PrintHelp("Missing values for --secret")
				os.Exit(1)
			}

			specs = append(specs, SSSpec{
				Filename: args[i+1],
				IdentKey: args[i+2],
				IdentVal: args[i+3],
				Attr:     args[i+4],
				Extra:    strings.Split(args[i+5], ","),
			})

			i += 5

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
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println(" ./kxc-fdos-fuse --mount [dir] --secret [...] --secret [...] --secret [...]")
	fmt.Println("")
	fmt.Println(" --mount [dir]")
	fmt.Println("")
	fmt.Println(" --secret [filename] [search-attribute-key] [search-attribute-value] [data-attribute-key] [extra]")
	fmt.Println("   // [filename]:                Output filename")
	fmt.Println("   // [search-attribute-key]:    Keepass/Secret-Service Attribute-Key used for searching the correct entry")
	fmt.Println("   // [search-attribute-value]:  Keepass/Secret-Service Attribute-Value used for searching the correct entry")
	fmt.Println("   // [data-attribute-key]:      Keepass/Secret-Service Attribute-Key, which contains the file data")
	fmt.Println("   // [extra]:                   Extra parameters")
	fmt.Println("")
	fmt.Println(" --version")
	fmt.Println("")
	fmt.Println(" --info")
	fmt.Println("")
	fmt.Println(" --help")
	fmt.Println("")
	fmt.Println("Possible [extra] params (comma separated):")
	fmt.Println("  'base64':  base64-decode attribute value")
	fmt.Println("  'plain':   directly output attribute value")
}
